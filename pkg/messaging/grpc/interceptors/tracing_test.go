package interceptors

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel/trace/noop"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TestTracing tests the tracing interceptor.
func TestTracing(t *testing.T) {
	tracer := noop.NewTracerProvider().Tracer("test")
	interceptor := Tracing(tracer)

	t.Run("successful request", func(t *testing.T) {
		handler := func(ctx context.Context, req interface{}) (interface{}, error) {
			return "response", nil
		}

		info := &grpc.UnaryServerInfo{FullMethod: "/test.Service/Method"}

		_, err := interceptor(context.Background(), nil, info, handler)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("failed request", func(t *testing.T) {
		handler := func(ctx context.Context, req interface{}) (interface{}, error) {
			return nil, status.Error(codes.Internal, "test error")
		}

		info := &grpc.UnaryServerInfo{FullMethod: "/test.Service/Method"}

		_, err := interceptor(context.Background(), nil, info, handler)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

// TestTracingStream tests the stream tracing interceptor.
func TestTracingStream(t *testing.T) {
	tracer := noop.NewTracerProvider().Tracer("test")
	interceptor := TracingStream(tracer)

	handler := func(srv interface{}, stream grpc.ServerStream) error {
		return nil
	}

	info := &grpc.StreamServerInfo{
		FullMethod:     "/test.Service/StreamMethod",
		IsClientStream: true,
		IsServerStream: true,
	}

	// Create a mock server stream with context
	mockStream := &mockServerStream{ctx: context.Background()}

	err := interceptor(nil, mockStream, info, handler)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// TestTracingClient tests the client tracing interceptor.
func TestTracingClient(t *testing.T) {
	tracer := noop.NewTracerProvider().Tracer("test")
	interceptor := TracingClient(tracer)

	t.Run("successful call", func(t *testing.T) {
		invoker := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
			return nil
		}

		err := interceptor(context.Background(), "/test.Service/Method", nil, nil, nil, invoker)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("failed call", func(t *testing.T) {
		invoker := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
			return status.Error(codes.Unavailable, "service unavailable")
		}

		err := interceptor(context.Background(), "/test.Service/Method", nil, nil, nil, invoker)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

// TestExtractServiceName tests the extractServiceName helper.
func TestExtractServiceName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"/helloworld.Greeter/SayHello", "helloworld.Greeter"},
		{"helloworld.Greeter/SayHello", "helloworld.Greeter"},
		{"/service/method", "service"},
		{"", ""},
		{"/nomethod", "nomethod"},
	}

	for _, tt := range tests {
		result := extractServiceName(tt.input)
		if result != tt.expected {
			t.Errorf("extractServiceName(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

// TestExtractMethodName tests the extractMethodName helper.
func TestExtractMethodName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"/helloworld.Greeter/SayHello", "SayHello"},
		{"helloworld.Greeter/SayHello", "SayHello"},
		{"/service/method", "method"},
		{"", ""},
		{"nomethod", "nomethod"},
	}

	for _, tt := range tests {
		result := extractMethodName(tt.input)
		if result != tt.expected {
			t.Errorf("extractMethodName(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

// mockServerStream implements grpc.ServerStream for testing.
type mockServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (m *mockServerStream) Context() context.Context {
	return m.ctx
}
