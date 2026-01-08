package nats

import (
	"context"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/sony/gobreaker"
	"go.uber.org/zap"
)

// NewCircuitBreaker creates a new circuit breaker from config
func NewCircuitBreaker(cfg CircuitBreakerConfig, logger *zap.Logger) *gobreaker.CircuitBreaker {
	settings := gobreaker.Settings{
		Name:        "NATS-Client-CB",
		MaxRequests: cfg.MaxRequests,
		Interval:    cfg.Interval,
		Timeout:     cfg.Timeout,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			threshold := cfg.TripThreshold
			if threshold == 0 {
				threshold = 5 // Default to 5 if not configured
			}
			return counts.ConsecutiveFailures >= threshold
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			if logger != nil {
				logger.Warn("circuit breaker state change",
					zap.String("name", name),
					zap.String("from", from.String()),
					zap.String("to", to.String()),
				)
			}
		},
	}
	return gobreaker.NewCircuitBreaker(settings)
}

// CircuitBreakerMiddleware returns a middleware that uses a circuit breaker for sync publishing
func CircuitBreakerMiddleware(cb *gobreaker.CircuitBreaker, logger *zap.Logger) PublisherMiddleware {
	return func(next PublisherFunc) PublisherFunc {
		return func(ctx context.Context, subject string, msgType string, data interface{}, opts *PublishOptions) error {
			_, err := cb.Execute(func() (interface{}, error) {
				return nil, next(ctx, subject, msgType, data, opts)
			})
			if err != nil && err == gobreaker.ErrOpenState {
				if logger != nil {
					logger.Warn("circuit breaker open",
						zap.String("subject", subject),
						zap.String("message_type", msgType),
					)
				}
				return fmt.Errorf("service temporarily unavailable")
			}
			return err
		}
	}
}

// RequestCircuitBreakerMiddleware returns a middleware for requests
func RequestCircuitBreakerMiddleware(cb *gobreaker.CircuitBreaker) RequestMiddleware {
	return func(next RequestFunc) RequestFunc {
		return func(ctx context.Context, subject string, msgType string, data interface{}, timeout time.Duration, opts *PublishOptions) (*MessageEnvelope, error) {
			resp, err := cb.Execute(func() (interface{}, error) {
				return next(ctx, subject, msgType, data, timeout, opts)
			})
			if err != nil {
				return nil, err
			}
			return resp.(*MessageEnvelope), nil
		}
	}
}

// JSPublisherCircuitBreakerMiddleware returns a middleware for JetStream sync publishing
func JSPublisherCircuitBreakerMiddleware(cb *gobreaker.CircuitBreaker) JSPublisherMiddleware {
	return func(next JSPublisherFunc) JSPublisherFunc {
		return func(ctx context.Context, subject string, msgType string, data interface{}, opts *PublishOptions) (*nats.PubAck, error) {
			ack, err := cb.Execute(func() (interface{}, error) {
				return next(ctx, subject, msgType, data, opts)
			})
			if err != nil {
				return nil, err
			}
			return ack.(*nats.PubAck), nil
		}
	}
}

// JSAsyncPublisherCircuitBreakerMiddleware returns a middleware for JetStream async publishing
// Note: This only guards against immediate submission errors, not async ack failures.
func JSAsyncPublisherCircuitBreakerMiddleware(cb *gobreaker.CircuitBreaker) JSAsyncPublisherMiddleware {
	return func(next JSAsyncPublisherFunc) JSAsyncPublisherFunc {
		return func(ctx context.Context, subject string, msgType string, data interface{}, opts *PublishOptions) (nats.PubAckFuture, error) {
			future, err := cb.Execute(func() (interface{}, error) {
				return next(ctx, subject, msgType, data, opts)
			})
			if err != nil {
				return nil, err
			}
			return future.(nats.PubAckFuture), nil
		}
	}
}
