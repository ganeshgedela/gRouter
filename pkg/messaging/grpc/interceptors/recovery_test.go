package interceptors

import (
	"context"
	"errors"
	"testing"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// TestRecovery tests the panic recovery interceptor.
func TestRecovery(t *testing.T) {
	logger := zap.NewNop()
	interceptor := Recovery(logger)

	tests := []struct {
		name        string
		handler     grpc.UnaryHandler
		shouldPanic bool
		wantErr     bool
	}{
		{
			name: "no panic",
			handler: func(ctx context.Context, req interface{}) (interface{}, error) {
				return "success", nil
			},
			shouldPanic: false,
			wantErr:     false,
		},
		{
			name: "panic with string",
			handler: func(ctx context.Context, req interface{}) (interface{}, error) {
				panic("test panic")
			},
			shouldPanic: true,
			wantErr:     true,
		},
		{
			name: "panic with error",
			handler: func(ctx context.Context, req interface{}) (interface{}, error) {
				panic(errors.New("test error"))
			},
			shouldPanic: true,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := &grpc.UnaryServerInfo{FullMethod: "/test.Service/Method"}

			resp, err := interceptor(context.Background(), nil, info, tt.handler)

			if tt.wantErr && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !tt.shouldPanic && resp == nil && err == nil {
				t.Error("expected response, got nil")
			}
		})
	}
}

// TestRecoveryStream tests the stream panic recovery interceptor.
func TestRecoveryStream(t *testing.T) {
	logger := zap.NewNop()
	interceptor := RecoveryStream(logger)

	tests := []struct {
		name        string
		handler     grpc.StreamHandler
		shouldPanic bool
		wantErr     bool
	}{
		{
			name: "no panic",
			handler: func(srv interface{}, stream grpc.ServerStream) error {
				return nil
			},
			shouldPanic: false,
			wantErr:     false,
		},
		{
			name: "panic in stream",
			handler: func(srv interface{}, stream grpc.ServerStream) error {
				panic("stream panic")
			},
			shouldPanic: true,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := &grpc.StreamServerInfo{FullMethod: "/test.Service/StreamMethod"}

			err := interceptor(nil, nil, info, tt.handler)

			if tt.wantErr && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
