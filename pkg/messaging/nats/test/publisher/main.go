package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	messaging "grouter/pkg/messaging/nats"

	"go.uber.org/zap"
)

func main() {
	// Custom logger config to avoid stack traces/caller info
	zapConfig := zap.NewDevelopmentConfig()
	zapConfig.DisableCaller = true
	zapConfig.DisableStacktrace = true
	logger, _ := zapConfig.Build()
	defer logger.Sync()

	// Configuration
	// Assuming NATS is running on localhost default port
	cfg := messaging.Config{
		URL:               "nats://localhost:4222",
		MaxReconnects:     5,
		ReconnectWait:     2 * time.Second,
		ConnectionTimeout: 5 * time.Second,
	}

	// Create Client
	client, err := messaging.NewNATSClient(cfg, logger)
	if err != nil {
		log.Fatalf("Failed to create NATS client: %v", err)
	}

	if err := client.Connect(); err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer client.Close()

	// Create Publisher
	pub := messaging.NewPublisher(client, "test-publisher")

	// Parse flags
	count := flag.Int("count", 1, "Number of messages to publish")
	delay := flag.Duration("delay", 100*time.Millisecond, "Delay between messages")
	sSubject := flag.String("subject", "gRouter.natsdemosvc.echo", "NATS subject to publish to")
	sType := flag.String("type", "echo.request", "Message type")
	sData := flag.String("data", "{}", "JSON data payload")
	flag.Parse()

	// Topic and Payload
	topic := *sSubject
	msgType := *sType

	for i := 0; i < *count; i++ {
		var payload interface{}
		// Basic string payload if not valid JSON, or pass generic map
		payload = map[string]interface{}{
			"message": fmt.Sprintf("Hello from manual publisher! #%d", i+1),
			"data":    *sData,
			"time":    time.Now().Format(time.RFC3339),
		}

		// Publish/Request
		logger.Info("Sending request",
			zap.Int("seq", i+1),
			zap.String("topic", topic),
			zap.Any("payload", payload),
		)

		// Use Request to wait for reply (health checks need this)
		response, err := pub.Request(context.Background(), topic, msgType, payload, 2*time.Second)
		if err != nil {
			logger.Error("Request failed", zap.Error(err))
		} else {
			logger.Info("Received response",
				zap.String("id", response.ID),
				zap.Any("data", string(response.Data)), // Log data as string for readability
			)
		}

		if *count > 1 {
			time.Sleep(*delay)
		}
	}

	// Wait a bit to ensure logs/metrics are flushed if needed
	time.Sleep(500 * time.Millisecond)
}
