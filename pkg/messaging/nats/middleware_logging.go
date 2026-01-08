package nats

import (
	"context"
	"time"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
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
				zap.Int("payload_size", len(env.Data)),
				zap.Time("timestamp", env.Timestamp),
			}

			if env.Metadata != nil {
				if correlationID, ok := env.Metadata["correlation_id"]; ok {
					fields = append(fields, zap.String("correlation_id", correlationID))
				}
			}

			if err != nil {
				logger.Error("Message processing failed", append(fields, zap.Error(err))...)
			} else {
				logger.Debug("Message processed successfully", fields...)
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

			if opts != nil && opts.Metadata != nil {
				if correlationID, ok := opts.Metadata["correlation_id"]; ok {
					fields = append(fields, zap.String("correlation_id", correlationID))
				}
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

// JSPublisherLoggingMiddleware returns a middleware that logs JetStream message publishing
func JSPublisherLoggingMiddleware(logger *zap.Logger) JSPublisherMiddleware {
	return func(next JSPublisherFunc) JSPublisherFunc {
		return func(ctx context.Context, subject string, msgType string, data interface{}, opts *PublishOptions) (*nats.PubAck, error) {
			start := time.Now()
			ack, err := next(ctx, subject, msgType, data, opts)
			duration := time.Since(start)

			fields := []zap.Field{
				zap.String("subject", subject),
				zap.String("type", msgType),
				zap.Duration("duration", duration),
			}

			if opts != nil && opts.Metadata != nil {
				if correlationID, ok := opts.Metadata["correlation_id"]; ok {
					fields = append(fields, zap.String("correlation_id", correlationID))
				}
			}

			if err != nil {
				logger.Error("JetStream message publishing failed", append(fields, zap.Error(err))...)
			} else {
				logger.Debug("JetStream message published successfully",
					append(fields, zap.Uint64("stream_seq", ack.Sequence))...)
			}

			return ack, err
		}
	}
}

// JSAsyncPublisherLoggingMiddleware returns a middleware that logs async JetStream message publishing
func JSAsyncPublisherLoggingMiddleware(logger *zap.Logger) JSAsyncPublisherMiddleware {
	return func(next JSAsyncPublisherFunc) JSAsyncPublisherFunc {
		return func(ctx context.Context, subject string, msgType string, data interface{}, opts *PublishOptions) (nats.PubAckFuture, error) {
			start := time.Now()
			future, err := next(ctx, subject, msgType, data, opts)
			duration := time.Since(start)

			fields := []zap.Field{
				zap.String("subject", subject),
				zap.String("type", msgType),
				zap.Duration("duration", duration),
			}

			if opts != nil && opts.Metadata != nil {
				if correlationID, ok := opts.Metadata["correlation_id"]; ok {
					fields = append(fields, zap.String("correlation_id", correlationID))
				}
			}

			if err != nil {
				logger.Error("Async JetStream message publishing failed", append(fields, zap.Error(err))...)
			} else {
				logger.Debug("Async JetStream message publish initiated", fields...)
			}

			return future, err
		}
	}
}

// RequestLoggingMiddleware returns a middleware that logs request-reply interactions
func RequestLoggingMiddleware(logger *zap.Logger) RequestMiddleware {
	return func(next RequestFunc) RequestFunc {
		return func(ctx context.Context, subject string, msgType string, data interface{}, timeout time.Duration, opts *PublishOptions) (*MessageEnvelope, error) {
			start := time.Now()
			resp, err := next(ctx, subject, msgType, data, timeout, opts)
			duration := time.Since(start)

			fields := []zap.Field{
				zap.String("subject", subject),
				zap.String("type", msgType),
				zap.Duration("duration", duration),
			}

			if opts != nil && opts.Metadata != nil {
				if correlationID, ok := opts.Metadata["correlation_id"]; ok {
					fields = append(fields, zap.String("correlation_id", correlationID))
				}
			}

			if err != nil {
				logger.Error("Request failed", append(fields, zap.Error(err))...)
			} else {
				logger.Debug("Request completed successfully",
					append(fields,
						zap.String("response_id", resp.ID),
						zap.String("response_type", resp.Type),
					)...,
				)
			}

			return resp, err
		}
	}
}
