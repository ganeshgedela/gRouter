package manager

import (
	"context"
)

// WebRoutable allows a service to register HTTP routes.
// Services implementing this interface will have RegisterRoutes called during startup.
type WebRoutable interface {
	// RegisterRoutes registers the service's HTTP endpoints.
	// The router argument is generic (interface{}) to decouple from the specific web framework (Gin),
	// but expected to be *gin.Engine in this implementation.
	RegisterRoutes(router interface{}) error
}

// GRPCRoutable allows a service to register gRPC servers.
type GRPCRoutable interface {
	// RegisterGRPC registers the gRPC server implementation.
	// server should be *grpc.Server.
	RegisterGRPC(server interface{}) error
}

// NATSRoutable allows a service to declare NATS interest.
// Services implementing this interface can encapsulate their subscription logic.
type NATSRoutable interface {
	// StartNATS is called to initialize NATS subscriptions.
	// This is often synonymous with Start(), but explicit interface allows for separation if needed.
	// For now, we rely on standard Start(), but this interface marks the capability.
	Subjects() []string
	Handle(ctx context.Context, subject string, payload []byte) ([]byte, error)
}
