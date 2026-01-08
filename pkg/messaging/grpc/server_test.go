package grpc

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func TestNewServer(t *testing.T) {
	logger := zap.NewNop()
	s := NewServer(logger, WithPort(8080))
	assert.NotNil(t, s)
	assert.Equal(t, 8080, s.port)
}

func TestServer_Start_Stop(t *testing.T) {
	logger := zap.NewNop()
	// Use port 0 for dynamic assignment to avoid conflicts
	s := NewServer(logger, WithPort(0))

	err := s.Start()
	require.NoError(t, err)
	assert.NotZero(t, s.port, "Port should be updated after Start")

	// Allow some time for startup
	time.Sleep(100 * time.Millisecond)

	s.Stop(context.Background())
}

func TestServer_WithGateway(t *testing.T) {
	logger := zap.NewNop()
	// Gateway requires a separate port, we'll let it be dynamic if possible or pick likely unused
	// Our implementation of WithGateway expects a port.
	// Since we can't easily predict a free port, we'll listen on 0 for gRPC
	// but for Gateway we also need a port.
	// Let's use 0 for gateway port? net/http Server Addr ":0" works.

	s := NewServer(logger, WithPort(0), WithGateway(context.Background(), 0))

	err := s.Start()
	require.NoError(t, err)

	// Gateway server should happen to have started.
	// However, we didn't capture Gateway port in server struct specifically for reading back in test
	// (only s.gatewayPorts which stored what we passed: 0).
	// If we passed 0, http.Server uses random port, but we don't know which one unless we capture listener of http server too.
	// My implementation of startGateway uses http.ListenAndServe/Serve which manages listener internally unless we create it properly.
	// In startGateway: s.gatewaySrv.ListenAndServe()
	// If Addr is ":0", it listens on random port. But we can't query it easily because we didn't store the listener.
	// For this test, we just verify it doesn't crash properly.

	time.Sleep(100 * time.Millisecond)
	s.Stop(context.Background())
}

// Mock service registration
func mockRegister(registrar grpc.ServiceRegistrar) {
	// Do nothing
}

func TestServer_RegisterService(t *testing.T) {
	logger := zap.NewNop()
	s := NewServer(logger)
	s.RegisterService(mockRegister)
	// Success if no panic
}
