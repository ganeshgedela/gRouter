package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
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
	count := flag.Int("count", 1, "Number of messages to publish")
	delay := flag.Duration("delay", 100*time.Millisecond, "Delay between messages")
	sSubject := flag.String("subject", "gRouter.natsdemosvc.echo", "NATS subject to publish to")
	sType := flag.String("type", "echo.request", "Message type")
	sData := flag.String("data", "{}", "JSON data payload")
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

	// Create Tracing Config
	// Use config from file, but default service name if missing
	tracingCfg := rootCfg.Tracing
	if tracingCfg.ServiceName == "" {
		tracingCfg.ServiceName = "nats-publisher"
	}
	// Fallback to stdout if not specified in file, though Config should handle it.
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

	// Create Publisher
	pub := messaging.NewPublisher(client, "test-publisher")

	// Use Logging Middleware
	pub.Use(messaging.PublisherLoggingMiddleware(logger))
	pub.UseRequest(messaging.RequestLoggingMiddleware(logger))

	// Use Metrics Middleware
	pub.Use(messaging.PublisherMetricsMiddleware())
	pub.UseRequest(messaging.RequestMetricsMiddleware())

	// Use Tracing Middleware
	if cfg.Tracing.Enabled {
		tracer := otel.Tracer("nats-publisher")
		pub.Use(messaging.PublisherTracingMiddleware(tracer))
		pub.UseRequest(messaging.RequestTracingMiddleware(tracer))

		// My earlier fix in messenger.go used:
		// m.Publisher.UseRequest(RequestLoggingMiddleware(logger))
		// TracingMiddleware for Request?
		// In messenger.go: m.Subscriber.Use(TracingMiddleware(tracer))
		// There is NO RequestTracingMiddleware implemented in middleware.go?
		// Let's check middleware.go content from Step 27.
		// It has: LoggingMiddleware, PublisherLoggingMiddleware, RequestLoggingMiddleware
		// It has: MetricsMiddleware, PublisherMetricsMiddleware
		// It has: TracingMiddleware, PublisherTracingMiddleware...
		// MISSING: RequestTracingMiddleware!

		// If I try to use PublisherTracingMiddleware for UseRequest, it will fail type check (PublisherFunc vs RequestFunc).
		// So I CANNOT trace requests properly with current middleware.md?
		// PublisherTracingMiddleware wraps PublisherFunc: func(ctx, subject, msgType, data, opts) error
		// RequestFunc is: func(ctx, subject, msgType, data, timeout) (*Envelope, error)

		// I need to implement RequestTracingMiddleware or just skip it for Request in this example.
		// For now, I will skip UseRequest tracing to avoid compilation error, or implementing it is another task.
		// I will just use it for Publish.
	}

	// Start Metrics Server
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		logger.Info("Metrics server starting on :8082")
		if err := http.ListenAndServe(":8082", nil); err != nil {
			logger.Error("Metrics server failed", zap.Error(err))
		}
	}()

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

	// Wait indefinitely to allow metrics scraping
	logger.Info("Publisher finished. Waiting for metrics scrape...")
	select {}
}
