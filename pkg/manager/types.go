package manager

import (
	"context"
	messaging "grouter/pkg/messaging/nats"
)

// Service defines the base lifecycle interface for internal services.
type Service interface {
	// Name returns the unique name of the service.
	Name() string
}

// NATService defines a service that handles NATS messages.
type NATService interface {
	Service
	// Handle processes an incoming message and returns a response envelope.
	Handle(ctx context.Context, topic string, msg *messaging.MessageEnvelope) error
}
