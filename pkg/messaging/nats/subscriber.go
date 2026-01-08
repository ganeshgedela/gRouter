package nats

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

// NATSSubscriber handles message subscriptions
type NATSSubscriber struct {
	client        *Client
	source        string
	validator     Validator
	subscriptions []*nats.Subscription
	middleware    []SubscriberMiddleware
	mu            sync.Mutex
	wg            sync.WaitGroup
	drainTimeout  time.Duration
}

// NewSubscriber creates a new subscriber
func NewSubscriber(client *Client) Subscriber {
	timeout := client.config.DrainTimeout
	if timeout == 0 {
		timeout = 30 * time.Second // Default safe timeout
	}

	return &NATSSubscriber{
		client:        client,
		source:        client.source,
		subscriptions: make([]*nats.Subscription, 0),
		middleware:    make([]SubscriberMiddleware, 0),
		drainTimeout:  timeout,
	}
}

// Use adds middleware to the subscriber
func (s *NATSSubscriber) Use(mw ...SubscriberMiddleware) {
	s.middleware = append(s.middleware, mw...)
}

// SetValidator sets the validator for the subscriber
func (s *NATSSubscriber) SetValidator(v Validator) {
	s.validator = v
}

// Subscribe subscribes to a subject with a handler
func (s *NATSSubscriber) Subscribe(ctx context.Context, subject string, handler HandlerFunc, opts *SubscribeOptions) (*nats.Subscription, error) {

	// Setup concurrency control if MaxWorkers is set
	var sem chan struct{}
	if opts != nil && opts.MaxWorkers > 0 {
		sem = make(chan struct{}, opts.MaxWorkers)
	}

	// Create message handler wrapper
	msgHandler := func(msg *nats.Msg) {
		s.wg.Add(1)
		defer s.wg.Done()

		if sem != nil {
			sem <- struct{}{}
			defer func() { <-sem }()
		}

		// Unmarshal envelope
		var envelope MessageEnvelope
		if err := json.Unmarshal(msg.Data, &envelope); err != nil {
			s.client.logger.Error("Failed to unmarshal message",
				zap.Error(err),
				zap.String("subject", msg.Subject),
			)
			return
		}

		// Capture NATS reply subject for request-reply
		if msg.Reply != "" {
			envelope.Reply = msg.Reply
		}

		// Apply middleware
		h := handler
		for i := len(s.middleware) - 1; i >= 0; i-- {
			h = s.middleware[i](h)
		}

		// Handle message (tracing is handled in middleware)
		if err := h(context.Background(), msg.Subject, &envelope); err != nil {
			s.client.logger.Error("Handler error",
				zap.Error(err),
				zap.String("subject", msg.Subject),
				zap.String("message_id", envelope.ID),
			)
		}
	}

	var sub *nats.Subscription
	var err error

	// Subscribe with or without queue group
	if opts != nil && opts.QueueGroup != "" {
		sub, err = s.client.Conn().QueueSubscribe(subject, opts.QueueGroup, msgHandler)
	} else {
		sub, err = s.client.Conn().Subscribe(subject, msgHandler)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to subscribe: %w", err)
	}

	// Store subscription
	s.mu.Lock()
	s.subscriptions = append(s.subscriptions, sub)
	s.mu.Unlock()

	s.client.logger.Debug("Subscribed to subject",
		zap.String("subject", subject),
		zap.String("queue_group", func() string {
			if opts != nil {
				return opts.QueueGroup
			}
			return ""
		}()),
	)

	return sub, nil
}

// Unsubscribe unsubscribes from all subscriptions
func (s *NATSSubscriber) Unsubscribe(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, sub := range s.subscriptions {
		if sub.IsValid() {
			if err := sub.Unsubscribe(); err != nil {
				if err != nats.ErrBadSubscription && err != nats.ErrConnectionClosed {
					s.client.logger.Error("Failed to unsubscribe", zap.Error(err))
				}
			}
		}
	}

	s.subscriptions = make([]*nats.Subscription, 0)
	s.client.logger.Debug("Unsubscribed from all subjects")
	return nil
}

