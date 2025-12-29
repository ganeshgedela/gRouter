package manager

import (
	"context"
	"fmt"
	"strings"

	messaging "grouter/pkg/messaging/nats"
)

// ServiceRouter routes messages to the appropriate service based on the topic.
type ServiceRouter struct {
	store *ServiceStore
}

// NewServiceRouter creates a new ServiceRouter.
func NewServiceRouter() *ServiceRouter {
	return &ServiceRouter{store: NewServiceStore()}
}

// Register adds a service to the router.
func (r *ServiceRouter) Register(name string, svc Service) {
	r.store.Add(name, svc)
}

// Unregister removes a service from the router.
func (r *ServiceRouter) Unregister(name string) {
	r.store.Delete(name)
}

// List returns a list of all registered service names.
func (r *ServiceRouter) List() []string {
	return r.store.List()
}

// RouteByTopic finds the service registered for the given topic.
func (r *ServiceRouter) RouteByTopic(topic string) (Service, error) {
	topic = strings.TrimSpace(topic)
	if topic == "" {
		return nil, fmt.Errorf("empty topic")
	}

	parts := strings.Split(topic, ".")
	if len(parts) < 2 {
		// If topic is just "natdemo", try to look it up directly or fail gracefully
		if len(parts) == 1 && parts[0] != "" {
			serviceName := parts[0]
			svc, ok := r.store.Get(serviceName)
			if ok {
				return svc, nil
			}
		}
		return nil, fmt.Errorf("invalid topic format: %q (expected service.action)", topic)
	}

	serviceName := parts[0]
	svc, ok := r.store.Get(serviceName)
	if !ok {
		return nil, fmt.Errorf("no service registered for topic: %q", serviceName)
	}
	return svc, nil
}

// HandleMessage routes the message to the correct service and calls its Handle method.
func (r *ServiceRouter) HandleMessage(ctx context.Context, topic string, env *messaging.MessageEnvelope) error {
	if env == nil {
		return fmt.Errorf("nil envelope")
	}
	svc, err := r.RouteByTopic(topic)
	if err != nil {
		return err
	}

	natSvc, ok := svc.(NATService)
	if !ok {
		return fmt.Errorf("service %q does not support NATS handling", svc.Name())
	}

	return natSvc.Handle(ctx, topic, env)
}
