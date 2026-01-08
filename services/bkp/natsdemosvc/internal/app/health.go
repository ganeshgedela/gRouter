package app

import (
	"context"
	"fmt"
	health "grouter/pkg/health"
	"grouter/pkg/manager"
	messaging "grouter/pkg/messaging/nats"

	"go.uber.org/zap"
)

// HealthService implements the Service interface for NATS health checks
type HealthService struct {
	health    *health.HealthService
	logger    *zap.Logger
	messenger *messaging.Messenger
	mappings  map[string]messaging.HandlerFunc

	status    manager.HealthStatus
	lastError error
}

// init registers the HealthService factory.
func init() {
	manager.RegisterFactory("health", func(deps manager.Deps) (manager.Service, error) {
		return NewHealthService(deps), nil
	})
}

// NewHealthService creates a new HealthService
func NewHealthService(deps manager.Deps) *HealthService {
	return &HealthService{
		logger:    deps.Logger,
		health:    health.NewHealthService(),
		messenger: deps.Messenger,
	}
}

// Name returns the service name
func (s *HealthService) Name() string {
	return "health"
}

func (s *HealthService) ID() string {
	return "health"
}

func (s *HealthService) Status() manager.HealthStatus {
	return s.status
}

func (s *HealthService) LastError() error {
	return s.lastError
}

// Dependencies returns the list of dependencies.
func (s *HealthService) Dependencies() []string {
	return nil
}

// Init initializes the service.
func (s *HealthService) Init(ctx context.Context) error {
	s.mappings = make(map[string]messaging.HandlerFunc)
	s.mappings["live"] = s.HandleLiveness
	s.mappings["ready"] = s.HandleReadiness

	// register liveness and readiness checks
	s.health.AddLivenessCheck("nats", func() error {
		if s.messenger == nil || !s.messenger.IsConnected() {
			return fmt.Errorf("nats helper not connected")
		}
		return nil
	})
	s.health.AddReadinessCheck("nats", func() error {
		if s.messenger == nil || !s.messenger.IsConnected() {
			return fmt.Errorf("nats helper not connected")
		}
		return nil
	})

	return nil
}

// Start starts the service.
func (s *HealthService) Start(ctx context.Context) error {
	// Typically we might start subscriptions here, but caller (App) seems to call Subscribe explicitly or we can do it here.
	// In BootstrapService, Start/Stop calls Subscribe/Unsubscribe.
	// But in App.go for other services, it calls InitServices then StartServices.
	// However, bootstrap service was treated special.
	// To be safe and compliant with Manager interface, Start should probably subscribe.
	// But wait, natsdemosvc seems to have specific flow in App.go for bootstrap.
	// For regular services:
	// App.LoadServices -> Register -> Init -> Start.
	// So if we register HealthService as a regular service, it should work.
	// But currently HealthService is not registered in App.go's LoadServices logic automatically unless we add it to config?
	// The current code didn't seem to register health service in App.go, let's check App.go after this.
	// For now, I implement standard Start.
	return s.Subscribe(ctx, nil)
}

// Stop stops the service.
func (s *HealthService) Stop(ctx context.Context) error {
	return s.Unsubscribe(ctx)
}

// Subscribe registers the health topics.
func (s *HealthService) Subscribe(ctx context.Context, opts *messaging.SubscribeOptions) error {
	if s.messenger == nil {
		return nil
	}

	for subject, handler := range s.mappings {
		topic := s.messenger.Source() + "." + s.Name() + "." + subject
		_, err := s.messenger.Subscriber.Subscribe(ctx, topic, handler, opts)
		s.logger.Debug("Subscribed to health topic", zap.String("topic", topic))
		if err != nil {
			return err
		}
	}
	return nil
}

// Unsubscribe unregisters the topic.
func (s *HealthService) Unsubscribe(ctx context.Context) error {
	if s.messenger == nil {
		return nil
	}
	for subject := range s.mappings {
		topic := s.messenger.Source() + "." + s.Name() + "." + subject
		if err := s.messenger.Subscriber.UnsubscribeSubject(ctx, topic); err != nil {
			s.logger.Error("Failed to unsubscribe", zap.Error(err), zap.String("subject", topic))
		}
	}
	return nil
}

// HandleLive processes health.live messages
func (s *HealthService) HandleLiveness(ctx context.Context, subject string, env *messaging.MessageEnvelope) error {
	if env == nil || env.Reply == "" {
		return nil // fire-and-forget
	}

	resp, err := s.health.CheckLiveness()
	s.logger.Debug("Liveness check", zap.Any("response", resp), zap.Error(err))
	if err != nil {
		return s.app.messenger.Publisher.Publish(ctx, env.Reply, subject+".error", resp, nil)
	}
	return s.app.messenger.Publisher.Publish(ctx, env.Reply, subject+".response", resp, nil)
}

// HandleReadiness processes health.ready messages
func (s *HealthService) HandleReadiness(ctx context.Context, subject string, env *messaging.MessageEnvelope) error {
	if env == nil || env.Reply == "" {
		return nil // fire-and-forget
	}

	resp, err := s.health.CheckReadiness()
	s.logger.Debug("Readiness check", zap.Any("response", resp), zap.Error(err))
	if err != nil {
		return s.app.messenger.Publisher.Publish(ctx, env.Reply, subject+".error", resp, nil)
	}
	return s.app.messenger.Publisher.Publish(ctx, env.Reply, subject+".response", resp, nil)
}
