package manager

import "context"

// HealthStatus represents the current state of a service.
type HealthStatus string

const (
	StatusCreated     HealthStatus = "created"
	StatusInitialized HealthStatus = "initialized"
	StatusStarting    HealthStatus = "starting"
	StatusRunning     HealthStatus = "running"
	StatusStopping    HealthStatus = "stopping"
	StatusStopped     HealthStatus = "stopped"
	StatusUnhealthy   HealthStatus = "unhealthy"
	StatusFailed      HealthStatus = "failed"
)

// Service defines the lifecycle interface for managed services.
type Service interface {
	// ID returns the unique identifier of the service (e.g., "natdemo").
	ID() string
	// Name returns the human-readable name of the service.
	Name() string

	// Lifecycle methods
	Init(ctx context.Context) error
	Start(ctx context.Context) error
	Stop(ctx context.Context) error

	// Status reporting
	Status() HealthStatus
	LastError() error

	// Dependencies returns the list of service IDs that this service depends on.
	Dependencies() []string
}

// ConfigAware is an interface for services that need configuration injection.
type ConfigAware interface {
	InitConfig(map[string]interface{}) error
}
