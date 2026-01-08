package interceptors

import (
	"context"
	"testing"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TestLogging tests the logging interceptor.
func TestLogging(t *testing.T) {
	logger := zap.NewNop()
	interceptor := Logging(logger)

	tests := []struct {
		name    string
		handler grpc.UnaryHandler
		wantErr bool
	}{
		{
			name: "successful request",
			handler: func(ctx context.Context, req interface{}) (interface{}, error) {
				return "response", nil
			},
			wantErr: false,
		},
		{
			name: "failed request",
			handler: func(ctx context.Context, req interface{}) (interface{}, error) {
				return nil, status.Error(codes.Internal, "test error")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := &grpc.UnaryServerInfo{FullMethod: "/test.Service/Method"}

			_, err := interceptor(context.Background(), nil, info, tt.handler)

			if tt.wantErr && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

// TestLoggingClient tests the client logging interceptor.
func TestLoggingClient(t *testing.T) {
	logger := zap.NewNop()
	interceptor := LoggingClient(logger)

	invoker := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		return nil
	}

	err := interceptor(context.Background(), "/test.Service/Method", nil, nil, nil, invoker)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
