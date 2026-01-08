package api

import (
	"context"
	"fmt"

	"grouter/pkg/manager"
	"grouter/pkg/messaging/nats"

	"go.uber.org/zap"
)

// Register factory on package import
func init() {
	manager.RegisterFactory("message_handler", func(deps manager.Deps) (manager.Service, error) {
		if deps.Messenger == nil {
			return nil, fmt.Errorf("messenger not provided in dependencies")
		}
		return NewMessageHandler(deps.Logger, deps.Messenger), nil
	})
}

// MessageHandler handles NATS messages and can trigger gRPC calls
type MessageHandler struct {
	id        string
	logger    *zap.Logger
	messenger *nats.Messenger
	status    manager.HealthStatus
}

// NewMessageHandler creates a new message handler
func NewMessageHandler(logger *zap.Logger, messenger *nats.Messenger) *MessageHandler {
	return &MessageHandler{
		id:        "message_handler",
		logger:    logger,
		messenger: messenger,
		status:    manager.StatusCreated,
	}
}

// ID returns the service ID
func (h *MessageHandler) ID() string { return h.id }

// Name returns the service name
func (h *MessageHandler) Name() string { return "Message Handler" }

// Status returns current health status
func (h *MessageHandler) Status() manager.HealthStatus { return h.status }

// LastError returns the last error
func (h *MessageHandler) LastError() error { return nil }

// Dependencies returns list of service IDs this depends on
func (h *MessageHandler) Dependencies() []string { return nil }

// Init initializes the service
func (h *MessageHandler) Init(ctx context.Context) error {
	h.logger.Info("initializing message handler")
	h.status = manager.StatusInitialized
	return nil
}

// Start starts the service
func (h *MessageHandler) Start(ctx context.Context) error {
	h.logger.Info("starting message handler")

	// Subscribe to event topics
	if _, err := h.messenger.Subscribe(ctx, "events.>", h.handleEvent, &nats.SubscribeOptions{}); err != nil {
		h.status = manager.StatusFailed
		return fmt.Errorf("failed to subscribe to events: %w", err)
	}

	// Subscribe to command topics
	if _, err := h.messenger.Subscribe(ctx, "commands.>", h.handleCommand, &nats.SubscribeOptions{}); err != nil {
		h.status = manager.StatusFailed
		return fmt.Errorf("failed to subscribe to commands: %w", err)
	}

	h.logger.Info("message handler subscriptions registered")
	h.status = manager.StatusRunning
	return nil
}

// Stop stops the service
func (h *MessageHandler) Stop(ctx context.Context) error {
	h.logger.Info("stopping message handler")
	h.status = manager.StatusStopped
	return nil
}

// handleEvent processes event messages
func (h *MessageHandler) handleEvent(ctx context.Context, subject string, envelope *nats.MessageEnvelope) error {
	h.logger.Info("received event",
		zap.String("subject", subject),
		zap.String("type", envelope.Type),
		zap.String("id", envelope.ID))

	// Process event - could trigger gRPC calls here
	// Example: Make gRPC call to another service based on event

	return nil
}

// handleCommand processes command messages
func (h *MessageHandler) handleCommand(ctx context.Context, subject string, envelope *nats.MessageEnvelope) error {
	h.logger.Info("received command",
		zap.String("subject", subject),
		zap.String("type", envelope.Type),
		zap.String("id", envelope.ID))

	// Process command - could trigger gRPC calls here
	// Example: Execute business logic and call gRPC services

	return nil
}
