package nats

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// topic = <service_manager_identity>.<service>.<operation>

// NATSPublisher handles message publishing
type NATSPublisher struct {
	client     *Client
	source     string
	validator  Validator
	middleware []PublisherMiddleware
}

// NewPublisher creates a new publisher
func NewPublisher(client *Client, source string) Publisher {
	return &NATSPublisher{
		client:     client,
		source:     source,
		middleware: make([]PublisherMiddleware, 0),
	}
}

// Use adds middleware to the publisher
func (p *NATSPublisher) Use(mw ...PublisherMiddleware) {
	p.middleware = append(p.middleware, mw...)
}

// SetValidator sets the validator for the publisher
func (p *NATSPublisher) SetValidator(v Validator) {
	p.validator = v
}

// Publish publishes a message to a subject
func (p *NATSPublisher) Publish(ctx context.Context, subject string, msgType string, data interface{}, opts *PublishOptions) error {
	publishFunc := p.publish

	// Apply middleware in reverse order
	for i := len(p.middleware) - 1; i >= 0; i-- {
		publishFunc = p.middleware[i](publishFunc)
	}

	return publishFunc(ctx, subject, msgType, data, opts)
}

func (p *NATSPublisher) publish(ctx context.Context, subject string, msgType string, data interface{}, opts *PublishOptions) error {
	// Marshal data
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	// Validate data if validator is set
	if p.validator != nil {
		if err := p.validator.Validate(msgType, dataBytes); err != nil {
			return fmt.Errorf("validation failed for type %s: %w", msgType, err)
		}
	}

	if !p.client.IsConnected() {
		return fmt.Errorf("not connected to NATS")
	}

	// Create envelope
	envelope := MessageEnvelope{
		ID:        uuid.New().String(),
		Type:      msgType,
		Timestamp: time.Now(),
		Source:    p.source,
		Data:      dataBytes,
		Metadata:  make(map[string]string),
	}

	// Inject trace context into metadata
	otel.GetTextMapPropagator().Inject(ctx, metadataCarrier(envelope.Metadata))

	// Marshal envelope
	envelopeBytes, err := json.Marshal(envelope)
	if err != nil {
		return fmt.Errorf("failed to marshal envelope: %w", err)
	}

	// Publish
	if opts != nil && opts.Async {
		// Async publish
		if err := p.client.Conn().Publish(subject, envelopeBytes); err != nil {
			return fmt.Errorf("failed to publish message: %w", err)
		}
	} else {
		// Sync publish with flush
		if err := p.client.Conn().Publish(subject, envelopeBytes); err != nil {
			return fmt.Errorf("failed to publish message: %w", err)
		}
		if err := p.client.Conn().Flush(); err != nil {
			return fmt.Errorf("failed to flush: %w", err)
		}
	}

	p.client.logger.Debug("Published message",
		zap.String("subject", subject),
		zap.String("type", msgType),
		zap.String("id", envelope.ID),
	)

	return nil
}

// PublishError publishes an error message to a reply subject
func (p *NATSPublisher) PublishError(ctx context.Context, subject string, errMsg string) error {
	if subject == "" {
		return nil
	}

	errorData := map[string]string{"error": errMsg}
	// Error messages should always be synchronous to ensure delivery before we close context or connection
	return p.Publish(ctx, subject, "error", errorData, &PublishOptions{Async: false})
}

// Request sends a request and waits for a response
func (p *NATSPublisher) Request(ctx context.Context, subject string, msgType string, data interface{}, timeout time.Duration) (*MessageEnvelope, error) {
	if !p.client.IsConnected() {
		return nil, fmt.Errorf("not connected to NATS")
	}

	// Marshal data
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	// Create envelope
	envelope := MessageEnvelope{
		ID:        uuid.New().String(),
		Type:      msgType,
		Timestamp: time.Now(),
		Source:    p.source,
		Data:      dataBytes,
		Metadata:  make(map[string]string),
	}

	// Inject trace context into metadata
	otel.GetTextMapPropagator().Inject(ctx, metadataCarrier(envelope.Metadata))

	// Marshal envelope
	envelopeBytes, err := json.Marshal(envelope)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal envelope: %w", err)
	}

	// Send request with context support
	msg, err := p.client.Conn().RequestWithContext(ctx, subject, envelopeBytes)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	// Unmarshal response
	var response MessageEnvelope
	if err := json.Unmarshal(msg.Data, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	p.client.logger.Debug("Request completed",
		zap.String("subject", subject),
		zap.String("request_id", envelope.ID),
		zap.String("response_id", response.ID),
	)

	return &response, nil
}

