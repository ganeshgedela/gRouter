package users

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
	manager.RegisterFactory("users", func(deps manager.Deps) (manager.Service, error) {
		if deps.Messenger == nil {
			return nil, fmt.Errorf("NATS messenger not provided in dependencies")
		}
		return NewService(deps.Logger, deps.Messenger), nil
	})
}

// UserCreatedEvent represents a user creation event
type UserCreatedEvent struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Timestamp string `json:"timestamp"`
}

// Service handles user-related NATS messages
type Service struct {
	id        string
	logger    *zap.Logger
	messenger *nats.Messenger
	status    manager.HealthStatus
}

// NewService creates a new user service
func NewService(logger *zap.Logger, messenger *nats.Messenger) *Service {
	return &Service{
		id:        "users",
		logger:    logger,
		messenger: messenger,
		status:    manager.StatusCreated,
	}
}

// ID returns the service ID
func (s *Service) ID() string { return s.id }

// Name returns the service name
func (s *Service) Name() string { return "User Service" }

// Status returns current health status
func (s *Service) Status() manager.HealthStatus { return s.status }

// LastError returns the last error
func (s *Service) LastError() error { return nil }

// Dependencies returns list of service IDs this depends on
func (s *Service) Dependencies() []string { return nil }

// Init initializes the service
func (s *Service) Init(ctx context.Context) error {
	s.logger.Info("initializing user service")
	s.status = manager.StatusInitialized
	return nil
}

// Start starts the service
func (s *Service) Start(ctx context.Context) error {
	s.logger.Info("starting user service")
	if err := s.Subscribe(ctx); err != nil {
		s.status = manager.StatusFailed
		return err
	}
	s.status = manager.StatusRunning
	return nil
}

// Stop stops the service
func (s *Service) Stop(ctx context.Context) error {
	s.logger.Info("stopping user service")
	_ = s.Unsubscribe(ctx)
	s.status = manager.StatusStopped
	return nil
}

// Subscribe registers all user-related subscriptions
func (s *Service) Subscribe(ctx context.Context) error {
	opts := &nats.SubscribeOptions{
		QueueGroup:  "user-workers",
		ServiceName: "user-service",
		Durable:     "user-processors",
	}

	_, err := s.messenger.Subscribe(ctx, "users.created", s.HandleUserCreated, opts)
	if err != nil {
		return fmt.Errorf("failed to subscribe to users.created: %w", err)
	}

	s.logger.Info("user service subscriptions registered")
	return nil
}

// HandleUserCreated processes user creation events
func (s *Service) HandleUserCreated(ctx context.Context, subject string, msg *nats.MessageEnvelope) error {
	var event UserCreatedEvent
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		s.logger.Error("failed to unmarshal user event", zap.Error(err))
		return err
	}

	s.logger.Info("processing user creation",
		zap.String("user_id", event.UserID),
		zap.String("email", event.Email),
	)

	// Simulate business logic
	// - Create user profile
	// - Send welcome email
	// - Initialize user preferences

	return nil
}

// Unsubscribe cleans up subscriptions
func (s *Service) Unsubscribe(ctx context.Context) error {
	return s.messenger.Unsubscribe(ctx)
}
