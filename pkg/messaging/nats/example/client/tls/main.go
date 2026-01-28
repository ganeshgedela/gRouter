package main

import (
	"log"
	"time"

	"grouter/pkg/config"
	"grouter/pkg/messaging/nats"

	"go.uber.org/zap"
)

// This example demonstrates connecting to NATS using Mutual TLS (mTLS).
//
// To run this example, you need certificates generated (CA, Server, Client).
// The NATS server must be running with TLS enabled and verifying client certs.

func main() {
	logger, _ := zap.NewDevelopment()

	cfg := &config.Config{
		App: config.AppConfig{
			Name: "tls-auth-client",
		},
		NATS: config.NATSConfig{
			// Note: Use 'tls://' scheme if you want to force TLS from the start,
			// though 'nats://' with TLS options usually works if connection upgrades.
			URL: "tls://localhost:4222",
			TLS: config.TLSConfig{
				Enabled:  true,
				CertFile: "assets/client-cert.pem",
				KeyFile:  "assets/client-key.pem",
				CAFile:   "assets/ca.pem",
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
		logger.Info("Successfully connected with Mutual TLS!")
	}
}
