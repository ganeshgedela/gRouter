package orders

import (
	"context"
	"encoding/json"
	"fmt"

	"grouter/pkg/manager"
	"grouter/pkg/messaging/nats"

	"go.uber.org/zap"
)

// Register factory on package import
func init() {
	manager.RegisterFactory("orders", func(deps manager.Deps) (manager.Service, error) {
		if deps.Messenger == nil {
			return nil, fmt.Errorf("NATS messenger not provided in dependencies")
		}
		return NewService(deps.Logger, deps.Messenger), nil
	})
}

// OrderCreatedEvent represents an order creation event
type OrderCreatedEvent struct {
	OrderID   string  `json:"order_id"`
	UserID    string  `json:"user_id"`
	Amount    float64 `json:"amount"`
	Timestamp string  `json:"timestamp"`
}

// Service handles order-related NATS messages
type Service struct {
	id        string
	logger    *zap.Logger
	messenger *nats.Messenger
	status    manager.HealthStatus
}

// NewService creates a new order service
func NewService(logger *zap.Logger, messenger *nats.Messenger) *Service {
	return &Service{
		id:        "orders",
		logger:    logger,
		messenger: messenger,
		status:    manager.StatusCreated,
	}
}

// ID returns the service ID
func (s *Service) ID() string { return s.id }

// Name returns the service name
func (s *Service) Name() string { return "Order Service" }

// Status returns current health status
func (s *Service) Status() manager.HealthStatus { return s.status }

// LastError returns the last error
func (s *Service) LastError() error { return nil }

// Dependencies returns list of service IDs this depends on
func (s *Service) Dependencies() []string { return nil }

// Init initializes the service
func (s *Service) Init(ctx context.Context) error {
	s.logger.Info("initializing order service")
	s.status = manager.StatusInitialized
	return nil
}

// Start starts the service
func (s *Service) Start(ctx context.Context) error {
	s.logger.Info("starting order service")
	if err := s.Subscribe(ctx); err != nil {
		s.status = manager.StatusFailed
		return err
	}
	s.status = manager.StatusRunning
	return nil
}

// Stop stops the service
func (s *Service) Stop(ctx context.Context) error {
	s.logger.Info("stopping order service")
	_ = s.Unsubscribe(ctx)
	s.status = manager.StatusStopped
	return nil
}

// Subscribe registers all order-related subscriptions
func (s *Service) Subscribe(ctx context.Context) error {
	opts := &nats.SubscribeOptions{
		QueueGroup:  "order-workers",
		ServiceName: "order-service",
		Durable:     "order-processors",
	}

	_, err := s.messenger.Subscribe(ctx, "orders.created", s.HandleOrderCreated, opts)
	if err != nil {
		return fmt.Errorf("failed to subscribe to orders.created: %w", err)
	}

	s.logger.Info("order service subscriptions registered")
	return nil
}

// HandleOrderCreated processes order creation events
func (s *Service) HandleOrderCreated(ctx context.Context, subject string, msg *nats.MessageEnvelope) error {
	var event OrderCreatedEvent
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		s.logger.Error("failed to unmarshal order event", zap.Error(err))
		return err
	}

	s.logger.Info("processing order",
		zap.String("order_id", event.OrderID),
		zap.String("user_id", event.UserID),
		zap.Float64("amount", event.Amount),
	)

	// Simulate business logic
	// - Validate order
	// - Update inventory
	// - Process payment
	// - Publish order.processed event

	return nil
}

// Unsubscribe cleans up subscriptions
func (s *Service) Unsubscribe(ctx context.Context) error {
	return s.messenger.Unsubscribe(ctx)
}
