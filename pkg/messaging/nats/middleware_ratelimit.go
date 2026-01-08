package nats

import (
	"context"
	"time"

	"github.com/nats-io/nats.go"
	"golang.org/x/time/rate"
)

// NewRateLimiter creates a new rate limiter from config
func NewRateLimiter(cfg RateLimitConfig) *rate.Limiter {
	limit := rate.Limit(cfg.RequestsPerSecond)
	// Safety: if limit is 0, defaults to Inf to avoid blocking forever if config is missing but enabled?
	// But enabled=true with 0 limit means "allow nothing".
	// Let's trust config.
	return rate.NewLimiter(limit, cfg.Burst)
}

// RateLimitMiddleware returns a middleware that waits for rate limit permission before publishing
func RateLimitMiddleware(limiter *rate.Limiter) PublisherMiddleware {
	return func(next PublisherFunc) PublisherFunc {
		return func(ctx context.Context, subject string, msgType string, data interface{}, opts *PublishOptions) error {
			if err := limiter.Wait(ctx); err != nil {
				return err // Context canceled or deadline exceeded
			}
			return next(ctx, subject, msgType, data, opts)
		}
	}
}

// RequestRateLimitMiddleware returns a middleware for requests
func RequestRateLimitMiddleware(limiter *rate.Limiter) RequestMiddleware {
	return func(next RequestFunc) RequestFunc {
		return func(ctx context.Context, subject string, msgType string, data interface{}, timeout time.Duration, opts *PublishOptions) (*MessageEnvelope, error) {
			if err := limiter.Wait(ctx); err != nil {
				return nil, err
			}
			return next(ctx, subject, msgType, data, timeout, opts)
		}
	}
}

// JSPublisherRateLimitMiddleware returns a middleware for JetStream sync publishing
func JSPublisherRateLimitMiddleware(limiter *rate.Limiter) JSPublisherMiddleware {
	return func(next JSPublisherFunc) JSPublisherFunc {
		return func(ctx context.Context, subject string, msgType string, data interface{}, opts *PublishOptions) (*nats.PubAck, error) {
			if err := limiter.Wait(ctx); err != nil {
				return nil, err
			}
			return next(ctx, subject, msgType, data, opts)
		}
	}
}

// JSAsyncPublisherRateLimitMiddleware returns a middleware for JetStream async publishing submission
func JSAsyncPublisherRateLimitMiddleware(limiter *rate.Limiter) JSAsyncPublisherMiddleware {
	return func(next JSAsyncPublisherFunc) JSAsyncPublisherFunc {
		return func(ctx context.Context, subject string, msgType string, data interface{}, opts *PublishOptions) (nats.PubAckFuture, error) {
			if err := limiter.Wait(ctx); err != nil {
				return nil, err
			}
			return next(ctx, subject, msgType, data, opts)
		}
	}
}
