package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"grouter/pkg/config"
	messaging "grouter/pkg/messaging/nats"
	"grouter/pkg/telemetry"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
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
		NATS    messaging.Config     `mapstructure:"nats"`
		Tracing config.TracingConfig `mapstructure:"tracing"`
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

	// Create Tracing Config locally for this test
	tracingCfg := rootCfg.Tracing
	if tracingCfg.ServiceName == "" {
		tracingCfg.ServiceName = "nats-subscriber"
	}
	if tracingCfg.Exporter == "" {
		tracingCfg.Exporter = "stdout"
	}

	// Initialize Tracer
	shutdown, err := telemetry.InitTracer(tracingCfg)
	if err != nil {
		log.Printf("Failed to init tracer: %v", err)
	}
	defer shutdown(context.Background())

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
	// Use Logging Middleware
	sub.Use(messaging.LoggingMiddleware(logger))

	// Use Metrics Middleware (if enabled in config, though we force enabled in this test logic for now or check rootCfg if available)
	// We will just enable it to test
	sub.Use(messaging.MetricsMiddleware())

	if cfg.Tracing.Enabled {
		tracer := otel.Tracer("nats-subscriber")
		sub.Use(messaging.TracingMiddleware(tracer))
	}

	// Subscribe to all topics for gRouter
	topic := "gRouter.>"
	logger.Info("Subscribing to topic",
		zap.String("topic", topic),
		zap.String("queue_group", *queueGroup),
		zap.Int("max_workers", *maxWorkers),
	)

	// Start Metrics Server
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		logger.Info("Metrics server starting on :8081")
		if err := http.ListenAndServe(":8081", nil); err != nil {
			logger.Error("Metrics server failed", zap.Error(err))
		}
	}()

	// Create Publisher for replies
	pub := messaging.NewPublisher(client, "test-subscriber")

	// Create handler with dependencies
	handler := &Handler{
		client:    client,
		logger:    logger,
		publisher: pub,
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
	client    *messaging.Client
	logger    *zap.Logger
	publisher messaging.Publisher
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

		// Use the Publisher interface to send the reply
		// We use Publish (Sync) or PublishAsync depending on need.
		// Since it is a reply, we usually want it to go out quickly.
		// Note: The subject is env.Reply.
		if err := h.publisher.Publish(ctx, env.Reply, "echo.response", responseData, nil); err != nil {
			h.logger.Error("Failed to reply", zap.Error(err))
		} else {
			h.logger.Info("Reply sent")
		}
	}

	return nil
}
