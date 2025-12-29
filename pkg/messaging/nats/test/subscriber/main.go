package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	messaging "grouter/pkg/messaging/nats"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	// Custom logger config to avoid stack traces/caller info
	zapConfig := zap.NewDevelopmentConfig()
	zapConfig.DisableCaller = true
	zapConfig.DisableStacktrace = true
	logger, _ := zapConfig.Build()
	defer logger.Sync()

	// Parse flags
	queueGroup := flag.String("queue", "", "Queue group name for load balancing")
	maxWorkers := flag.Int("workers", 0, "Max concurrent workers")
	flag.Parse()

	// Load Configuration
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: Config file not found, using defaults: %v", err)
	}

	var rootCfg struct {
		NATS messaging.Config `mapstructure:"nats"`
	}

	// Default values if config missing
	rootCfg.NATS = messaging.Config{
		URL:               "nats://localhost:4222",
		MaxReconnects:     5,
		ReconnectWait:     2 * time.Second,
		ConnectionTimeout: 5 * time.Second,
	}

	if err := viper.Unmarshal(&rootCfg); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	cfg := rootCfg.NATS

	// Create Client
	client, err := messaging.NewNATSClient(cfg, logger)
	if err != nil {
		log.Fatalf("Failed to create NATS client: %v", err)
	}

	if err := client.Connect(); err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer client.Close()

	// Create Subscriber
	sub := messaging.NewSubscriber(client, "test-subscriber")

	// Subscribe to all topics for gRouter
	topic := "gRouter.>"
	logger.Info("Subscribing to topic",
		zap.String("topic", topic),
		zap.String("queue_group", *queueGroup),
		zap.Int("max_workers", *maxWorkers),
	)

	// Create handler with dependencies
	handler := &Handler{
		client: client,
		logger: logger,
	}

	opts := &messaging.SubscribeOptions{
		QueueGroup: *queueGroup,
		MaxWorkers: *maxWorkers,
	}

	err = sub.Subscribe(topic, handler.HandleMessage, opts)

	if err != nil {
		log.Fatalf("Failed to subscribe: %v", err)
	}

	logger.Info("Subscriber running. Press Ctrl+C to stop.")

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	logger.Info("Shutting down...")
}

// Handler encapsulates message handling logic and dependencies
type Handler struct {
	client *messaging.Client
	logger *zap.Logger
}

// HandleMessage processes incoming NATS messages
func (h *Handler) HandleMessage(ctx context.Context, subject string, env *messaging.MessageEnvelope) error {
	h.logger.Info("Received message",
		zap.String("subject", subject),
		zap.String("type", env.Type),
		zap.String("id", env.ID),
		zap.Any("data", env.Data),
	)

	// Check if it's a request and reply
	if env.Reply != "" {
		h.logger.Info("Received request, sending reply", zap.String("reply_to", env.Reply))
		// Echo back
		responseData := map[string]string{"reply": "echo response"}
		dataBytes, _ := json.Marshal(responseData)

		responseEnvelope := &messaging.MessageEnvelope{
			ID:        "response-id",
			Type:      "echo.response",
			Source:    "test-subscriber",
			Timestamp: time.Now(),
			Data:      dataBytes,
		}
		respBytes, _ := json.Marshal(responseEnvelope)

		// Using client connection directly
		if err := h.client.Conn().Publish(env.Reply, respBytes); err != nil {
			h.logger.Error("Failed to reply", zap.Error(err))
		} else {
			h.logger.Info("Reply sent")
		}
	}

	return nil
}
