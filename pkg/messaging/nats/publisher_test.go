package nats

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

func TestNewPublisher(t *testing.T) {
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

	publisher := NewPublisher(client, "test-service")
	if publisher == nil {
		t.Fatalf("NewPublisher() returned nil")
	}

	// Type assertion to access struct fields for testing
	natsPub, ok := publisher.(*NATSPublisher)
	if !ok {
		t.Fatalf("Publisher is not *NATSPublisher")
	}

	if natsPub.source != "test-service" {
		t.Errorf("Publisher source = %v, want %v", natsPub.source, "test-service")
	}
}

func TestPublisher_Publish_NotConnected(t *testing.T) {
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

	publisher := NewPublisher(client, "test-service")

	// Try to publish without connection
	err = publisher.Publish(context.Background(), "test.subject", "test.event", map[string]string{"key": "value"}, nil)
	if err == nil {
		t.Error("Publish() should return error when not connected")
	}
}

func TestPublisher_Publish_Integration(t *testing.T) {
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

	err = client.Connect()
	if err != nil || !client.IsConnected() {
		t.Skipf("NATS server not available or not connected: %v", err)
		return
	}
	defer client.Close()

	publisher := NewPublisher(client, "test-service")

	tests := []struct {
		name    string
		subject string
		msgType string
		data    interface{}
		opts    *PublishOptions
		wantErr bool
	}{
		{
			name:    "sync publish",
			subject: "test.sync",
			msgType: "test.event",
			data:    map[string]string{"key": "value"},
			opts:    nil,
			wantErr: false,
		},
		{
			name:    "async publish",
			subject: "test.async",
			msgType: "test.event",
			data:    map[string]string{"key": "value"},
			opts:    &PublishOptions{Async: true},
			wantErr: false,
		},
		{
			name:    "publish with struct",
			subject: "test.struct",
			msgType: "test.event",
			data: struct {
				Name  string
				Value int
			}{Name: "test", Value: 42},
			opts:    nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := publisher.Publish(context.Background(), tt.subject, tt.msgType, tt.data, tt.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("Publish() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPublisher_Request_Integration(t *testing.T) {
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

	err = client.Connect()
	if err != nil || !client.IsConnected() {
		t.Skipf("NATS server not available or not connected: %v", err)
		return
	}
	defer client.Close()

	// Set up a responder
	_, err = client.Conn().Subscribe("test.request", func(msg *nats.Msg) {
		response := MessageEnvelope{
			ID:        "response-1",
			Type:      "test.response",
			Timestamp: time.Now(),
			Source:    "responder",
			Data:      json.RawMessage(`{"result":"success"}`),
		}
		data, _ := json.Marshal(response)
		msg.Respond(data)
	})
	if err != nil {
		t.Fatalf("Failed to set up responder: %v", err)
	}

	publisher := NewPublisher(client, "test-service")

	// Test request
	response, err := publisher.Request(context.Background(), "test.request", "test.request", map[string]string{"key": "value"}, 2*time.Second)
	if err != nil {
		t.Errorf("Request() error = %v", err)
		return
	}

	if response == nil {
		t.Error("Request() returned nil response")
		return
	}

	if response.Type != "test.response" {
		t.Errorf("Response type = %v, want %v", response.Type, "test.response")
	}
}

func TestPublisher_Request_Timeout(t *testing.T) {
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

	err = client.Connect()
	if err != nil || !client.IsConnected() {
		t.Skipf("NATS server not available or not connected: %v", err)
		return
	}
	defer client.Close()

	publisher := NewPublisher(client, "test-service")

	// Request to non-existent subject should timeout
	_, err = publisher.Request(context.Background(), "test.nonexistent", "test.request", map[string]string{"key": "value"}, 100*time.Millisecond)
	if err == nil {
		t.Error("Request() should timeout when no responder exists")
	}
}

func TestPublisher_InvalidData(t *testing.T) {
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

	err = client.Connect()
	if err != nil || !client.IsConnected() {
		t.Skipf("NATS server not available or not connected: %v", err)
		return
	}
	defer client.Close()

	publisher := NewPublisher(client, "test-service")

	// Try to publish unmarshalable data
	err = publisher.Publish(context.Background(), "test.subject", "test.event", make(chan int), nil)
	if err == nil {
		t.Error("Publish() should return error for unmarshalable data")
	}
}
