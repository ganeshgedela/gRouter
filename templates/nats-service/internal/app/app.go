package app

import (
	"context"
	"fmt"

	"grouter/pkg/manager"
	"grouter/pkg/messaging/nats"
	appconfig "grouter/templates/nats-service/internal/config"

	// Import service packages to register factories
	_ "grouter/templates/nats-service/internal/pkg/orders"
	_ "grouter/templates/nats-service/internal/pkg/users"
)

// App manages the application lifecycle
type App struct {
	deps      manager.Deps
	messenger *nats.Messenger
	manager   *manager.ServiceManager
	appConfig *appconfig.AppConfig
}

// NewApp creates a new application instance
func NewApp(deps manager.Deps) *App {
	return &App{
		deps:      deps,
		appConfig: appconfig.DefaultConfig(),
	}
}

// Init initializes the application
func (a *App) Init(ctx context.Context) error {
	a.deps.Logger.Info("initializing NATS worker application")

	// Initialize NATS Messenger
	messenger := nats.NewMessenger(nil, nil, nil)
	if err := messenger.InitFromConfig(a.deps.Config, a.deps.Logger); err != nil {
		return fmt.Errorf("failed to init messenger: %w", err)
	}

	if err := messenger.Start(); err != nil {
		return fmt.Errorf("failed to start messenger: %w", err)
	}

	a.messenger = messenger

	// Add messenger to deps for service factories
	a.deps.Messenger = messenger

	// Create ServiceManager
	a.manager = manager.NewServiceManager(a.deps)

	// Build services from config using factory pattern
	// The factory will check if "enabled: true" in config
	if err := a.manager.BuildFromConfig(); err != nil {
		return fmt.Errorf("failed to build services: %w", err)
	}

	a.deps.Logger.Info("NATS worker application initialized")
	return nil
}

// Start starts the application
func (a *App) Start(ctx context.Context) error {
	a.deps.Logger.Info("starting NATS worker application")

	// Initialize all services
	if err := a.manager.InitServices(ctx); err != nil {
		return fmt.Errorf("failed to initialize services: %w", err)
	}

	// Start all services
	if err := a.manager.StartServices(ctx); err != nil {
		return fmt.Errorf("failed to start services: %w", err)
	}

	a.deps.Logger.Info("NATS worker application started successfully")
	return nil
}

// Stop stops the application
func (a *App) Stop(ctx context.Context) error {
	a.deps.Logger.Info("stopping NATS worker application")

	// Stop all services via manager
	if a.manager != nil {
		_ = a.manager.StopServices(ctx)
	}

	// Close messenger
	if a.messenger != nil {
		_ = a.messenger.Close()
	}

	a.deps.Logger.Info("NATS worker application stopped")
	return nil
}

// HealthCheck performs a health check
func (a *App) HealthCheck(ctx context.Context) error {
	if a.messenger == nil {
		return fmt.Errorf("messenger not initialized")
	}

	if !a.messenger.IsConnected() {
		return fmt.Errorf("messenger not connected")
	}

	return nil
}

// ID returns the service ID
func (a *App) ID() string { return "nats-worker" }

// Name returns the service name
func (a *App) Name() string { return "NATS Worker" }

// Status returns current health status
func (a *App) Status() manager.HealthStatus { return manager.StatusRunning }

// LastError returns the last error
func (a *App) LastError() error { return nil }

// Dependencies returns list of service IDs this depends on
func (a *App) Dependencies() []string { return nil }
