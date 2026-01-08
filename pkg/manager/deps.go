package manager

import (
	"grouter/pkg/config"
	messaging "grouter/pkg/messaging/nats"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// Deps holds the shared infrastructure dependencies that are injected into services.
// This container prevents "constructor explosion" and ensures consistent access to core components.
type Deps struct {
	// Config is the global application configuration.
	Config *config.Config

	// Logger is the configured Zap logger.
	Logger *zap.Logger

	// Messenger handles NATS communication (Pub/Sub).
	Messenger *messaging.Messenger

	// TracerProvider is the OpenTelemetry tracer provider.
	TracerProvider trace.TracerProvider

	// Store provides access to registered services.
	Store *ServiceStore

	// WebRouter is the HTTP router (e.g., *gin.Engine).
	WebRouter interface{}

	// GRPCServer is the gRPC server instance (e.g., *grpc.Server).
	GRPCServer interface{}
}
