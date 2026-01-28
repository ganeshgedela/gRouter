package main

import (
	"log"
	"time"

	"grouter/pkg/config"
	"grouter/pkg/messaging/nats"

	"go.uber.org/zap"
)

// This example demonstrates connecting to NATS using a simple Auth Token.
//
// To run this example, you need a NATS server configured with a token:
// nats-server -auth "my-secret-token"

func main() {
	// Initialize logger
	logger, _ := zap.NewDevelopment()

	// Configuration
	cfg := &config.Config{
		App: config.AppConfig{
			Name: "token-auth-client",
		},
		NATS: config.NATSConfig{
			URL: "nats://localhost:4222",
			Auth: config.AuthConfig{
				Token: "my-secret-token", // Set the token here
			},
			ConnectionTimeout: 5 * time.Second,
			MaxReconnects:     5,
			ReconnectWait:     2 * time.Second,
		},
	}

	// Create Messenger
	messenger := &nats.Messenger{}
	if err := messenger.InitFromConfig(cfg, logger); err != nil {
		log.Fatalf("Failed to initialize messenger: %v", err)
	}
	defer messenger.Close()

	if messenger.IsConnected() {
		logger.Info("Successfully connected with Token Auth!")

		// Example usage
		// messenger.Publish(context.Background(), "foo", "test", "hello", nil)
	} else {
		logger.Error("Failed to connect")
	}
}
