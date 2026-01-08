package interceptors

import (
	"context"
	"fmt"
	"runtime/debug"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Recovery creates a server interceptor that recovers from panics.
func Recovery(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				stack := debug.Stack()
				logger.Error("Panic recovered in gRPC handler",
					zap.String("method", info.FullMethod),
					zap.Any("panic", r),
					zap.ByteString("stack", stack),
				)

				// Convert panic to gRPC error
				err = status.Errorf(codes.Internal, "internal server error: %v", r)
			}
		}()

		return handler(ctx, req)
	}
}

// RecoveryStream creates a stream interceptor that recovers from panics.
func RecoveryStream(logger *zap.Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		defer func() {
			if r := recover(); r != nil {
				stack := debug.Stack()
				logger.Error("Panic recovered in gRPC stream handler",
					zap.String("method", info.FullMethod),
					zap.Any("panic", r),
					zap.ByteString("stack", stack),
				)

				// Convert panic to gRPC error
				err = status.Errorf(codes.Internal, "internal server error: %v", r)
			}
		}()

		return handler(srv, ss)
	}
}

// recoverFrom converts a panic into a gRPC error with appropriate status code.
func recoverFrom(p interface{}) error {
	var msg string
	switch v := p.(type) {
	case string:
		msg = v
	case error:
		msg = v.Error()
	default:
		msg = fmt.Sprintf("%v", v)
	}

	return status.Errorf(codes.Internal, "panic: %s", msg)
}
