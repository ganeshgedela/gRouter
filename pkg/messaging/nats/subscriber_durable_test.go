package nats

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

func TestSubscriber_DurableConsumer(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	logger, _ := zap.NewDevelopment()
	config := Config{
		URL:               GetTestNATSURL(nil),
		MaxReconnects:     5,
		ReconnectWait:     1 * time.Second,
		ConnectionTimeout: 2 * time.Second,
	}

	client, err := NewNATSClient(config, logger)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	err = client.Connect()
	if err != nil || !client.IsConnected() {
		t.Skipf("NATS server not available: %v", err)
		return
	}
	defer client.Close()

	js, err := client.JetStream()
	if err != nil {
		t.Fatalf("Failed to get JetStream context: %v", err)
	}

	// Create Stream
	streamName := "TEST_DURABLE_STREAM"
	subject := "test.durable"
	// Ensure stream clean state
	_ = js.DeleteStream(streamName)

	_, err = js.AddStream(&nats.StreamConfig{
		Name:     streamName,
		Subjects: []string{subject},
	})
	if err != nil {
		t.Fatalf("Failed to create stream: %v", err)
	}
	defer js.DeleteStream(streamName)

	subscriber := NewSubscriber(client)
	handler := func(ctx context.Context, subject string, msg *MessageEnvelope) error {
		return nil
	}

	durableName := "test-durable-consumer"
	opts := &SubscribeOptions{
		Durable:       durableName,
		AckWait:       5 * time.Second,
		MaxDeliveries: 3,
	}

	sub, err := subscriber.SubscribePush(context.Background(), subject, handler, opts)
	if err != nil {
		t.Fatalf("SubscribePush failed: %v", err)
	}
	defer sub.Unsubscribe()

	// Verify Consumer Info
	ci, err := js.ConsumerInfo(streamName, durableName)
	if err != nil {
		t.Fatalf("Failed to get ConsumerInfo: %v", err)
	}

	if ci.Config.Durable != durableName {
		t.Errorf("Durable name = %v, want %v", ci.Config.Durable, durableName)
	}
	if ci.Config.AckWait != 5*time.Second {
		t.Errorf("AckWait = %v, want %v", ci.Config.AckWait, 5*time.Second)
	}
	if ci.Config.MaxDeliver != 3 {
		t.Errorf("MaxDeliver = %v, want %v", ci.Config.MaxDeliver, 3)
	}
}

func TestSubscriber_DLQ(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	logger, _ := zap.NewDevelopment()
	config := Config{
		URL:               GetTestNATSURL(nil),
		MaxReconnects:     5,
		ReconnectWait:     1 * time.Second,
		ConnectionTimeout: 2 * time.Second,
	}

	client, err := NewNATSClient(config, logger)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	err = client.Connect()
	if err != nil || !client.IsConnected() {
		t.Skipf("NATS server not available: %v", err)
		return
	}
	defer client.Close()

	js, err := client.JetStream()
	if err != nil {
		t.Fatalf("Failed to get JetStream context: %v", err)
	}

	// Setup Streams
	streamName := "TEST_DLQ_STREAM"
	subject := "test.dlq.source"
	dlqSubject := "test.dlq.target"

	_ = js.DeleteStream(streamName)
	_, err = js.AddStream(&nats.StreamConfig{
		Name:     streamName,
		Subjects: []string{subject},
	})
	if err != nil {
		t.Fatalf("Failed to create stream: %v", err)
	}
	defer js.DeleteStream(streamName)

	subscriber := NewSubscriber(client)
	// Handler that ALWAYS fails
	handler := func(ctx context.Context, subject string, msg *MessageEnvelope) error {
		return fmt.Errorf("poison message")
	}

	// Create a subscriber for the DLQ subject to verify receipt
	dlqReceived := make(chan struct{})
	_, err = client.Conn().Subscribe(dlqSubject, func(msg *nats.Msg) {
		close(dlqReceived)
	})
	if err != nil {
		t.Fatalf("Failed to subscribe to DLQ subject: %v", err)
	}

	// Subscribe with DLQ options
	opts := &SubscribeOptions{
		Durable:       "test-dlq-consumer",
		AckWait:       1 * time.Second,
		MaxDeliveries: 3,
		DLQSubject:    dlqSubject,
	}

	sub, err := subscriber.SubscribePush(context.Background(), subject, handler, opts)
	if err != nil {
		t.Fatalf("SubscribePush failed: %v", err)
	}
	defer sub.Unsubscribe()

	// Publish a message
	// We need to construct a valid envelope manually since we are using raw publish or use Publisher
	// Let's use raw publish with a manually created envelope JSON
	env := MessageEnvelope{
		ID:        "poison-id",
		Type:      "poison.pill",
		Source:    "test",
		Timestamp: time.Now(),
		Data:      json.RawMessage(`{"foo":"bar"}`),
	}
	bytes, _ := json.Marshal(env)

	_, err = js.Publish(subject, bytes)
	if err != nil {
		t.Fatalf("Failed to publish: %v", err)
	}

	// Wait for DLQ
	select {
	case <-dlqReceived:
		// Success!
	case <-time.After(10 * time.Second): // Needs to retry 3 times (immediate Nak handling)
		t.Fatal("Timeout waiting for message to arrive in DLQ")
	}
}
