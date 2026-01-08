package nats

import (
	"context"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

// --- Tracing Middleware ---

// metadataCarrier implements propagation.TextMapCarrier for MessageEnvelope.Metadata
type metadataCarrier map[string]string

func (c metadataCarrier) Get(key string) string {
	return c[key]
}

func (c metadataCarrier) Set(key string, value string) {
	c[key] = value
}

func (c metadataCarrier) Keys() []string {
	keys := make([]string, 0, len(c))
	for k := range c {
		keys = append(keys, k)
	}
	return keys
}

// TracingMiddleware returns a middleware that extracts trace context from message metadata
func TracingMiddleware(tracer trace.Tracer) SubscriberMiddleware {
	propagator := otel.GetTextMapPropagator()

	return func(next HandlerFunc) HandlerFunc {
		return func(ctx context.Context, subject string, env *MessageEnvelope) error {
			// Extract context from metadata
			if env.Metadata == nil {
				env.Metadata = make(map[string]string)
			}

			ctx = propagator.Extract(ctx, metadataCarrier(env.Metadata))

			// Start span
			// Use <destination> <operation> naming convention
			ctx, span := tracer.Start(ctx, fmt.Sprintf("%s process", subject),
				trace.WithSpanKind(trace.SpanKindConsumer),
				trace.WithAttributes(
					semconv.MessagingSystem(systemName),
					semconv.MessagingDestinationName(subject),
					semconv.MessagingOperationProcess,
					semconv.MessagingMessageID(env.ID), // Optional but good for correlation
					attribute.String("messaging.source", env.Source),
					attribute.String("messaging.message_type", env.Type),
				),
			)
			defer span.End()

			err := next(ctx, subject, env)
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
			} else {
				span.SetStatus(codes.Ok, "OK")
			}

			return err
		}
	}
}

// PublisherTracingMiddleware returns a middleware that injects trace context into message metadata
func PublisherTracingMiddleware(tracer trace.Tracer) PublisherMiddleware {
	return func(next PublisherFunc) PublisherFunc {
		return func(ctx context.Context, subject string, msgType string, data interface{}, opts *PublishOptions) error {
			// Start span
			ctx, span := tracer.Start(ctx, fmt.Sprintf("%s send", subject),
				trace.WithSpanKind(trace.SpanKindProducer),
				trace.WithAttributes(
					semconv.MessagingSystem(systemName),
					semconv.MessagingDestinationName(subject),
					// Use generic publish operation or custom if valid constant not found
					attribute.String("messaging.operation", "publish"),
					attribute.String("messaging.message_type", msgType),
				),
			)
			defer span.End()

			// Inject context into metadata for propagation
			if opts == nil {
				opts = &PublishOptions{}
			}
			if opts.Metadata == nil {
				opts.Metadata = make(map[string]string)
			}
			otel.GetTextMapPropagator().Inject(ctx, metadataCarrier(opts.Metadata))

			err := next(ctx, subject, msgType, data, opts)
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
			} else {
				span.SetStatus(codes.Ok, "OK")
			}

			return err
		}
	}
}

// JSPublisherTracingMiddleware returns a middleware that injects trace context into JetStream message metadata
func JSPublisherTracingMiddleware(tracer trace.Tracer) JSPublisherMiddleware {
	return func(next JSPublisherFunc) JSPublisherFunc {
		return func(ctx context.Context, subject string, msgType string, data interface{}, opts *PublishOptions) (*nats.PubAck, error) {
			// Start span
			ctx, span := tracer.Start(ctx, fmt.Sprintf("%s send", subject),
				trace.WithSpanKind(trace.SpanKindProducer),
				trace.WithAttributes(
					semconv.MessagingSystem(systemName),
					semconv.MessagingDestinationName(subject),
					attribute.String("messaging.operation", "publish_js"),
					attribute.String("messaging.message_type", msgType),
				),
			)
			defer span.End()

			// Inject context into metadata for propagation
			if opts == nil {
				opts = &PublishOptions{}
			}
			if opts.Metadata == nil {
				opts.Metadata = make(map[string]string)
			}
			otel.GetTextMapPropagator().Inject(ctx, metadataCarrier(opts.Metadata))

			ack, err := next(ctx, subject, msgType, data, opts)
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
			} else {
				span.SetStatus(codes.Ok, "OK")
				span.SetAttributes(attribute.Int64("messaging.js.sequence", int64(ack.Sequence)))
			}

			return ack, err
		}
	}
}

// TracingPubAckFuture wraps nats.PubAckFuture to end the span when the future resolves
type TracingPubAckFuture struct {
	nats.PubAckFuture
	span       trace.Span
	proxyAckCh chan *nats.PubAck
	proxyErrCh chan error
}

