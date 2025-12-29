package app

import (
	"context"

	messaging "grouter/pkg/messaging/nats"
)

// StopService waits for a stop signal to initiate application shutdown.
type StopService struct {
	trigger chan struct{}
}

// NewStopService creates a new StopService.
func NewStopService(trigger chan struct{}) *StopService {
	return &StopService{
		trigger: trigger,
	}
}

// Name returns the service name "stop".
func (s *StopService) Name() string {
	return "stop"
}

// Handle processes the stop message.
func (s *StopService) Handle(ctx context.Context, topic string, env *messaging.MessageEnvelope) error {
	select {
	case s.trigger <- struct{}{}:
	default:
	}
	return nil
}
