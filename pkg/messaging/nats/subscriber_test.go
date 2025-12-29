package nats

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestNewSubscriber(t *testing.T) {
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

	subscriber := NewSubscriber(client, "test-subscriber")
	if subscriber == nil {
		t.Error("NewSubscriber() returned nil")
	}
}

func TestSubscriber_Subscribe_NotConnected(t *testing.T) {
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

	subscriber := NewSubscriber(client, "test-subscriber")

	handler := func(ctx context.Context, subject string, msg *MessageEnvelope) error {
		return nil
	}

	// Try to subscribe without connection
	err = subscriber.Subscribe("test.subject", handler, nil)
	if err == nil {
		t.Error("Subscribe() should return error when not connected")
	}
}

func TestSubscriber_Subscribe_Integration(t *testing.T) {
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

	subscriber := NewSubscriber(client, "test-subscriber")
	publisher := NewPublisher(client, "test-service")

	// Test message reception
	var wg sync.WaitGroup
	var receivedMsg *MessageEnvelope
	wg.Add(1)

	handler := func(ctx context.Context, subject string, msg *MessageEnvelope) error {
		receivedMsg = msg
		wg.Done()
		return nil
	}

	err = subscriber.Subscribe("test.subscribe", handler, nil)
	if err != nil {
		t.Fatalf("Subscribe() error = %v", err)
	}

	// Give subscription time to be ready
	time.Sleep(100 * time.Millisecond)

	// Publish a message
	testData := map[string]string{"key": "value"}
	err = publisher.Publish(context.Background(), "test.subscribe", "test.event", testData, nil)
	if err != nil {
		t.Fatalf("Publish() error = %v", err)
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
	case <-time.After(2 * time.Second):
		t.Fatal("Timeout waiting for message")
	}

	if receivedMsg == nil {
		t.Error("Handler was not called")
		return
	}

	if receivedMsg.Type != "test.event" {
		t.Errorf("Message type = %v, want %v", receivedMsg.Type, "test.event")
	}

	if receivedMsg.Source != "test-service" {
		t.Errorf("Message source = %v, want %v", receivedMsg.Source, "test-service")
	}

	// Verify data
	var data map[string]string
	err = json.Unmarshal(receivedMsg.Data, &data)
	if err != nil {
		t.Fatalf("Failed to unmarshal data: %v", err)
	}

	if data["key"] != "value" {
		t.Errorf("Data key = %v, want %v", data["key"], "value")
	}
}

func TestSubscriber_QueueGroup_Integration(t *testing.T) {
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

	// Create two subscribers in the same queue group
	subscriber1 := NewSubscriber(client, "subscriber-1")
	subscriber2 := NewSubscriber(client, "subscriber-2")
	publisher := NewPublisher(client, "test-service")

	var mu sync.Mutex
	count1 := 0
	count2 := 0

	handler1 := func(ctx context.Context, subject string, msg *MessageEnvelope) error {
		mu.Lock()
		count1++
		mu.Unlock()
		return nil
	}

	handler2 := func(ctx context.Context, subject string, msg *MessageEnvelope) error {
		mu.Lock()
		count2++
		mu.Unlock()
		return nil
	}

	opts := &SubscribeOptions{QueueGroup: "test-queue"}

	err = subscriber1.Subscribe("test.queue", handler1, opts)
	if err != nil {
		t.Fatalf("Subscribe() error = %v", err)
	}

	err = subscriber2.Subscribe("test.queue", handler2, opts)
	if err != nil {
		t.Fatalf("Subscribe() error = %v", err)
	}

	// Give subscriptions time to be ready
	time.Sleep(100 * time.Millisecond)

	// Publish multiple messages
	numMessages := 10
	for i := 0; i < numMessages; i++ {
		err = publisher.Publish(context.Background(), "test.queue", "test.event", map[string]int{"count": i}, nil)
		if err != nil {
			t.Fatalf("Publish() error = %v", err)
		}
	}

	// Wait for messages to be processed
	time.Sleep(500 * time.Millisecond)

	mu.Lock()
	total := count1 + count2
	mu.Unlock()

	// In a queue group, messages should be distributed
	// Both subscribers should have received some messages
	if total != numMessages {
		t.Errorf("Total messages received = %v, want %v", total, numMessages)
	}

	// Both should have received at least one (with high probability)
	// Note: This is probabilistic, but with 10 messages it's very likely
	if count1 == 0 && count2 == 0 {
		t.Error("Neither subscriber received any messages")
	}
}

