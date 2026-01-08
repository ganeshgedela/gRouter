package app

import (
	"context"

	"grouter/pkg/manager"
	messaging "grouter/pkg/messaging/nats"

	"go.uber.org/zap"
)

// BootstrapService waits for a start signal.
type BootstrapService struct {
	messenger *messaging.Messenger
	logger    *zap.Logger
	appName   string

	start chan struct{}
	stop  chan struct{}

	mappings map[string]messaging.HandlerFunc

	status    manager.HealthStatus
	lastError error
}

// init registers the BootstrapService factory.
func init() {
	manager.RegisterFactory("bootstrap", func(deps manager.Deps) (manager.Service, error) {
		return NewBootstrapService(deps), nil
	})
}

// NewBootstrapService creates a new BootstrapService.
func NewBootstrapService(deps manager.Deps) *BootstrapService {
	return &BootstrapService{
		messenger: deps.Messenger,
		logger:    deps.Logger,
		appName:   deps.Config.App.Name,
		start:     make(chan struct{}, 1),
		stop:      make(chan struct{}, 1),
	}
}

func (s *BootstrapService) WaitForStart() <-chan struct{} {
	return s.start
}

func (s *BootstrapService) WaitForStop() <-chan struct{} {
	return s.stop
}

// Name returns the service name "bootstrap".
func (s *BootstrapService) Name() string {
	return "bootstrap"
}

func (s *BootstrapService) ID() string {
	return "bootstrap"
}

func (s *BootstrapService) Status() manager.HealthStatus {
	return s.status
}

func (s *BootstrapService) LastError() error {
	return s.lastError
}

// Dependencies returns the list of dependencies.
func (s *BootstrapService) Dependencies() []string {
	return nil
}

// Init initializes the service.
func (s *BootstrapService) Init(ctx context.Context) error {
	s.mappings = make(map[string]messaging.HandlerFunc)
	s.mappings["start"] = s.HandleStart
	s.mappings["stop"] = s.HandleStop
	return nil
}

func (s *BootstrapService) Start(ctx context.Context) error {
	return s.Subscribe(ctx, nil)
}

func (s *BootstrapService) Stop(ctx context.Context) error {
	return s.Unsubscribe(ctx)
}

// Subscribe registers the start topic.
func (s *BootstrapService) Subscribe(ctx context.Context, options *messaging.SubscribeOptions) error {
	if s.messenger == nil {
		return nil
	}

	for subject, handler := range s.mappings {
		topic := s.appName + "." + subject
		_, err := s.messenger.Subscriber.Subscribe(ctx, topic, handler, options)
		s.logger.Debug("Subscribed to bootstrap topic", zap.String("topic", topic))
		if err != nil {
			return err
		}
	}
	return nil
}

// Unsubscribe unregisters the topic.
func (s *BootstrapService) Unsubscribe(ctx context.Context) error {
	if s.messenger == nil {
		return nil
	}
	for subject := range s.mappings {
		topic := s.appName + "." + subject
		if err := s.messenger.Subscriber.UnsubscribeSubject(ctx, topic); err != nil {
			s.logger.Error("Failed to unsubscribe", zap.Error(err), zap.String("subject", topic))
		}
	}
	return nil
}

// HandleStart processes the start message.
func (s *BootstrapService) HandleStart(ctx context.Context, topic string, env *messaging.MessageEnvelope) error {
	s.logger.Debug("BootstrapService received start message", zap.String("subject", env.Type), zap.String("id", env.ID))
	select {
	case s.start <- struct{}{}:
	default:
	}
	return nil
}

func (s *BootstrapService) HandleStop(ctx context.Context, topic string, env *messaging.MessageEnvelope) error {
	s.logger.Debug("BootstrapService received stop message", zap.String("subject", env.Type), zap.String("id", env.ID))
	select {
	case s.stop <- struct{}{}:
	default:
	}
	return nil
}
