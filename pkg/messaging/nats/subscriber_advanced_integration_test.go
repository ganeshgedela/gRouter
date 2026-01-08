package nats

import (
	"context"
	"sync"
	"testing"
	"time"

	server "github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestSubscriber_Advanced_Integration(t *testing.T) {
	// Start embedded server
	s := RunTestServer(&server.Options{Port: -1, JetStream: true})
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

	js, err := client.JetStream()
	assert.NoError(t, err)

	// Create Stream
	streamName := "ADVANCED_SUB_TEST"
	subject := "advanced.sub.*"
	_, err = js.AddStream(&nats.StreamConfig{
		Name:     streamName,
		Subjects: []string{subject},
		Storage:  nats.MemoryStorage,
	})
	if err != nil {
		t.Logf("Warning: AddStream failed: %v", err)
	}

	subscriber := NewSubscriber(client)
	publisher := NewPublisher(client)

	t.Run("Pull Consumer Subscription", func(t *testing.T) {
		// Publish messages first
		for i := 0; i < 5; i++ {
			_, err := publisher.PublishAsyncJS(context.Background(), "advanced.sub.pull", "event.pull", map[string]int{"seq": i}, nil)
			assert.NoError(t, err)
		}

		// Subscribe Pull
		// SubscribePull takes (ctx, subject, handler, opts)
		// It manages the pull loop internally if handler is provided.
		// We expect the handler to be called.

		wg := sync.WaitGroup{}
		wg.Add(5)

		handler := func(ctx context.Context, subject string, msg *MessageEnvelope) error {
			wg.Done()
			return nil
		}

		opts := &SubscribeOptions{
			Durable: "TEST_PULL_DURABLE",
		}

		err := subscriber.SubscribePull(context.Background(), "advanced.sub.pull", handler, opts)
		if err != nil {
			t.Fatalf("SubscribePull failed: %v", err)
		}

		// Wait for handler to process messages
		done := make(chan struct{})
		go func() {
			wg.Wait()
			close(done)
		}()

		select {
		case <-done:
			// Success
		case <-time.After(5 * time.Second):
			t.Fatal("Timeout waiting for pulled messages")
		}
	})

	t.Run("Queue Group Load Balancing (Push)", func(t *testing.T) {
		count := 10
		wg := sync.WaitGroup{}
		wg.Add(count)

		// Create two subscribers in same queue group
		// Handler handles message
		handler := func(ctx context.Context, subject string, msg *MessageEnvelope) error {
			wg.Done()
			return nil
		}

		// Subscribe with queue group in options
		opts := &SubscribeOptions{
			QueueGroup: "test-lb-group",
		}

		sub1, err := subscriber.Subscribe(context.Background(), "advanced.sub.lb", handler, opts)
		assert.NoError(t, err)
		defer sub1.Unsubscribe()

		sub2, err := subscriber.Subscribe(context.Background(), "advanced.sub.lb", handler, opts)
		assert.NoError(t, err)
		defer sub2.Unsubscribe()

		// Publish messages
		for i := 0; i < count; i++ {
			err := publisher.Publish(context.Background(), "advanced.sub.lb", "event.lb", nil, nil)
			assert.NoError(t, err)
		}

		// Wait
		done := make(chan struct{})
		go func() {
			wg.Wait()
			close(done)
		}()

		select {
		case <-done:
			// Success
		case <-time.After(3 * time.Second):
			t.Fatal("Timeout waiting for messages")
		}
	})
}