// UnsubscribeSubject unsubscribes from a specific subject
func (s *NATSSubscriber) UnsubscribeSubject(ctx context.Context, subject string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var active []*nats.Subscription
	for _, sub := range s.subscriptions {
		if sub.Subject == subject {
			if sub.IsValid() {
				if err := sub.Unsubscribe(); err != nil {
					if err != nats.ErrBadSubscription && err != nats.ErrConnectionClosed {
						s.client.logger.Error("Failed to unsubscribe from subject", zap.String("subject", subject), zap.Error(err))
					}
				}
			}
			s.client.logger.Debug("Unsubscribed from subject", zap.String("subject", subject))
		} else {
			active = append(active, sub)
		}
	}
	s.subscriptions = active
	return nil
}

// SubscribePush subscribes to a JetStream subject with a handler
func (s *NATSSubscriber) SubscribePush(ctx context.Context, subject string, handler HandlerFunc, opts *SubscribeOptions) (*nats.Subscription, error) {
	js, err := s.client.JetStream()
	if err != nil {
		return nil, err
	}

	var subOpts []nats.SubOpt
	if opts != nil {
		subOpts = opts.NatsOptions
		if opts.Durable != "" {
			subOpts = append(subOpts, nats.Durable(opts.Durable))
		}
		if opts.AckWait > 0 {
			subOpts = append(subOpts, nats.AckWait(opts.AckWait))
		}
		if opts.MaxDeliveries > 0 {
			subOpts = append(subOpts, nats.MaxDeliver(opts.MaxDeliveries))
		}
		// Always use explicit ack
		subOpts = append(subOpts, nats.AckExplicit())
	} else {
		subOpts = []nats.SubOpt{nats.AckExplicit()}
	}

	// Create message handler wrapper
	msgHandler := func(msg *nats.Msg) {
		s.wg.Add(1)
		defer s.wg.Done()
		s.processJetStreamMessage(msg, handler, opts)
	}

	sub, err := js.Subscribe(subject, msgHandler, subOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to JetStream: %w", err)
	}

	// Store subscription
	s.mu.Lock()
	s.subscriptions = append(s.subscriptions, sub)
	s.mu.Unlock()

	s.client.logger.Debug("Subscribed to JetStream subject",
		zap.String("subject", subject),
	)

	return sub, nil
}

// SubscribePull subscribes to a JetStream subject using a pull consumer
func (s *NATSSubscriber) SubscribePull(ctx context.Context, subject string, handler HandlerFunc, opts *SubscribeOptions) error {
	js, err := s.client.JetStream()
	if err != nil {
		return err
	}

	if opts == nil || opts.Durable == "" {
		return fmt.Errorf("durable name required for pull subscription")
	}

	// Default options
	batchSize := 10
	fetchTimeout := 5 * time.Second

	if opts.BatchSize > 0 {
		batchSize = opts.BatchSize
	}
	if opts.FetchTimeout > 0 {
		fetchTimeout = opts.FetchTimeout
	}

	// Create pull subscription
	sub, err := js.PullSubscribe(subject, opts.Durable)
	if err != nil {
		return fmt.Errorf("failed to create pull subscription: %w", err)
	}

	// Store subscription
	s.mu.Lock()
	s.subscriptions = append(s.subscriptions, sub)
	s.mu.Unlock()

	s.client.logger.Debug("Created pull subscription",
		zap.String("subject", subject),
		zap.String("durable", opts.Durable),
		zap.Int("batch_size", batchSize),
	)

	// Start background worker
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		for {
			// Check if subscription is valid
			if !sub.IsValid() {
				s.client.logger.Warn("Pull subscription invalid, stopping worker",
					zap.String("subject", subject),
					zap.String("durable", opts.Durable),
				)
				return
			}

			// Fetch batch
			msgs, err := sub.Fetch(batchSize, nats.MaxWait(fetchTimeout))
			if err != nil {
				if err == nats.ErrTimeout {
					// Timeout is normal if no messages, just continue
					// Check for context cancellation or closure if needed, but loop continues
					continue
				}
				if err == nats.ErrConnectionClosed || err == nats.ErrBadSubscription {
					// Stop on terminal errors
					return
				}
				// Log but continue (backoff?)
				s.client.logger.Debug("Failed to fetch messages (retrying)", zap.Error(err))
				time.Sleep(1 * time.Second) // Backoff
				continue
			}

			// Process batch
			for _, msg := range msgs {
				s.processJetStreamMessage(msg, handler, opts)
			}
		}
	}()

	return nil
}

