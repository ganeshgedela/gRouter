package nats

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

// ValidatorMiddleware returns a middleware that validates data using the publisher's validator
func ValidatorMiddleware(p *NATSPublisher) PublisherMiddleware {
	return func(next PublisherFunc) PublisherFunc {
		return func(ctx context.Context, subject string, msgType string, data interface{}, opts *PublishOptions) error {
			if p.validator == nil {
				return next(ctx, subject, msgType, data, opts)
			}

			// Marshal for validation
			var dataBytes []byte
			var err error
			switch d := data.(type) {
			case []byte:
				dataBytes = d
			default:
				dataBytes, err = json.Marshal(data)
				if err != nil {
					return fmt.Errorf("failed to marshal data for validation: %w", err)
				}
			}

			// Validate
			if err := p.validator.Validate(msgType, dataBytes); err != nil {
				return fmt.Errorf("validation failed for type %s: %w", msgType, err)
			}

			// Pass marshaled bytes to next middleware to avoid re-marshalling
			return next(ctx, subject, msgType, dataBytes, opts)
		}
	}
}

// RequestValidatorMiddleware returns a middleware that validates data for requests
func RequestValidatorMiddleware(p *NATSPublisher) RequestMiddleware {
	return func(next RequestFunc) RequestFunc {
		return func(ctx context.Context, subject string, msgType string, data interface{}, timeout time.Duration, opts *PublishOptions) (*MessageEnvelope, error) {
			if p.validator == nil {
				return next(ctx, subject, msgType, data, timeout, opts)
			}

			var dataBytes []byte
			var err error
			switch d := data.(type) {
			case []byte:
				dataBytes = d
			default:
				dataBytes, err = json.Marshal(data)
				if err != nil {
					return nil, fmt.Errorf("failed to marshal data for validation: %w", err)
				}
			}

			if err := p.validator.Validate(msgType, dataBytes); err != nil {
				return nil, fmt.Errorf("validation failed for type %s: %w", msgType, err)
			}

			return next(ctx, subject, msgType, dataBytes, timeout, opts)
		}
	}
}

// JSPublisherValidatorMiddleware returns a middleware for JetStream publishing
func JSPublisherValidatorMiddleware(p *NATSPublisher) JSPublisherMiddleware {
	return func(next JSPublisherFunc) JSPublisherFunc {
		return func(ctx context.Context, subject string, msgType string, data interface{}, opts *PublishOptions) (*nats.PubAck, error) {
			if p.validator == nil {
				return next(ctx, subject, msgType, data, opts)
			}

			var dataBytes []byte
			var err error
			switch d := data.(type) {
			case []byte:
				dataBytes = d
			default:
				dataBytes, err = json.Marshal(data)
				if err != nil {
					return nil, fmt.Errorf("failed to marshal data for validation: %w", err)
				}
			}

			if err := p.validator.Validate(msgType, dataBytes); err != nil {
				return nil, fmt.Errorf("validation failed for type %s: %w", msgType, err)
			}

			return next(ctx, subject, msgType, dataBytes, opts)
		}
	}
}

// JSAsyncPublisherValidatorMiddleware returns a middleware for JetStream async publishing
func JSAsyncPublisherValidatorMiddleware(p *NATSPublisher) JSAsyncPublisherMiddleware {
	return func(next JSAsyncPublisherFunc) JSAsyncPublisherFunc {
		return func(ctx context.Context, subject string, msgType string, data interface{}, opts *PublishOptions) (nats.PubAckFuture, error) {
			if p.validator == nil {
				return next(ctx, subject, msgType, data, opts)
			}

			var dataBytes []byte
			var err error
			switch d := data.(type) {
			case []byte:
				dataBytes = d
			default:
				dataBytes, err = json.Marshal(data)
				if err != nil {
					return nil, fmt.Errorf("failed to marshal data for validation: %w", err)
				}
			}

			if err := p.validator.Validate(msgType, dataBytes); err != nil {
				return nil, fmt.Errorf("validation failed for type %s: %w", msgType, err)
			}

			return next(ctx, subject, msgType, dataBytes, opts)
		}
	}
}

// SubscriberValidatorMiddleware returns a middleware that validates incoming messages
func SubscriberValidatorMiddleware(s *NATSSubscriber) SubscriberMiddleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx context.Context, subject string, env *MessageEnvelope) error {
			if s.validator == nil {
				return next(ctx, subject, env)
			}

			// Validate
			if err := s.validator.Validate(env.Type, env.Data); err != nil {
				// Log detailed validation error
				if s.client != nil && s.client.logger != nil {
					s.client.logger.Warn("message validation failed",
						zap.String("subject", subject),
						zap.String("message_type", env.Type),
						zap.Error(err),
					)
				}
				// Return generic error (don't leak validation details)
				return fmt.Errorf("message validation failed")
			}

			return next(ctx, subject, env)
		}
	}
}