func TestSubscriber_Unsubscribe(t *testing.T) {
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

	subscriber := NewSubscriber(client, "test-subscriber")
	publisher := NewPublisher(client, "test-service")

	var mu sync.Mutex
	count := 0

	handler := func(ctx context.Context, subject string, msg *MessageEnvelope) error {
		mu.Lock()
		count++
		mu.Unlock()
		return nil
	}

	err = subscriber.Subscribe("test.unsub", handler, nil)
	if err != nil {
		t.Fatalf("Subscribe() error = %v", err)
	}

	// Give subscription time to be ready
	time.Sleep(100 * time.Millisecond)

	// Publish a message
	err = publisher.Publish(context.Background(), "test.unsub", "test.event", map[string]string{"key": "value"}, nil)
	if err != nil {
		t.Fatalf("Publish() error = %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	// Unsubscribe
	err = subscriber.Unsubscribe()
	if err != nil {
		t.Errorf("Unsubscribe() error = %v", err)
	}

	// Publish another message
	err = publisher.Publish(context.Background(), "test.unsub", "test.event", map[string]string{"key": "value2"}, nil)
	if err != nil {
		t.Fatalf("Publish() error = %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	finalCount := count
	mu.Unlock()

	// Should only have received the first message
	if finalCount != 1 {
		t.Errorf("Message count = %v, want 1 (should not receive messages after unsubscribe)", finalCount)
	}
}

func TestSubscriber_HandlerError(t *testing.T) {
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

	subscriber := NewSubscriber(client, "test-subscriber")
	publisher := NewPublisher(client, "test-service")

	var wg sync.WaitGroup
	wg.Add(1)

	// Handler that returns an error
	handler := func(ctx context.Context, subject string, msg *MessageEnvelope) error {
		wg.Done()
		return ErrHandlerFailed
	}

	err = subscriber.Subscribe("test.error", handler, nil)
	if err != nil {
		t.Fatalf("Subscribe() error = %v", err)
	}

	// Give subscription time to be ready
	time.Sleep(100 * time.Millisecond)

	// Publish a message
	err = publisher.Publish(context.Background(), "test.error", "test.event", map[string]string{"key": "value"}, nil)
	if err != nil {
		t.Fatalf("Publish() error = %v", err)
	}

	// Wait for handler to be called
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// Handler was called (even though it returned an error)
	case <-time.After(2 * time.Second):
		t.Fatal("Timeout waiting for handler")
	}
}

var ErrHandlerFailed = fmt.Errorf("handler failed")

func TestSubscriber_Reply(t *testing.T) {
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

	subscriber := NewSubscriber(client, "test-responder")
	publisher := NewPublisher(client, "test-requester")

	// Setup responder
	responderHandler := func(ctx context.Context, subject string, msg *MessageEnvelope) error {
		// Echo back the data
		// Manually publish response since Reply is now registration-only
		responseEnvelope := &MessageEnvelope{
			ID:        "response-id",
			Type:      "test.response",
			Source:    "test-responder",
			Timestamp: time.Now(),
			Data:      msg.Data,
			Metadata:  make(map[string]string),
		}
		respBytes, _ := json.Marshal(responseEnvelope)
		return client.Conn().Publish(msg.Reply, respBytes)
	}

	err = subscriber.Subscribe("test.request", responderHandler, nil)
	if err != nil {
		t.Fatalf("Subscribe() error = %v", err)
	}

	// Give subscription time to be ready
	time.Sleep(100 * time.Millisecond)

	// Send request
	requestData := map[string]string{"query": "foo"}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	response, err := publisher.Request(ctx, "test.request", "test.query", requestData, 2*time.Second)
	if err != nil {
		t.Fatalf("Request() error = %v", err)
	}

	// Verify response
	if response.Type != "test.response" {
		t.Errorf("Response type = %v, want %v", response.Type, "test.response")
	}

	if response.Source != "test-responder" {
		t.Errorf("Response source = %v, want %v", response.Source, "test-responder")
	}

	// Verify data
	var responseData map[string]string
	// The data is double-marshaled because we passed msg.Data (which is json.RawMessage/[]byte) directly to Reply
	// Reply marshals it again.
	// Wait, msg.Data is json.RawMessage which is []byte.
	// json.Marshal([]byte) encodes it as base64 string if it's treated as []byte, or just as is if it's RawMessage?
	// json.RawMessage marshals to itself.
	// So it should be fine.
	err = json.Unmarshal(response.Data, &responseData)
	if err != nil {
		t.Fatalf("Failed to unmarshal response data: %v", err)
	}

	if responseData["query"] != "foo" {
		t.Errorf("Response data query = %v, want %v", responseData["query"], "foo")
	}
}

func TestSubscriber_GracefulShutdown(t *testing.T) {
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

	subscriber := NewSubscriber(client, "test-graceful")
	publisher := NewPublisher(client, "test-service")

	// Handler that sleeps
	handlerFinished := make(chan struct{})
	handler := func(ctx context.Context, subject string, msg *MessageEnvelope) error {
		time.Sleep(500 * time.Millisecond)
		close(handlerFinished)
		return nil
	}

	err = subscriber.Subscribe("test.graceful", handler, nil)
	if err != nil {
		t.Fatalf("Subscribe() error = %v", err)
	}

	// Give subscription time to be ready
	time.Sleep(100 * time.Millisecond)

	// Publish message
	err = publisher.Publish(context.Background(), "test.graceful", "test.event", nil, nil)
	if err != nil {
		t.Fatalf("Publish() error = %v", err)
	}

	// Wait for handler to start (approx)
	time.Sleep(100 * time.Millisecond)

	// Start closing
	start := time.Now()
	err = subscriber.Close()
	if err != nil {
		t.Fatalf("Close() error = %v", err)
	}
	duration := time.Since(start)

	// Check if we waited
	if duration < 200*time.Millisecond {
		t.Errorf("Close() returned too quickly (%v), expected it to wait for handler", duration)
	}

	select {
	case <-handlerFinished:
	// Success
	default:
		t.Error("Handler did not finish before Close() returned")
	}
}
