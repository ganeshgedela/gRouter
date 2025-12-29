package nats

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

func TestJetStream_Integration(t *testing.T) {
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

	js, err := client.JetStream()
	if err != nil {
		t.Fatalf("Failed to get JetStream context: %v", err)
	}

	// Create a stream for testing
	streamName := "TEST_STREAM"
	subject := "test.js.subject"
	_, err = js.AddStream(&nats.StreamConfig{
		Name:     streamName,
		Subjects: []string{subject},
		Storage:  nats.MemoryStorage,
	})
	if err != nil {
		t.Fatalf("Failed to create stream: %v", err)
	}
	defer js.DeleteStream(streamName)

	publisher := NewPublisher(client, "test-service")
	subscriber := NewSubscriber(client, "test-subscriber")

	// Test message reception
	var wg sync.WaitGroup
	var receivedMsg *MessageEnvelope
	wg.Add(1)

	handler := func(ctx context.Context, sub string, msg *MessageEnvelope) error {
		receivedMsg = msg
		wg.Done()
		return nil
	}

	// Subscribe to JetStream
	err = subscriber.SubscribePush(subject, handler, nats.Durable("test-consumer"))
	if err != nil {
		t.Fatalf("SubscribePush() error = %v", err)
	}

	// Publish a message to JetStream
	testData := map[string]string{"key": "js-value"}
	ack, err := publisher.PublishJS(context.Background(), subject, "test.js.event", testData)
	if err != nil {
		t.Fatalf("PublishJS() error = %v", err)
	}
	if ack == nil || ack.Sequence == 0 {
		t.Errorf("Invalid PubAck: %v", ack)
	}

	// Wait for message with timeout
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// Success
	case <-time.After(5 * time.Second):
		t.Fatal("Timeout waiting for JetStream message")
	}

	if receivedMsg == nil {
		t.Error("Handler was not called")
		return
	}

	if receivedMsg.Type != "test.js.event" {
		t.Errorf("Message type = %v, want %v", receivedMsg.Type, "test.js.event")
	}
}

func TestJetStream_Redelivery_Integration(t *testing.T) {
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

	js, err := client.JetStream()
	if err != nil {
		t.Fatalf("Failed to get JetStream context: %v", err)
	}

	// Create a stream for testing
	streamName := "RETRY_STREAM"
	subject := "test.retry.subject"
	_, err = js.AddStream(&nats.StreamConfig{
		Name:     streamName,
		Subjects: []string{subject},
		Storage:  nats.MemoryStorage,
	})
	if err != nil {
		t.Fatalf("Failed to create stream: %v", err)
	}
	defer js.DeleteStream(streamName)

	publisher := NewPublisher(client, "test-service")
	subscriber := NewSubscriber(client, "test-subscriber")

	var mu sync.Mutex
	attempts := 0
	var wg sync.WaitGroup
	wg.Add(2) // We expect 2 attempts (1 fail, 1 success)

	handler := func(ctx context.Context, sub string, msg *MessageEnvelope) error {
		mu.Lock()
		attempts++
		currentAttempt := attempts
		mu.Unlock()

		if currentAttempt == 1 {
			wg.Done()
			return fmt.Errorf("simulated failure")
		}
		if currentAttempt == 2 {
			wg.Done()
			return nil
		}
		return nil
	}

	// Subscribe to JetStream with a short AckWait for faster redelivery
	err = subscriber.SubscribePush(subject, handler,
		nats.Durable("retry-consumer"),
		nats.AckWait(1*time.Second),
		nats.MaxDeliver(3),
	)
	if err != nil {
		t.Fatalf("SubscribePush() error = %v", err)
	}

	// Publish a message
	_, err = publisher.PublishJS(context.Background(), subject, "test.retry.event", map[string]string{"foo": "bar"})
	if err != nil {
		t.Fatalf("PublishJS() error = %v", err)
	}

	// Wait for redelivery
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// Success
	case <-time.After(10 * time.Second):
		t.Fatalf("Timeout waiting for redelivery, attempts = %d", attempts)
	}

	mu.Lock()
	finalAttempts := attempts
	mu.Unlock()

	if finalAttempts < 2 {
		t.Errorf("Expected at least 2 attempts, got %d", finalAttempts)
	}
}
