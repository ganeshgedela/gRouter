package grpc

import (
	"context"
	"testing"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// TestServerOptions tests all server option functions.
func TestServerOptions(t *testing.T) {
	logger := zap.NewNop()

	t.Run("WithReflection", func(t *testing.T) {
		server := NewServer(logger, WithReflection())
		if !server.reflection {
			t.Error("expected reflection enabled")
		}
	})

	t.Run("WithInterceptors", func(t *testing.T) {
		unary := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
			return handler(ctx, req)
		}
		stream := func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
			return handler(srv, ss)
		}

		server := NewServer(logger,
			WithInterceptors([]UnaryServerInterceptor{unary}, []StreamServerInterceptor{stream}),
		)

		if len(server.unaryInterceptors) != 1 {
			t.Errorf("expected 1 unary interceptor, got %d", len(server.unaryInterceptors))
		}
		if len(server.streamInterceptors) != 1 {
			t.Errorf("expected 1 stream interceptor, got %d", len(server.streamInterceptors))
		}
	})

	t.Run("WithUnaryInterceptor", func(t *testing.T) {
		unary := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
			return handler(ctx, req)
		}

		server := NewServer(logger, WithUnaryInterceptor(unary))

		if len(server.unaryInterceptors) != 1 {
			t.Errorf("expected 1 unary interceptor, got %d", len(server.unaryInterceptors))
		}
	})

	t.Run("WithStreamInterceptor", func(t *testing.T) {
		stream := func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
			return handler(srv, ss)
		}

		server := NewServer(logger, WithStreamInterceptor(stream))

		if len(server.streamInterceptors) != 1 {
			t.Errorf("expected 1 stream interceptor, got %d", len(server.streamInterceptors))
		}
	})

	t.Run("WithServerOptions", func(t *testing.T) {
		server := NewServer(logger,
			WithServerOptions(grpc.ConnectionTimeout(10)),
		)

		if len(server.serverOptions) != 1 {
			t.Errorf("expected 1 server option, got %d", len(server.serverOptions))
		}
	})
}