// processJetStreamMessage handles a single JetStream message
func (s *NATSSubscriber) processJetStreamMessage(msg *nats.Msg, handler HandlerFunc, opts *SubscribeOptions) {
	// Unmarshal envelope
	var envelope MessageEnvelope
	if err := json.Unmarshal(msg.Data, &envelope); err != nil {
		s.client.logger.Error("Failed to unmarshal JetStream message",
			zap.Error(err),
			zap.String("subject", msg.Subject),
		)
		// We don't Ack here, so it will be redelivered based on AckWait
		return
	}

	// Capture NATS reply subject
	if msg.Reply != "" {
		envelope.Reply = msg.Reply
	}

	// Apply middleware
	h := handler
	for i := len(s.middleware) - 1; i >= 0; i-- {
		h = s.middleware[i](h)
	}

	// Handle message (tracing is handled in middleware)
	if err := h(context.Background(), msg.Subject, &envelope); err != nil {
		s.client.logger.Error("JetStream handler error",
			zap.Error(err),
			zap.String("subject", msg.Subject),
			zap.String("message_id", envelope.ID),
		)

		// Check for DLQ
		meta, metaErr := msg.Metadata()
		if metaErr == nil && opts != nil && opts.DLQSubject != "" && opts.MaxDeliveries > 0 {
			if int(meta.NumDelivered) >= opts.MaxDeliveries {
				s.client.logger.Warn("Message exceeded max deliveries, moving to DLQ",
					zap.String("subject", msg.Subject),
					zap.String("dlq", opts.DLQSubject),
					zap.String("message_id", envelope.ID),
					zap.Uint64("deliveries", meta.NumDelivered),
				)

				// Publish to DLQ
				// We reuse the client connection to publish.
				// Note: We publish the RAW data to preserve original content.
				if err := s.client.Conn().Publish(opts.DLQSubject, msg.Data); err != nil {
					s.client.logger.Error("Failed to publish to DLQ", zap.Error(err))
					// If DLQ publish fails, we Nak to retry (hopefully transient)
					_ = msg.Nak()
					return
				}

				// Ack original message so it doesn't redeliver
				if err := msg.Ack(); err != nil {
					s.client.logger.Error("Failed to ack message after DLQ move", zap.Error(err))
				}
				return
			}
		}

		// Explicitly Nak to trigger redelivery
		if err := msg.Nak(); err != nil {
			s.client.logger.Error("Failed to nak JetStream message", zap.Error(err))
		}
		return
	}

	// Acknowledge message
	if err := msg.Ack(); err != nil {
		s.client.logger.Error("Failed to ack JetStream message",
			zap.Error(err),
			zap.String("subject", msg.Subject),
			zap.String("message_id", envelope.ID),
		)
	}
}

// Close closes the subscriber and unsubscribes from all subjects
func (s *NATSSubscriber) Close() error {
	if err := s.Unsubscribe(context.Background()); err != nil {
		return err
	}

	// Wait for active handlers
	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		s.client.logger.Debug("Subscriber closed gracefully")
	case <-time.After(s.drainTimeout):
		s.client.logger.Warn("Subscriber closed with active handlers (timeout)",
			zap.Duration("timeout", s.drainTimeout))
	}

	return nil
}
