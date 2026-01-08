package interceptors

import (
	"context"
	"testing"

	"google.golang.org/grpc"
)

// TestMetrics tests the metrics interceptor.
func TestMetrics(t *testing.T) {
	interceptor := Metrics()

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "response", nil
	}

	info := &grpc.UnaryServerInfo{FullMethod: "/test.Service/Method"}

	_, err := interceptor(context.Background(), nil, info, handler)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Metrics are recorded, we can't easily assert on Prometheus metrics in unit tests
	// but we can verify the interceptor doesn't error
}

// TestMetricsStream tests the stream metrics interceptor.
func TestMetricsStream(t *testing.T) {
	interceptor := MetricsStream()

	handler := func(srv interface{}, stream grpc.ServerStream) error {
		return nil
	}

	info := &grpc.StreamServerInfo{FullMethod: "/test.Service/StreamMethod"}

	err := interceptor(nil, nil, info, handler)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// TestMetricsClient tests the client metrics interceptor.
func TestMetricsClient(t *testing.T) {
	interceptor := MetricsClient()

	invoker := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		return nil
	}

	err := interceptor(context.Background(), "/test.Service/Method", nil, nil, nil, invoker)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
