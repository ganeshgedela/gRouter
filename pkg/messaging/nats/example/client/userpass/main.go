package main

import (
	"log"
	"time"

	"grouter/pkg/config"
	"grouter/pkg/messaging/nats"

	"go.uber.org/zap"
)

// This example demonstrates connecting to NATS using Username and Password.
//
// To run this example, you need a NATS server configured with user/pass:
// nats-server --user myuser --pass mypassword

func main() {
	logger, _ := zap.NewDevelopment()

	cfg := &config.Config{
		App: config.AppConfig{
			Name: "user-pass-client",
		},
		NATS: config.NATSConfig{
			URL: "nats://localhost:4222",
			Auth: config.AuthConfig{
				Username: "myuser",     // Username
				Password: "mypassword", // Password
			},
			ConnectionTimeout: 5 * time.Second,
			MaxReconnects:     5,
		},
	}

	messenger := &nats.Messenger{}
	if err := messenger.InitFromConfig(cfg, logger); err != nil {
		log.Fatalf("Failed to initialize messenger: %v", err)
	}
	defer messenger.Close()

	if messenger.IsConnected() {
		logger.Info("Successfully connected with User/Pass Auth!")
	}
}
