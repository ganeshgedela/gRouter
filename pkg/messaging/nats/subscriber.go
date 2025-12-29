package nats

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
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
}

// NewSubscriber creates a new subscriber
func NewSubscriber(client *Client, source string) Subscriber {
	return &NATSSubscriber{
		client:        client,
		source:        source,
		subscriptions: make([]*nats.Subscription, 0),
		middleware:    make([]SubscriberMiddleware, 0),
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
func (s *NATSSubscriber) Subscribe(subject string, handler HandlerFunc, opts *SubscribeOptions) error {
	if !s.client.IsConnected() {
		return fmt.Errorf("not connected to NATS")
	}

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

		// Extract trace context
		ctx := otel.GetTextMapPropagator().Extract(context.Background(), metadataCarrier(envelope.Metadata))

		// Start Span
		ctx, span := tracer.Start(ctx, spanNameProcess+" "+msg.Subject,
			trace.WithSpanKind(trace.SpanKindConsumer),
			trace.WithAttributes(
				semconv.MessagingSystem(systemName),
				semconv.MessagingDestinationName(msg.Subject),
				semconv.MessagingOperationProcess,
				semconv.MessagingMessageID(envelope.ID),
			),
		)
		defer span.End()

		// ✅ capture NATS reply subject for request-reply
		if msg.Reply != "" {
			envelope.Reply = msg.Reply
		}

		// Validate data if validator is set
		if s.validator != nil {
			if err := s.validator.Validate(envelope.Type, envelope.Data); err != nil {
				s.client.logger.Error("Validation failed",
					zap.Error(err),
					zap.String("subject", msg.Subject),
					zap.String("type", envelope.Type),
					zap.String("id", envelope.ID),
				)
				return
			}
		}

		s.client.logger.Debug("Received message",
			zap.String("subject", msg.Subject),
			zap.String("type", envelope.Type),
			zap.String("id", envelope.ID),
			zap.String("reply", envelope.Reply),
		)

		// Apply middleware
		h := handler
		for i := len(s.middleware) - 1; i >= 0; i-- {
			h = s.middleware[i](h)
		}

		// Handle message
		if err := h(ctx, msg.Subject, &envelope); err != nil {
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
		return fmt.Errorf("failed to subscribe: %w", err)
	}

	// Store subscription
	s.mu.Lock()
	s.subscriptions = append(s.subscriptions, sub)
	s.mu.Unlock()

	s.client.logger.Info("Subscribed to subject",
		zap.String("subject", subject),
		zap.String("queue_group", func() string {
			if opts != nil {
				return opts.QueueGroup
			}
			return ""
		}()),
	)

	return nil
}

// Unsubscribe unsubscribes from all subscriptions
func (s *NATSSubscriber) Unsubscribe() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, sub := range s.subscriptions {
		if err := sub.Unsubscribe(); err != nil {
			s.client.logger.Error("Failed to unsubscribe", zap.Error(err))
		}
	}

	s.subscriptions = make([]*nats.Subscription, 0)
	s.client.logger.Info("Unsubscribed from all subjects")
	return nil
}

// SubscribePush subscribes to a JetStream subject with a handler
func (s *NATSSubscriber) SubscribePush(subject string, handler HandlerFunc, opts ...nats.SubOpt) error {
	js, err := s.client.JetStream()
	if err != nil {
		return err
	}

	// Create message handler wrapper
	msgHandler := func(msg *nats.Msg) {
		s.wg.Add(1)
		defer s.wg.Done()

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

		// Extract trace context
		ctx := otel.GetTextMapPropagator().Extract(context.Background(), metadataCarrier(envelope.Metadata))

		// Start Span
		ctx, span := tracer.Start(ctx, spanNameProcess+" "+msg.Subject,
			trace.WithSpanKind(trace.SpanKindConsumer),
			trace.WithAttributes(
				semconv.MessagingSystem(systemName),
				semconv.MessagingDestinationName(msg.Subject),
				semconv.MessagingOperationProcess,
				semconv.MessagingMessageID(envelope.ID),
			),
		)
		defer span.End()

		// Capture NATS reply subject
		if msg.Reply != "" {
			envelope.Reply = msg.Reply
		}

		// Validate data if validator is set
		if s.validator != nil {
			if err := s.validator.Validate(envelope.Type, envelope.Data); err != nil {
				s.client.logger.Error("JetStream validation failed",
					zap.Error(err),
					zap.String("subject", msg.Subject),
					zap.String("type", envelope.Type),
					zap.String("id", envelope.ID),
				)
				// We don't Ack here, so it will be redelivered or go to DLQ
				return
			}
		}

		s.client.logger.Debug("Received JetStream message",
			zap.String("subject", msg.Subject),
			zap.String("type", envelope.Type),
			zap.String("id", envelope.ID),
		)

		// Apply middleware
		h := handler
		for i := len(s.middleware) - 1; i >= 0; i-- {
			h = s.middleware[i](h)
		}

		// Handle message
		if err := h(ctx, msg.Subject, &envelope); err != nil {
			s.client.logger.Error("JetStream handler error",
				zap.Error(err),
				zap.String("subject", msg.Subject),
				zap.String("message_id", envelope.ID),
			)
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

	sub, err := js.Subscribe(subject, msgHandler, opts...)
	if err != nil {
		return fmt.Errorf("failed to subscribe to JetStream: %w", err)
	}

	// Store subscription
	s.mu.Lock()
	s.subscriptions = append(s.subscriptions, sub)
	s.mu.Unlock()

	s.client.logger.Info("Subscribed to JetStream subject",
		zap.String("subject", subject),
	)

	return nil
}

// SubscribePull subscribes to a JetStream subject using a pull consumer
func (s *NATSSubscriber) SubscribePull(subject, durable string, handler HandlerFunc, opts ...PullOption) error {
	js, err := s.client.JetStream()
	if err != nil {
		return err
	}

	// Default options
	options := &PullOptions{
		BatchSize:    10,
		FetchTimeout: 5 * time.Second,
	}

	// Apply options
	for _, opt := range opts {
		opt(options)
	}

	// Create pull subscription
	sub, err := js.PullSubscribe(subject, durable)
	if err != nil {
		return fmt.Errorf("failed to create pull subscription: %w", err)
	}

	// Store subscription
	s.mu.Lock()
	s.subscriptions = append(s.subscriptions, sub)
	s.mu.Unlock()

	s.client.logger.Info("Created pull subscription",
		zap.String("subject", subject),
		zap.String("durable", durable),
		zap.Int("batch_size", options.BatchSize),
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
					zap.String("durable", durable),
				)
				return
			}

			// Fetch batch
			msgs, err := sub.Fetch(options.BatchSize, nats.MaxWait(options.FetchTimeout))
			if err != nil {
				if err == nats.ErrTimeout {
					// Timeout is normal if no messages, just continue
					continue
				}
				if err == nats.ErrConnectionClosed || err == nats.ErrBadSubscription {
					// Stop on terminal errors
					return
				}
				s.client.logger.Error("Failed to fetch messages", zap.Error(err))
				time.Sleep(1 * time.Second) // Backoff
				continue
			}

			// Process batch
			for _, msg := range msgs {
				s.processJetStreamMessage(msg, handler)
			}
		}
	}()

	return nil
}

