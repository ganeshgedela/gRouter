package main

import (
	"log"
	"time"

	"grouter/pkg/config"
	"grouter/pkg/messaging/nats"

	"go.uber.org/zap"
)

// This example demonstrates connecting to NATS using a Credentials File (NKEYs/JWT).
// This is the standard for NATS 2.0+ Authentication/Authorization (Account security).
//
// To run this example:
// 1. Generate an operator, account, and user using `nsc`.
// 2. Export the user credentials file: `nsc generate creds -o user.creds`
// 3. Point `Auth.CredsFile` to that file path.

func main() {
	logger, _ := zap.NewDevelopment()

	cfg := &config.Config{
		App: config.AppConfig{
			Name: "creds-auth-client",
		},
		NATS: config.NATSConfig{
			URL: "nats://localhost:4222",
			Auth: config.AuthConfig{
				CredsFile: "assets/user.creds", // Path relative to CWD (setup_and_verify script dir)
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
		logger.Info("Successfully connected with Credentials File!")
	}
}
