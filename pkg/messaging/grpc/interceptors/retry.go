package interceptors

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RetryConfig holds configuration for the retry interceptor.
type RetryConfig struct {
	MaxAttempts    int
	InitialWait    time.Duration
	MaxWait        time.Duration
	Multiplier     float64
	RetryableCodes map[codes.Code]bool
}

// DefaultRetryConfig returns a sensible default retry configuration.
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxAttempts: 3,
		InitialWait: 100 * time.Millisecond,
		MaxWait:     5 * time.Second,
		Multiplier:  2.0,
		RetryableCodes: map[codes.Code]bool{
			codes.Unavailable:       true,
			codes.ResourceExhausted: true,
			codes.Aborted:           true,
		},
	}
}

// Retry creates a client interceptor that retries failed requests with exponential backoff.
func Retry(config RetryConfig) grpc.UnaryClientInterceptor {
	if config.RetryableCodes == nil {
		config = DefaultRetryConfig()
	}

	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		var lastErr error
		wait := config.InitialWait

		for attempt := 1; attempt <= config.MaxAttempts; attempt++ {
			// Try the request
			lastErr = invoker(ctx, method, req, reply, cc, opts...)

			// If successful, return
			if lastErr == nil {
				return nil
			}

			// Check if error is retryable
			st, ok := status.FromError(lastErr)
			if !ok || !config.RetryableCodes[st.Code()] {
				// Not retryable, return error
				return lastErr
			}

			// Don't retry on last attempt
			if attempt == config.MaxAttempts {
				break
			}

			// Wait before retrying
			select {
			case <-time.After(wait):
				// Calculate next wait with exponential backoff
				wait = time.Duration(float64(wait) * config.Multiplier)
				if wait > config.MaxWait {
					wait = config.MaxWait
				}
			case <-ctx.Done():
				return ctx.Err()
			}
		}

		return lastErr
	}
}
