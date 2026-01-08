package interceptors

import (
	"context"
	"testing"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// mockServerStream implements grpc.ServerStream for testing.
type mockServerStreamForLogging struct {
	grpc.ServerStream
	ctx context.Context
}

func (m *mockServerStreamForLogging) Context() context.Context {
	return m.ctx
}

// TestLoggingStream tests the stream logging interceptor.
func TestLoggingStream(t *testing.T) {
	logger := zap.NewNop()
	interceptor := LoggingStream(logger)

	t.Run("successful stream", func(t *testing.T) {
		handler := func(srv interface{}, stream grpc.ServerStream) error {
			return nil
		}

		info := &grpc.StreamServerInfo{
			FullMethod:     "/test.Service/StreamMethod",
			IsClientStream: true,
			IsServerStream: true,
		}

		mockStream := &mockServerStreamForLogging{ctx: context.Background()}

		err := interceptor(nil, mockStream, info, handler)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("failed stream", func(t *testing.T) {
		handler := func(srv interface{}, stream grpc.ServerStream) error {
			return status.Error(codes.Internal, "stream error")
		}

		info := &grpc.StreamServerInfo{
			FullMethod:     "/test.Service/StreamMethod",
			IsClientStream: false,
			IsServerStream: true,
		}

		mockStream := &mockServerStreamForLogging{ctx: context.Background()}

		err := interceptor(nil, mockStream, info, handler)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}