// processJetStreamMessage handles a single JetStream message
func (s *NATSSubscriber) processJetStreamMessage(msg *nats.Msg, handler HandlerFunc) {
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

	// Extract trace context
	ctx := otel.GetTextMapPropagator().Extract(context.Background(), metadataCarrier(envelope.Metadata))

	// Start Span
	ctx, span := tracer.Start(ctx, spanNameProcess+" "+msg.Subject,
		trace.WithSpanKind(trace.SpanKindConsumer),
		trace.WithAttributes(
			semconv.MessagingSystem(systemName),
			semconv.MessagingDestinationName(msg.Subject),
			semconv.MessagingOperationProcess,
			semconv.MessagingMessageID(envelope.ID),
		),
	)
	defer span.End()

	// Capture NATS reply subject
	if msg.Reply != "" {
		envelope.Reply = msg.Reply
	}

	// Validate data if validator is set
	if s.validator != nil {
		if err := s.validator.Validate(envelope.Type, envelope.Data); err != nil {
			s.client.logger.Error("JetStream validation failed",
				zap.Error(err),
				zap.String("subject", msg.Subject),
				zap.String("type", envelope.Type),
				zap.String("id", envelope.ID),
			)
			// We don't Ack here, so it will be redelivered or go to DLQ
			return
		}
	}

	s.client.logger.Debug("Received JetStream message",
		zap.String("subject", msg.Subject),
		zap.String("type", envelope.Type),
		zap.String("id", envelope.ID),
	)

	// Apply middleware
	h := handler
	for i := len(s.middleware) - 1; i >= 0; i-- {
		h = s.middleware[i](h)
	}

	// Handle message
	if err := h(ctx, msg.Subject, &envelope); err != nil {
		s.client.logger.Error("JetStream handler error",
			zap.Error(err),
			zap.String("subject", msg.Subject),
			zap.String("message_id", envelope.ID),
		)
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
	if err := s.Unsubscribe(); err != nil {
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
		s.client.logger.Info("Subscriber closed gracefully")
	case <-time.After(5 * time.Second):
		s.client.logger.Warn("Subscriber closed with active handlers (timeout)")
	}

	return nil
}
