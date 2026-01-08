package nats

import (
	"context"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

// --- Recovery Middleware ---

// RecoveryMiddleware returns a middleware that recovers from panics in subscribers
func RecoveryMiddleware(logger *zap.Logger) SubscriberMiddleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx context.Context, subject string, env *MessageEnvelope) (err error) {
			defer func() {
				if r := recover(); r != nil {
					// Log the panic with stack trace
					logger.Error("Panic recovered in subscriber",
						zap.Any("panic", r),
						zap.String("subject", subject),
						zap.String("stack", string(debug.Stack())),
					)
					// Return an error to indicate failure
					err = fmt.Errorf("panic recovered: %v", r)
				}
			}()
			return next(ctx, subject, env)
		}
	}
}

// PublisherRecoveryMiddleware returns a middleware that recovers from panics in publishers
func PublisherRecoveryMiddleware(logger *zap.Logger) PublisherMiddleware {
	return func(next PublisherFunc) PublisherFunc {
		return func(ctx context.Context, subject string, msgType string, data interface{}, opts *PublishOptions) (err error) {
			defer func() {
				if r := recover(); r != nil {
					logger.Error("Panic recovered in publisher",
						zap.Any("panic", r),
						zap.String("subject", subject),
						zap.String("stack", string(debug.Stack())),
					)
					err = fmt.Errorf("panic recovered: %v", r)
				}
			}()
			return next(ctx, subject, msgType, data, opts)
		}
	}
}

// JSPublisherRecoveryMiddleware returns a middleware that recovers from panics in JetStream publishers
func JSPublisherRecoveryMiddleware(logger *zap.Logger) JSPublisherMiddleware {
	return func(next JSPublisherFunc) JSPublisherFunc {
		return func(ctx context.Context, subject string, msgType string, data interface{}, opts *PublishOptions) (ack *nats.PubAck, err error) {
			defer func() {
				if r := recover(); r != nil {
					logger.Error("Panic recovered in JetStream publisher",
						zap.Any("panic", r),
						zap.String("subject", subject),
						zap.String("stack", string(debug.Stack())),
					)
					err = fmt.Errorf("panic recovered: %v", r)
				}
			}()
			return next(ctx, subject, msgType, data, opts)
		}
	}
}

// JSAsyncPublisherRecoveryMiddleware returns a middleware that recovers from panics in async JetStream publishers
func JSAsyncPublisherRecoveryMiddleware(logger *zap.Logger) JSAsyncPublisherMiddleware {
	return func(next JSAsyncPublisherFunc) JSAsyncPublisherFunc {
		return func(ctx context.Context, subject string, msgType string, data interface{}, opts *PublishOptions) (future nats.PubAckFuture, err error) {
			defer func() {
				if r := recover(); r != nil {
					logger.Error("Panic recovered in async JetStream publisher",
						zap.Any("panic", r),
						zap.String("subject", subject),
						zap.String("stack", string(debug.Stack())),
					)
					err = fmt.Errorf("panic recovered: %v", r)
				}
			}()
			return next(ctx, subject, msgType, data, opts)
		}
	}
}

// RequestRecoveryMiddleware returns a middleware that recovers from panics in requests
func RequestRecoveryMiddleware(logger *zap.Logger) RequestMiddleware {
	return func(next RequestFunc) RequestFunc {
		return func(ctx context.Context, subject string, msgType string, data interface{}, timeout time.Duration, opts *PublishOptions) (resp *MessageEnvelope, err error) {
			defer func() {
				if r := recover(); r != nil {
					logger.Error("Panic recovered in request",
						zap.Any("panic", r),
						zap.String("subject", subject),
						zap.String("stack", string(debug.Stack())),
					)
					err = fmt.Errorf("panic recovered: %v", r)
				}
			}()
			return next(ctx, subject, msgType, data, timeout, opts)
		}
	}
}
