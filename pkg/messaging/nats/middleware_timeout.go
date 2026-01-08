package nats

import (
	"context"
	"time"

	"github.com/nats-io/nats.go"
)

// TimeoutMiddleware returns a middleware that sets a default timeout on the context if not present
func TimeoutMiddleware(cfg TimeoutConfig) PublisherMiddleware {
	return func(next PublisherFunc) PublisherFunc {
		return func(ctx context.Context, subject string, msgType string, data interface{}, opts *PublishOptions) error {
			if _, ok := ctx.Deadline(); !ok && cfg.Default > 0 {
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(ctx, cfg.Default)
				defer cancel()
			}
			return next(ctx, subject, msgType, data, opts)
		}
	}
}

// RequestTimeoutMiddleware returns a middleware for requests
// Note: Requests also have an explicit timeout argument. This middleware sets the context deadline.
func RequestTimeoutMiddleware(cfg TimeoutConfig) RequestMiddleware {
	return func(next RequestFunc) RequestFunc {
		return func(ctx context.Context, subject string, msgType string, data interface{}, timeout time.Duration, opts *PublishOptions) (*MessageEnvelope, error) {
			// Logic: If timeout arg is 0, use default config.
			// Also ensure context has deadline.
			if timeout == 0 && cfg.Default > 0 {
				timeout = cfg.Default
			}

			// Even if timeout arg is present, we should update context if it lacks deadline
			if _, ok := ctx.Deadline(); !ok && timeout > 0 {
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(ctx, timeout)
				defer cancel()
			}

			return next(ctx, subject, msgType, data, timeout, opts)
		}
	}
}

// JSPublisherTimeoutMiddleware returns a middleware for JetStream sync publishing
func JSPublisherTimeoutMiddleware(cfg TimeoutConfig) JSPublisherMiddleware {
	return func(next JSPublisherFunc) JSPublisherFunc {
		return func(ctx context.Context, subject string, msgType string, data interface{}, opts *PublishOptions) (*nats.PubAck, error) {
			if _, ok := ctx.Deadline(); !ok && cfg.Default > 0 {
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(ctx, cfg.Default)
				defer cancel()
			}
			return next(ctx, subject, msgType, data, opts)
		}
	}
}

// JSAsyncPublisherTimeoutMiddleware returns a middleware for JetStream async publishing submission
func JSAsyncPublisherTimeoutMiddleware(cfg TimeoutConfig) JSAsyncPublisherMiddleware {
	return func(next JSAsyncPublisherFunc) JSAsyncPublisherFunc {
		return func(ctx context.Context, subject string, msgType string, data interface{}, opts *PublishOptions) (nats.PubAckFuture, error) {
			if _, ok := ctx.Deadline(); !ok && cfg.Default > 0 {
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(ctx, cfg.Default)
				defer cancel()
			}
			return next(ctx, subject, msgType, data, opts)
		}
	}
}
