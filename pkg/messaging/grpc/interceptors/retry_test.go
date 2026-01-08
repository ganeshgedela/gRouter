package interceptors

import (
	"context"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TestRetry tests the retry interceptor.
func TestRetry(t *testing.T) {
	config := RetryConfig{
		MaxAttempts: 3,
		InitialWait: 10 * time.Millisecond,
		MaxWait:     100 * time.Millisecond,
		Multiplier:  2.0,
		RetryableCodes: map[codes.Code]bool{
			codes.Unavailable: true,
		},
	}

	interceptor := Retry(config)

	t.Run("success on first attempt", func(t *testing.T) {
		attempts := 0
		invoker := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
			attempts++
			return nil
		}

		err := interceptor(context.Background(), "/test", nil, nil, nil, invoker)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if attempts != 1 {
			t.Errorf("expected 1 attempt, got %d", attempts)
		}
	})

	t.Run("retry on retryable error", func(t *testing.T) {
		attempts := 0
		invoker := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
			attempts++
			if attempts < 3 {
				return status.Error(codes.Unavailable, "unavailable")
			}
			return nil
		}

		err := interceptor(context.Background(), "/test", nil, nil, nil, invoker)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if attempts != 3 {
			t.Errorf("expected 3 attempts, got %d", attempts)
		}
	})

	t.Run("no retry on non-retryable error", func(t *testing.T) {
		attempts := 0
		invoker := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
			attempts++
			return status.Error(codes.InvalidArgument, "invalid")
		}

		err := interceptor(context.Background(), "/test", nil, nil, nil, invoker)
		if err == nil {
			t.Error("expected error, got nil")
		}
		if attempts != 1 {
			t.Errorf("expected 1 attempt, got %d", attempts)
		}
	})
}

// TestTimeout tests the timeout interceptor.
func TestTimeout(t *testing.T) {
	t.Run("completes before timeout", func(t *testing.T) {
		interceptor := Timeout(100 * time.Millisecond)
		handler := func(ctx context.Context, req interface{}) (interface{}, error) {
			time.Sleep(10 * time.Millisecond)
			return "success", nil
		}

		info := &grpc.UnaryServerInfo{FullMethod: "/test"}
		_, err := interceptor(context.Background(), nil, info, handler)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("times out", func(t *testing.T) {
		interceptor := Timeout(50 * time.Millisecond)
		handler := func(ctx context.Context, req interface{}) (interface{}, error) {
			time.Sleep(200 * time.Millisecond)
			return "success", nil
		}

		info := &grpc.UnaryServerInfo{FullMethod: "/test"}
		_, err := interceptor(context.Background(), nil, info, handler)
		if err == nil {
			t.Error("expected timeout error, got nil")
		}
		if status.Code(err) != codes.DeadlineExceeded {
			t.Errorf("expected DeadlineExceeded, got %v", status.Code(err))
		}
	})
}
