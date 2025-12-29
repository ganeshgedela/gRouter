package nats

import (
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestNewClient(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: Config{
				URL:               "nats://localhost:4222",
				MaxReconnects:     10,
				ReconnectWait:     2 * time.Second,
				ConnectionTimeout: 5 * time.Second,
			},
			wantErr: false,
		},
		{
			name:    "nil logger",
			config:  Config{URL: "nats://localhost:4222"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var client *Client
			var err error

			if tt.name == "nil logger" {
				client, err = NewNATSClient(tt.config, nil)
			} else {
				client, err = NewNATSClient(tt.config, logger)
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("NewNATSClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && client == nil {
				t.Error("NewNATSClient() returned nil client")
			}
		})
	}
}

func TestClient_IsConnected(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	config := Config{
		URL:               "nats://localhost:4222",
		MaxReconnects:     10,
		ReconnectWait:     2 * time.Second,
		ConnectionTimeout: 5 * time.Second,
	}

	client, err := NewNATSClient(config, logger)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Should not be connected initially
	if client.IsConnected() {
		t.Error("Client should not be connected before Connect() is called")
	}
}

func TestClient_ConnectAndClose(t *testing.T) {
	// Skip if NATS server is not available
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	logger, _ := zap.NewDevelopment()
	config := Config{
		URL:               "nats://localhost:4222",
		MaxReconnects:     10,
		ReconnectWait:     2 * time.Second,
		ConnectionTimeout: 5 * time.Second,
	}

	client, err := NewNATSClient(config, logger)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test connection
	err = client.Connect()
	if err != nil || !client.IsConnected() {
		t.Skipf("NATS server not available or not connected: %v", err)
		return
	}

	// Should be connected
	if !client.IsConnected() {
		t.Error("Client should be connected after Connect()")
	}

	// Test close
	err = client.Close()
	if err != nil {
		t.Errorf("Close() error = %v", err)
	}

	// Give it a moment to close
	time.Sleep(100 * time.Millisecond)

	// Should not be connected after close
	if client.IsConnected() {
		t.Error("Client should not be connected after Close()")
	}
}

func TestClient_Conn(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	config := Config{
		URL:               "nats://localhost:4222",
		MaxReconnects:     10,
		ReconnectWait:     2 * time.Second,
		ConnectionTimeout: 5 * time.Second,
	}

	client, err := NewNATSClient(config, logger)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Should return nil before connection
	if client.Conn() != nil {
		t.Error("Conn() should return nil before Connect()")
	}
}

func TestClient_WithAuthentication(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	tests := []struct {
		name   string
		config Config
	}{
		{
			name: "with token",
			config: Config{
				URL:               "nats://localhost:4222",
				Token:             "test-token",
				MaxReconnects:     10,
				ReconnectWait:     2 * time.Second,
				ConnectionTimeout: 5 * time.Second,
			},
		},
		{
			name: "with username and password",
			config: Config{
				URL:               "nats://localhost:4222",
				Username:          "testuser",
				Password:          "testpass",
				MaxReconnects:     10,
				ReconnectWait:     2 * time.Second,
				ConnectionTimeout: 5 * time.Second,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewNATSClient(tt.config, logger)
			if err != nil {
				t.Fatalf("Failed to create client: %v", err)
			}

			// Just verify client was created with auth config
			// Actual connection would fail without proper NATS server setup
			if client == nil {
				t.Error("Client should not be nil")
			}
		})
	}
}
