package app

import (
	"context"

	messaging "grouter/pkg/messaging/nats"
)

// BootstrapService waits for a start signal.
type BootstrapService struct {
	trigger chan struct{}
}

// NewBootstrapService creates a new BootstrapService.
func NewBootstrapService(trigger chan struct{}) *BootstrapService {
	return &BootstrapService{
		trigger: trigger,
	}
}

// Name returns the service name "start".
// This matches the topic "natsdemosvc.start" where "start" is the second segment.
func (s *BootstrapService) Name() string {
	return "start"
}

// Handle processes the start message.
func (s *BootstrapService) Handle(ctx context.Context, topic string, env *messaging.MessageEnvelope) error {
	select {
	case s.trigger <- struct{}{}:
	default:
		// If channel is full (already started), do nothing
		// Since we buffered it, this means a signal is already pending.
	}
	return nil
}
