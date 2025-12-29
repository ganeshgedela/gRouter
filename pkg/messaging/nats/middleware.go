package nats

import (
	"context"
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

var (
	// Metrics for publishers
	publishCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "messaging_publish_total",
		Help: "Total number of messages published",
	}, []string{"subject", "type", "status"})

	publishDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "messaging_publish_duration_seconds",
		Help:    "Duration of message publishing in seconds",
		Buckets: prometheus.DefBuckets,
	}, []string{"subject", "type"})

	// Metrics for subscribers
	subscribeCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "messaging_subscribe_total",
		Help: "Total number of messages received",
	}, []string{"subject", "type", "status"})

	subscribeDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "messaging_subscribe_duration_seconds",
		Help:    "Duration of message processing in seconds",
		Buckets: prometheus.DefBuckets,
	}, []string{"subject", "type"})
)

// --- Logging Middleware ---

// LoggingMiddleware returns a middleware that logs message processing
func LoggingMiddleware(logger *zap.Logger) SubscriberMiddleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx context.Context, subject string, env *MessageEnvelope) error {
			start := time.Now()
			err := next(ctx, subject, env)
			duration := time.Since(start)

			fields := []zap.Field{
				zap.String("subject", subject),
				zap.String("type", env.Type),
				zap.String("id", env.ID),
				zap.String("source", env.Source),
				zap.Duration("duration", duration),
			}

			if err != nil {
				logger.Error("Message processing failed", append(fields, zap.Error(err))...)
			} else {
				logger.Info("Message processed successfully", fields...)
			}

			return err
		}
	}
}

// PublisherLoggingMiddleware returns a middleware that logs message publishing
func PublisherLoggingMiddleware(logger *zap.Logger) PublisherMiddleware {
	return func(next PublisherFunc) PublisherFunc {
		return func(ctx context.Context, subject string, msgType string, data interface{}, opts *PublishOptions) error {
			start := time.Now()
			err := next(ctx, subject, msgType, data, opts)
			duration := time.Since(start)

			fields := []zap.Field{
				zap.String("subject", subject),
				zap.String("type", msgType),
				zap.Duration("duration", duration),
			}

			if err != nil {
				logger.Error("Message publishing failed", append(fields, zap.Error(err))...)
			} else {
				logger.Debug("Message published successfully", fields...)
			}

			return err
		}
	}
}

// --- Metrics Middleware ---

// MetricsMiddleware returns a middleware that tracks message processing metrics
func MetricsMiddleware() SubscriberMiddleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx context.Context, subject string, env *MessageEnvelope) error {
			start := time.Now()
			err := next(ctx, subject, env)
			duration := time.Since(start)

			status := "success"
			if err != nil {
				status = "error"
			}

			subscribeCounter.WithLabelValues(subject, env.Type, status).Inc()
			subscribeDuration.WithLabelValues(subject, env.Type).Observe(duration.Seconds())

			return err
		}
	}
}

// PublisherMetricsMiddleware returns a middleware that tracks message publishing metrics
func PublisherMetricsMiddleware() PublisherMiddleware {
	return func(next PublisherFunc) PublisherFunc {
		return func(ctx context.Context, subject string, msgType string, data interface{}, opts *PublishOptions) error {
			start := time.Now()
			err := next(ctx, subject, msgType, data, opts)
			duration := time.Since(start)

			status := "success"
			if err != nil {
				status = "error"
			}

			publishCounter.WithLabelValues(subject, msgType, status).Inc()
			publishDuration.WithLabelValues(subject, msgType).Observe(duration.Seconds())

			return err
		}
	}
}

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
			ctx, span := tracer.Start(ctx, fmt.Sprintf("messaging.receive %s", subject),
				trace.WithSpanKind(trace.SpanKindConsumer),
				trace.WithAttributes(
					attribute.String("messaging.subject", subject),
					attribute.String("messaging.message_id", env.ID),
					attribute.String("messaging.message_type", env.Type),
					attribute.String("messaging.source", env.Source),
				),
			)
			defer span.End()

			err := next(ctx, subject, env)
			if err != nil {
				span.RecordError(err)
				span.SetAttributes(attribute.String("error", err.Error()))
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
			ctx, span := tracer.Start(ctx, fmt.Sprintf("messaging.send %s", subject),
				trace.WithSpanKind(trace.SpanKindProducer),
				trace.WithAttributes(
					attribute.String("messaging.subject", subject),
					attribute.String("messaging.message_type", msgType),
				),
			)
			defer span.End()

			// We need to inject context into the envelope metadata.
			// However, PublisherFunc doesn't give us access to the envelope directly.
			// The envelope is created inside the publisher.publish method.
			// This is a design limitation. To fix this, we would need to refactor
			// the publisher to allow middleware to modify the envelope or pass metadata.

			// For now, we'll just wrap the call. Tracing will work for the local process,
			// but propagation to the subscriber will require a refactor of the Publisher.

			err := next(ctx, subject, msgType, data, opts)
			if err != nil {
				span.RecordError(err)
				span.SetAttributes(attribute.String("error", err.Error()))
			}

			return err
		}
	}
}

// Note: To fully support trace propagation, we should update MessagePublisher interface
// and Publisher implementation to accept metadata or a context that can be used to
// populate the envelope's metadata.
