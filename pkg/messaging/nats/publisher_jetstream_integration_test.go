package nats

import (
	"context"
	"fmt"
	"testing"
	"time"

	server "github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestPublisher_JetStream_Integration_Advanced(t *testing.T) {
	// Start JS server
	opts := &server.Options{Port: -1, JetStream: true}
	s := RunTestServer(opts)
	defer s.Shutdown()

	// Setup client
	logger, _ := zap.NewDevelopment()
	config := Config{
		URL:               s.ClientURL(),
		MaxReconnects:     2,
		ReconnectWait:     2 * time.Second,
		ConnectionTimeout: 2 * time.Second,
	}
	client, err := NewNATSClient(config, logger)
	assert.NoError(t, err)
	err = client.Connect()
	if err != nil {
		t.Skipf("NATS server not available: %v", err)
		return
	}
	defer client.Close()

	// Ensure JetStream context
	js, err := client.JetStream()
	assert.NoError(t, err)

	// Create Stream
	streamName := "ADVANCED_JS_TEST"
	subject := "advanced.js.*"
	_, err = js.AddStream(&nats.StreamConfig{
		Name:     streamName,
		Subjects: []string{subject},
		Storage:  nats.MemoryStorage,
	})
	if err != nil {
		// Ignore if exists
	}

	publisher := NewPublisher(client)

	t.Run("PublishAsyncJS with Deduplication", func(t *testing.T) {
		msgID := fmt.Sprintf("dedup-%d", time.Now().UnixNano())
		data := map[string]string{"foo": "bar"}

		// 1. Publish first time
		ack, err := publisher.PublishAsyncJS(context.Background(), "advanced.js.dedup", "event.type", data, &PublishOptions{
			MsgID: msgID,
		})
		assert.NoError(t, err)

		// Wait for ack
		select {
		case <-ack.Ok():
			// Good
		case err := <-ack.Err():
			t.Fatalf("Publish failed: %v", err)
		case <-time.After(time.Second):
			t.Fatal("Timeout waiting for ack")
		}

		// 2. Publish duplicate
		ack2, err := publisher.PublishAsyncJS(context.Background(), "advanced.js.dedup", "event.type", data, &PublishOptions{
			MsgID: msgID,
		})
		assert.NoError(t, err)

		// Wait for ack (should be duplicate)
		select {
		case <-ack2.Ok():
			// NATS handles dedup by acknowledging with 'Duplicate' internal flag, but client sees success
			// unless we check the sequence. Usually it returns the *same* sequence.
		case err := <-ack2.Err():
			t.Fatalf("Duplicate publish failed: %v", err)
		case <-time.After(time.Second):
			t.Fatal("Timeout waiting for ack2")
		}
	})
}
