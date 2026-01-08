package interceptors

import (
	"context"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TestRateLimiter tests the rate limiter interceptor.
func TestRateLimiter(t *testing.T) {
	// Use very high rate to avoid timing issues
	interceptor := RateLimiter(1000, 100)

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "response", nil
	}

	info := &grpc.UnaryServerInfo{FullMethod: "/test.Service/Method"}

	// All requests should succeed with high limits
	for i := 0; i < 10; i++ {
		_, err := interceptor(context.Background(), nil, info, handler)
		if err != nil {
			t.Errorf("request %d failed: %v", i+1, err)
		}
	}
}

// TestRateLimiterExceeded tests rate limit enforcement.
func TestRateLimiterExceeded(t *testing.T) {
	// Very low rate to guarantee limit hit
	interceptor := RateLimiter(1, 1) // 1 req/s, 1 burst

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "response", nil
	}

	info := &grpc.UnaryServerInfo{FullMethod: "/test.Service/Method"}

	// First request should succeed
	_, err := interceptor(context.Background(), nil, info, handler)
	if err != nil {
		t.Errorf("first request failed: %v", err)
	}

	// Immediate second request should be rate limited
	_, err = interceptor(context.Background(), nil, info, handler)
	if err == nil {
		t.Error("expected rate limit error, got nil")
	}
	if status.Code(err) != codes.ResourceExhausted {
		t.Errorf("expected ResourceExhausted, got %v", status.Code(err))
	}
}

// TestPerMethodRateLimiter tests per-method rate limiting.
func TestPerMethodRateLimiter(t *testing.T) {
	limits := NewMethodRateLimits()
	limits.Add("/test.Service/Method1", 10, 1)
	limits.Add("/test.Service/Method2", 5, 1)

	interceptor := PerMethodRateLimiterInterceptor(limits, 100, 10)

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "response", nil
	}

	// Test Method1
	info1 := &grpc.UnaryServerInfo{FullMethod: "/test.Service/Method1"}
	_, err := interceptor(context.Background(), nil, info1, handler)
	if err != nil {
		t.Errorf("Method1 failed: %v", err)
	}

	// Test Method2
	info2 := &grpc.UnaryServerInfo{FullMethod: "/test.Service/Method2"}
	_, err = interceptor(context.Background(), nil, info2, handler)
	if err != nil {
		t.Errorf("Method2 failed: %v", err)
	}

	// Test unknown method (uses default)
	info3 := &grpc.UnaryServerInfo{FullMethod: "/test.Service/Unknown"}
	_, err = interceptor(context.Background(), nil, info3, handler)
	if err != nil {
		t.Errorf("Unknown method failed: %v", err)
	}
}

// TestCircuitBreaker tests the circuit breaker.
func TestCircuitBreaker(t *testing.T) {
	cb := NewCircuitBreaker(3, 100*time.Millisecond, 1)

	t.Run("closed state allows requests", func(t *testing.T) {
		err := cb.Call(func() error {
			return nil
		})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("failures open circuit", func(t *testing.T) {
		// Trigger 3 failures to open circuit
		for i := 0; i < 3; i++ {
			cb.Call(func() error {
				return status.Error(codes.Unavailable, "error")
			})
		}

		// Next call should be rejected
		err := cb.Call(func() error {
			return nil
		})
		if err == nil {
			t.Error("expected circuit open error, got nil")
		}
		if status.Code(err) != codes.Unavailable {
			t.Errorf("expected Unavailable, got %v", status.Code(err))
		}
	})

	t.Run("circuit transitions to half-open after timeout", func(t *testing.T) {
		// Wait for timeout
		time.Sleep(150 * time.Millisecond)

		// Should allow request in half-open state
		err := cb.Call(func() error {
			return nil
		})
		if err != nil {
			t.Errorf("half-open request failed: %v", err)
		}

		// Success should close circuit
		err = cb.Call(func() error {
			return nil
		})
		if err != nil {
			t.Errorf("closed circuit request failed: %v", err)
		}
	})
}

// TestCircuitBreakerInterceptor tests the circuit breaker interceptor.
func TestCircuitBreakerInterceptor(t *testing.T) {
	cb := NewCircuitBreaker(2, 100*time.Millisecond, 1)
	interceptor := CircuitBreakerInterceptor(cb)

	successInvoker := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		return nil
	}

	failInvoker := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		return status.Error(codes.Internal, "error")
	}

	// First request succeeds
	err := interceptor(context.Background(), "/test", nil, nil, nil, successInvoker)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Trigger failures to open circuit
	interceptor(context.Background(), "/test", nil, nil, nil, failInvoker)
	interceptor(context.Background(), "/test", nil, nil, nil, failInvoker)

	// Circuit should be open
	err = interceptor(context.Background(), "/test", nil, nil, nil, successInvoker)
	if err == nil {
		t.Error("expected circuit open error, got nil")
	}
}

// TestTimeoutClient tests the client timeout interceptor.
func TestTimeoutClient(t *testing.T) {
	interceptor := TimeoutClient(100 * time.Millisecond)

	t.Run("respects existing deadline", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		invoker := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
			_, ok := ctx.Deadline()
			if !ok {
				t.Error("expected deadline, got none")
			}
			return nil
		}

		err := interceptor(ctx, "/test", nil, nil, nil, invoker)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("sets deadline if none exists", func(t *testing.T) {
		invoker := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
			deadline, ok := ctx.Deadline()
			if !ok {
				t.Error("expected deadline, got none")
			}
			if time.Until(deadline) > 100*time.Millisecond {
				t.Error("deadline not set correctly")
			}
			return nil
		}

		err := interceptor(context.Background(), "/test", nil, nil, nil, invoker)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}