// According to NATS docs, either Ok() OR Err() will receive, never both
// But we need to handle the case where consumer never reads from channels

func newTracingPubAckFuture(future nats.PubAckFuture, span trace.Span) *TracingPubAckFuture {
	t := &TracingPubAckFuture{
		PubAckFuture: future,
		span:         span,
		proxyAckCh:   make(chan *nats.PubAck, 1), // Buffered!
		proxyErrCh:   make(chan error, 1),        // Buffered!
	}
	go t.monitor()
	return t
}

func (f *TracingPubAckFuture) monitor() {
	defer f.span.End()

	select {
	case ack, ok := <-f.PubAckFuture.Ok():
		if ok {
			f.span.SetStatus(codes.Ok, "OK")
			f.span.SetAttributes(
				attribute.Int64("messaging.js.sequence", int64(ack.Sequence)),
				attribute.String("messaging.js.stream", ack.Stream),
				attribute.Bool("messaging.js.duplicate", ack.Duplicate),
			)
			f.proxyAckCh <- ack
		}
		close(f.proxyAckCh)
		close(f.proxyErrCh)

	case err, ok := <-f.PubAckFuture.Err():
		if ok && err != nil {
			f.span.RecordError(err)
			f.span.SetStatus(codes.Error, err.Error())
			f.proxyErrCh <- err
		}
		close(f.proxyAckCh)
		close(f.proxyErrCh)
	}
}

func (f *TracingPubAckFuture) Ok() <-chan *nats.PubAck {
	return f.proxyAckCh
}

func (f *TracingPubAckFuture) Err() <-chan error {
	return f.proxyErrCh
}

// JSAsyncPublisherTracingMiddleware returns a middleware that injects trace context into async JetStream message metadata
func JSAsyncPublisherTracingMiddleware(tracer trace.Tracer) JSAsyncPublisherMiddleware {
	return func(next JSAsyncPublisherFunc) JSAsyncPublisherFunc {
		return func(ctx context.Context, subject string, msgType string, data interface{}, opts *PublishOptions) (nats.PubAckFuture, error) {
			// Start span
			ctx, span := tracer.Start(ctx, fmt.Sprintf("%s send", subject),
				trace.WithSpanKind(trace.SpanKindProducer),
				trace.WithAttributes(
					semconv.MessagingSystem(systemName),
					semconv.MessagingDestinationName(subject),
					attribute.String("messaging.operation", "publish_js_async"),
					attribute.String("messaging.message_type", msgType),
				),
			)
			// Do NOT defer span.End() here; it's handled in the future wrapper

			// Inject context into metadata for propagation
			if opts != nil {
				if opts.Metadata == nil {
					opts.Metadata = make(map[string]string)
				}
				otel.GetTextMapPropagator().Inject(ctx, metadataCarrier(opts.Metadata))
			}

			future, err := next(ctx, subject, msgType, data, opts)
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
				span.End() // End immediately if submission failed
				return future, err
			}

			return newTracingPubAckFuture(future, span), nil
		}
	}
}

// RequestTracingMiddleware returns a middleware that injects trace context into request metadata
func RequestTracingMiddleware(tracer trace.Tracer) RequestMiddleware {
	return func(next RequestFunc) RequestFunc {
		return func(ctx context.Context, subject string, msgType string, data interface{}, timeout time.Duration, opts *PublishOptions) (*MessageEnvelope, error) {
			// Start span
			ctx, span := tracer.Start(ctx, fmt.Sprintf("%s request", subject),
				trace.WithSpanKind(trace.SpanKindProducer),
				trace.WithAttributes(
					semconv.MessagingSystem(systemName),
					semconv.MessagingDestinationName(subject),
					attribute.String("messaging.operation", "request"),
					attribute.String("messaging.message_type", msgType),
				),
			)
			defer span.End()

			// Inject context into metadata for propagation
			if opts == nil {
				opts = &PublishOptions{}
			}
			if opts.Metadata == nil {
				opts.Metadata = make(map[string]string)
			}
			otel.GetTextMapPropagator().Inject(ctx, metadataCarrier(opts.Metadata))

			resp, err := next(ctx, subject, msgType, data, timeout, opts)
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
			} else {
				span.SetStatus(codes.Ok, "OK")
				span.SetAttributes(
					attribute.String("messaging.response_id", resp.ID),
					attribute.String("messaging.response_type", resp.Type),
				)
			}

			return resp, err
		}
	}
}

// Note: To fully support trace propagation, we should update MessagePublisher interface
// and Publisher implementation to accept metadata or a context that can be used to
// populate the envelope's metadata.
