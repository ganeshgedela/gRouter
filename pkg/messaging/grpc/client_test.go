package grpc

import (
	"context"
	"testing"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func TestNewClient(t *testing.T) {
	logger := zap.NewNop()

	config := ClientConfig{
		Target:  "localhost:9090",
		Timeout: 5 * time.Second,
	}

	client, err := NewClient(config, logger)
	if err != nil {
		// Connection might fail if server not running, but client should be created
		// In real tests, we'd mock the connection
	}

	if client == nil {
		t.Error("expected client, got nil")
	}

	if client != nil {
		defer client.Close()
	}
}

func TestClient_GetConn(t *testing.T) {
	logger := zap.NewNop()

	config := ClientConfig{
		Target:  "localhost:9090",
		Timeout: 5 * time.Second,
	}

	client, err := NewClient(config, logger)
	if err != nil {
		t.Skip("Skipping test - server not available")
	}
	defer client.Close()

	conn := client.GetConn()
	if conn == nil {
		t.Error("expected connection, got nil")
	}
}

func TestClient_WithKeepAlive(t *testing.T) {
	logger := zap.NewNop()

	config := ClientConfig{
		Target:           "localhost:9090",
		Timeout:          5 * time.Second,
		KeepAliveTime:    30 * time.Second,
		KeepAliveTimeout: 10 * time.Second,
	}

	client, err := NewClient(config, logger)
	if err != nil {
		t.Skip("Skipping test - server not available")
	}
	defer client.Close()

	if client.config.KeepAliveTime != 30*time.Second {
		t.Errorf("expected KeepAliveTime 30s, got %v", client.config.KeepAliveTime)
	}
}

func TestClient_WithInterceptors(t *testing.T) {
	logger := zap.NewNop()

	// Create a simple test interceptor
	testInterceptor := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		return invoker(ctx, method, req, reply, cc, opts...)
	}

	config := ClientConfig{
		Target:  "localhost:9090",
		Timeout: 5 * time.Second,
	}

	client, err := NewClient(config, logger,
		WithClientUnaryInterceptor(testInterceptor),
	)

	if err != nil {
		t.Skip("Skipping test - server not available")
	}
	defer client.Close()

	if len(client.unaryInterceptors) != 1 {
		t.Errorf("expected 1 unary interceptor, got %d", len(client.unaryInterceptors))
	}
}

func TestClient_Close(t *testing.T) {
	logger := zap.NewNop()

	config := ClientConfig{
		Target:  "localhost:9090",
		Timeout: 5 * time.Second,
	}

	client, err := NewClient(config, logger)
	if err != nil {
		t.Skip("Skipping test - server not available")
	}

	err = client.Close()
	if err != nil {
		t.Errorf("unexpected error closing client: %v", err)
	}
}

func TestClientOptions(t *testing.T) {
	logger := zap.NewNop()

	config := ClientConfig{
		Target:  "localhost:9090",
		Timeout: 5 * time.Second,
	}

	unaryInterceptor := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		return invoker(ctx, method, req, reply, cc, opts...)
	}

	streamInterceptor := func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		return streamer(ctx, desc, cc, method, opts...)
	}

	client, err := NewClient(config, logger,
		WithClientInterceptors([]UnaryClientInterceptor{unaryInterceptor}, []StreamClientInterceptor{streamInterceptor}),
	)

	if err != nil {
		t.Skip("Skipping test - server not available")
	}
	defer client.Close()

	if len(client.unaryInterceptors) != 1 {
		t.Errorf("expected 1 unary interceptor, got %d", len(client.unaryInterceptors))
	}

	if len(client.streamInterceptors) != 1 {
		t.Errorf("expected 1 stream interceptor, got %d", len(client.streamInterceptors))
	}
}
