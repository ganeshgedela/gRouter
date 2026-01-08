package nats

import (
	"context"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/nats-io/nats.go"
)

// withBackoff executes an operation with exponential backoff
func withBackoff(ctx context.Context, cfg RetryConfig, op func() error) error {
	b := backoff.NewExponentialBackOff()
	b.InitialInterval = cfg.InitialInterval
	b.Multiplier = cfg.Multiplier
	b.MaxInterval = cfg.MaxInterval
	b.MaxElapsedTime = 0 // Disable time limit, relying on attempts/context

	// Ensure defaults if zero
	if b.InitialInterval == 0 {
		b.InitialInterval = 100 * time.Millisecond
	}
	if b.Multiplier == 0 {
		b.Multiplier = 2.0
	}
	if b.MaxInterval == 0 {
		b.MaxInterval = 1 * time.Second
	}

	var bOff backoff.BackOff = b

	// MaxAttempts includes the initial attempt. Retries = MaxAttempts - 1
	maxRetries := uint64(2) // Default 3 attempts total
	if cfg.MaxAttempts > 0 {
		maxRetries = uint64(cfg.MaxAttempts - 1)
	}
	bOff = backoff.WithMaxRetries(b, maxRetries)
	bOff = backoff.WithContext(bOff, ctx)

	return backoff.Retry(op, bOff)
}

// RetryMiddleware returns a middleware that retries sync publish operations
func RetryMiddleware(cfg RetryConfig) PublisherMiddleware {
	return func(next PublisherFunc) PublisherFunc {
		return func(ctx context.Context, subject string, msgType string, data interface{}, opts *PublishOptions) error {
			op := func() error {
				return next(ctx, subject, msgType, data, opts)
			}
			return withBackoff(ctx, cfg, op)
		}
	}
}

// RequestRetryMiddleware returns a middleware that retries request operations
func RequestRetryMiddleware(cfg RetryConfig) RequestMiddleware {
	return func(next RequestFunc) RequestFunc {
		return func(ctx context.Context, subject string, msgType string, data interface{}, timeout time.Duration, opts *PublishOptions) (*MessageEnvelope, error) {
			var resp *MessageEnvelope
			op := func() error {
				var err error
				resp, err = next(ctx, subject, msgType, data, timeout, opts)
				return err
			}
			if err := withBackoff(ctx, cfg, op); err != nil {
				return nil, err
			}
			return resp, nil
		}
	}
}

// JSPublisherRetryMiddleware returns a middleware that retries JetStream sync publish operations
func JSPublisherRetryMiddleware(cfg RetryConfig) JSPublisherMiddleware {
	return func(next JSPublisherFunc) JSPublisherFunc {
		return func(ctx context.Context, subject string, msgType string, data interface{}, opts *PublishOptions) (*nats.PubAck, error) {
			var ack *nats.PubAck
			op := func() error {
				var err error
				ack, err = next(ctx, subject, msgType, data, opts)
				return err
			}
			if err := withBackoff(ctx, cfg, op); err != nil {
				return nil, err
			}
			return ack, nil
		}
	}
}

// JSAsyncPublisherRetryMiddleware returns a middleware for retry of JetStream async publish SUBMISSION
// It does NOT retry on async ack failures.
func JSAsyncPublisherRetryMiddleware(cfg RetryConfig) JSAsyncPublisherMiddleware {
	return func(next JSAsyncPublisherFunc) JSAsyncPublisherFunc {
		return func(ctx context.Context, subject string, msgType string, data interface{}, opts *PublishOptions) (nats.PubAckFuture, error) {
			var future nats.PubAckFuture
			op := func() error {
				var err error
				future, err = next(ctx, subject, msgType, data, opts)
				return err
			}
			if err := withBackoff(ctx, cfg, op); err != nil {
				return nil, err
			}
			return future, nil
		}
	}
}
