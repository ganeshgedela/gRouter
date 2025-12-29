package app

import (
	"context"
	"fmt"

	messaging "grouter/pkg/messaging/nats"
)

// HealthService implements the Service interface for NATS health checks
type HealthService struct {
	app *App
}

// NewHealthService creates a new HealthService
func NewHealthService(app *App) *HealthService {

	//need to register live and ready checks
	app.manager.Health().AddLivenessCheck("app.live", func() error {
		return nil
	})
	app.manager.Health().AddReadinessCheck("nats.ready", func() error {
		if !app.manager.Messenger().IsConnected() {
			return fmt.Errorf("nats not connected")
		}
		return nil
	})
	return &HealthService{
		app: app,
	}
}

// Name returns the service name
func (s *HealthService) Name() string {
	return "health"
}

// Handle processes health check messages (request-reply)
func (s *HealthService) Handle(ctx context.Context, _ string, env *messaging.MessageEnvelope) error {
	if env == nil || env.Reply == "" {
		return nil // fire-and-forget, nothing to respond
	}

	msgType := env.Type

	var (
		checks map[string]string
		err    error
		status string
	)

	switch msgType {
	case "health.live":
		checks, err = s.app.manager.Health().CheckLiveness()
		status = "up"
		if err != nil {
			status = "down"
		}

	case "health.ready":
		checks, err = s.app.manager.Health().CheckReadiness()
		status = "ready"
		if err != nil {
			status = "not ready"
		}

	default:
		err = fmt.Errorf("unknown health type: %s", msgType)
		status = "error"
		checks = map[string]string{}
	}

	resp := map[string]interface{}{
		"status": status,
		"checks": checks,
	}
	if err != nil {
		resp["error"] = err.Error()
	}

	return s.app.manager.Publisher().Publish(ctx, env.Reply, msgType+".response", resp, nil)
}