// PublishJS publishes a message to a JetStream subject
func (p *NATSPublisher) PublishJS(ctx context.Context, subject string, msgType string, data interface{}, opts ...nats.PubOpt) (*nats.PubAck, error) {
	// Marshal data
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	// Validate data if validator is set
	if p.validator != nil {
		if err := p.validator.Validate(msgType, dataBytes); err != nil {
			return nil, fmt.Errorf("validation failed for type %s: %w", msgType, err)
		}
	}

	js, err := p.client.JetStream()
	if err != nil {
		return nil, err
	}

	// Create envelope
	envelope := MessageEnvelope{
		ID:        uuid.New().String(),
		Type:      msgType,
		Timestamp: time.Now(),
		Source:    p.source,
		Data:      dataBytes,
		Metadata:  make(map[string]string),
	}

	// Inject trace context into metadata
	otel.GetTextMapPropagator().Inject(ctx, metadataCarrier(envelope.Metadata))

	// Marshal envelope
	envelopeBytes, err := json.Marshal(envelope)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal envelope: %w", err)
	}

	// Publish to JetStream with context
	ack, err := js.PublishMsg(&nats.Msg{
		Subject: subject,
		Data:    envelopeBytes,
	}, append(opts, nats.Context(ctx))...)
	if err != nil {
		return nil, fmt.Errorf("failed to publish to JetStream: %w", err)
	}

	p.client.logger.Debug("Published JetStream message",
		zap.String("subject", subject),
		zap.String("type", msgType),
		zap.String("id", envelope.ID),
		zap.Uint64("stream_seq", ack.Sequence),
	)

	return ack, nil
}

// PublishAsyncJS publishes a message to a JetStream subject asynchronously
func (p *NATSPublisher) PublishAsyncJS(ctx context.Context, subject string, msgType string, data interface{}, opts ...nats.PubOpt) (nats.PubAckFuture, error) {
	// Marshal data
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	// Validate data if validator is set
	if p.validator != nil {
		if err := p.validator.Validate(msgType, dataBytes); err != nil {
			return nil, fmt.Errorf("validation failed for type %s: %w", msgType, err)
		}
	}

	js, err := p.client.JetStream()
	if err != nil {
		return nil, err
	}

	// Start Span
	ctx, span := tracer.Start(ctx, spanNamePublish+" "+subject,
		trace.WithSpanKind(trace.SpanKindProducer),
		trace.WithAttributes(
			semconv.MessagingSystem(systemName),
			semconv.MessagingDestinationName(subject),
			semconv.MessagingOperationPublish,
		),
	)
	defer span.End()

	// Create envelope
	envelope := MessageEnvelope{
		ID:        uuid.New().String(),
		Type:      msgType,
		Timestamp: time.Now(),
		Source:    p.source,
		Data:      dataBytes,
		Metadata:  make(map[string]string),
	}

	// Inject trace context into metadata
	otel.GetTextMapPropagator().Inject(ctx, metadataCarrier(envelope.Metadata))

	// Marshal envelope
	envelopeBytes, err := json.Marshal(envelope)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal envelope: %w", err)
	}

	// Publish to JetStream asynchronously
	future, err := js.PublishAsync(subject, envelopeBytes, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to publish async to JetStream: %w", err)
	}

	p.client.logger.Debug("Published JetStream message asynchronously",
		zap.String("subject", subject),
		zap.String("type", msgType),
		zap.String("id", envelope.ID),
	)

	return future, nil
}
