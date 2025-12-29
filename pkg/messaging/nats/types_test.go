package nats

import (
	"encoding/json"
	"testing"
	"time"
)

func TestMessageEnvelope_Marshal(t *testing.T) {
	data := map[string]string{"key": "value"}
	dataBytes, _ := json.Marshal(data)

	envelope := MessageEnvelope{
		ID:        "test-id-123",
		Type:      "test.event",
		Timestamp: time.Now(),
		Source:    "test-service",
		Data:      dataBytes,
	}

	// Marshal envelope
	envelopeBytes, err := json.Marshal(envelope)
	if err != nil {
		t.Fatalf("Failed to marshal envelope: %v", err)
	}

	// Unmarshal envelope
	var decoded MessageEnvelope
	err = json.Unmarshal(envelopeBytes, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal envelope: %v", err)
	}

	// Verify fields
	if decoded.ID != envelope.ID {
		t.Errorf("ID = %v, want %v", decoded.ID, envelope.ID)
	}

	if decoded.Type != envelope.Type {
		t.Errorf("Type = %v, want %v", decoded.Type, envelope.Type)
	}

	if decoded.Source != envelope.Source {
		t.Errorf("Source = %v, want %v", decoded.Source, envelope.Source)
	}

	// Verify data
	var decodedData map[string]string
	err = json.Unmarshal(decoded.Data, &decodedData)
	if err != nil {
		t.Fatalf("Failed to unmarshal data: %v", err)
	}

	if decodedData["key"] != "value" {
		t.Errorf("Data key = %v, want %v", decodedData["key"], "value")
	}
}

func TestPublishOptions(t *testing.T) {
	tests := []struct {
		name string
		opts PublishOptions
	}{
		{
			name: "sync publish",
			opts: PublishOptions{
				Async:   false,
				Timeout: 5 * time.Second,
			},
		},
		{
			name: "async publish",
			opts: PublishOptions{
				Async:   true,
				Timeout: 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just verify the struct can be created
			if tt.opts.Async && tt.name != "async publish" {
				t.Error("Async flag mismatch")
			}
		})
	}
}

func TestSubscribeOptions(t *testing.T) {
	tests := []struct {
		name string
		opts SubscribeOptions
	}{
		{
			name: "no queue group",
			opts: SubscribeOptions{
				QueueGroup: "",
				MaxWorkers: 1,
			},
		},
		{
			name: "with queue group",
			opts: SubscribeOptions{
				QueueGroup: "test-queue",
				MaxWorkers: 5,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just verify the struct can be created
			if tt.opts.QueueGroup == "" && tt.name == "with queue group" {
				t.Error("QueueGroup mismatch")
			}
		})
	}
}
