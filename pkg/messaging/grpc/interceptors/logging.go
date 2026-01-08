package interceptors

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// Logging creates a server interceptor that logs all RPC calls.
func Logging(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()

		// Extract peer information
		p, _ := peer.FromContext(ctx)
		peerAddr := "unknown"
		if p != nil {
			peerAddr = p.Addr.String()
		}

		// Log the incoming request
		logger.Debug("gRPC request received",
			zap.String("method", info.FullMethod),
			zap.String("peer", peerAddr),
		)

		// Call the handler
		resp, err := handler(ctx, req)

		// Calculate duration
		duration := time.Since(start)

		// Extract status code
		code := codes.OK
		if err != nil {
			if st, ok := status.FromError(err); ok {
				code = st.Code()
			}
		}

		// Log the response
		if err != nil {
			logger.Warn("gRPC request completed with error",
				zap.String("method", info.FullMethod),
				zap.String("peer", peerAddr),
				zap.Duration("duration", duration),
				zap.String("code", code.String()),
				zap.Error(err),
			)
		} else {
			logger.Debug("gRPC request completed",
				zap.String("method", info.FullMethod),
				zap.String("peer", peerAddr),
				zap.Duration("duration", duration),
				zap.String("code", code.String()),
			)
		}

		return resp, err
	}
}

// LoggingStream creates a stream interceptor that logs stream lifecycles.
func LoggingStream(logger *zap.Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		start := time.Now()

		// Extract peer information
		p, _ := peer.FromContext(ss.Context())
		peerAddr := "unknown"
		if p != nil {
			peerAddr = p.Addr.String()
		}

		// Log stream start
		logger.Debug("gRPC stream started",
			zap.String("method", info.FullMethod),
			zap.String("peer", peerAddr),
			zap.Bool("is_client_stream", info.IsClientStream),
			zap.Bool("is_server_stream", info.IsServerStream),
		)

		// Call the handler
		err := handler(srv, ss)

		// Calculate duration
		duration := time.Since(start)

		// Extract status code
		code := codes.OK
		if err != nil {
			if st, ok := status.FromError(err); ok {
				code = st.Code()
			}
		}

		// Log stream completion
		if err != nil {
			logger.Warn("gRPC stream completed with error",
				zap.String("method", info.FullMethod),
				zap.String("peer", peerAddr),
				zap.Duration("duration", duration),
				zap.String("code", code.String()),
				zap.Error(err),
			)
		} else {
			logger.Debug("gRPC stream completed",
				zap.String("method", info.FullMethod),
				zap.String("peer", peerAddr),
				zap.Duration("duration", duration),
				zap.String("code", code.String()),
			)
		}

		return err
	}
}

// LoggingClient creates a client interceptor that logs all outgoing RPC calls.
func LoggingClient(logger *zap.Logger) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		start := time.Now()

		logger.Debug("gRPC client request",
			zap.String("method", method),
		)

		err := invoker(ctx, method, req, reply, cc, opts...)

		duration := time.Since(start)
		code := codes.OK
		if err != nil {
			if st, ok := status.FromError(err); ok {
				code = st.Code()
			}
		}

		if err != nil {
			logger.Warn("gRPC client request failed",
				zap.String("method", method),
				zap.Duration("duration", duration),
				zap.String("code", code.String()),
				zap.Error(err),
			)
		} else {
			logger.Debug("gRPC client request completed",
				zap.String("method", method),
				zap.Duration("duration", duration),
				zap.String("code", code.String()),
			)
		}

		return err
	}
}
