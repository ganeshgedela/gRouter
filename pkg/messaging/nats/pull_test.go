package nats

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

func TestSubscriber_Pull_Integration(t *testing.T) {
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

	// Ensure JetStream is enabled
	js, err := client.JetStream()
	if err != nil {
		t.Fatalf("Failed to get JetStream context: %v", err)
	}

	// Create a stream for testing
	streamName := "TEST_PULL_STREAM"
	subject := "test.pull.>"
	_, err = js.AddStream(&nats.StreamConfig{
		Name:     streamName,
		Subjects: []string{subject},
	})
	if err != nil {
		t.Fatalf("Failed to create stream: %v", err)
	}
	defer js.DeleteStream(streamName)

	publisher := NewPublisher(client, "test-service")
	subscriber := NewSubscriber(client, "test-service")
	defer subscriber.Close()

	// Publish some messages
	for i := 0; i < 5; i++ {
		err := publisher.Publish(context.Background(), "test.pull.event", "test.event", map[string]int{"id": i}, nil)
		if err != nil {
			t.Fatalf("Failed to publish message: %v", err)
		}
	}

	// Subscribe using Pull Consumer
	received := make(chan int, 5)
	err = subscriber.SubscribePull("test.pull.event", "test-durable", func(ctx context.Context, subject string, msg *MessageEnvelope) error {
		var data map[string]int
		if err := json.Unmarshal(msg.Data, &data); err != nil {
			return err
		}
		received <- data["id"]
		return nil
	}, WithBatchSize(2), WithFetchTimeout(1*time.Second))

	if err != nil {
		t.Fatalf("SubscribePull failed: %v", err)
	}

	// Verify messages received
	timeout := time.After(5 * time.Second)
	count := 0
	for count < 5 {
		select {
		case <-received:
			count++
		case <-timeout:
			t.Fatalf("Timed out waiting for messages, received %d/5", count)
		}
	}
}
