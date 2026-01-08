package interceptors

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Timeout creates a server interceptor that enforces a timeout on requests.
func Timeout(timeout time.Duration) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Create context with timeout
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		// Channel to capture handler result
		type result struct {
			resp interface{}
			err  error
		}
		resultChan := make(chan result, 1)

		// Run handler in goroutine
		go func() {
			resp, err := handler(ctx, req)
			resultChan <- result{resp, err}
		}()

		// Wait for either handler completion or timeout
		select {
		case res := <-resultChan:
			return res.resp, res.err
		case <-ctx.Done():
			return nil, status.Errorf(codes.DeadlineExceeded, "request timeout exceeded")
		}
	}
}

// TimeoutClient creates a client interceptor that sets a timeout on requests.
func TimeoutClient(timeout time.Duration) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		// Create context with timeout if not already set
		if _, ok := ctx.Deadline(); !ok {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, timeout)
			defer cancel()
		}

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
