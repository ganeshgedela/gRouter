package app

import (
	"context"
	"fmt"

	"grouter/pkg/manager"
	"grouter/pkg/messaging/nats"
	"grouter/pkg/web"
	appconfig "grouter/templates/rest-nats-service/internal/config"

	"go.uber.org/zap"

	// Import service packages to register factories
	_ "grouter/templates/rest-nats-service/internal/pkg/api"
)

// App manages the REST-NATS application lifecycle (REST + NATS)
type App struct {
	deps      manager.Deps
	webServer web.Server
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
	a.deps.Logger.Info("initializing REST-NATS application")

	// Initialize NATS Messenger
	messenger := nats.NewMessenger(nil, nil, nil)
	if err := messenger.InitFromConfig(a.deps.Config, a.deps.Logger); err != nil {
		return fmt.Errorf("failed to init messenger: %w", err)
	}

	if err := messenger.Start(); err != nil {
		return fmt.Errorf("failed to start messenger: %w", err)
	}
	a.messenger = messenger

	// Add messenger to deps
	a.deps.Messenger = messenger

	// Create Web Server
	webServer, err := web.NewServer("hybrid-api", "REST-NATS API", a.deps.Config.Web, a.deps.Logger)
	if err != nil {
		return fmt.Errorf("failed to create web server: %w", err)
	}
	a.webServer = webServer

	// Add router to deps
	a.deps.WebRouter = a.webServer.Engine()

	// Create ServiceManager
	a.manager = manager.NewServiceManager(a.deps)

	// Build services from config
	if err := a.manager.BuildFromConfig(); err != nil {
		return fmt.Errorf("failed to build services: %w", err)
	}

	// Initialize all services
	if err := a.manager.InitServices(ctx); err != nil {
		return fmt.Errorf("failed to initialize services: %w", err)
	}

	// Note: Route registration happens in Start() after webServer.Init()

	a.deps.Logger.Info("REST-NATS application initialized")
	return nil
}

// registerServiceRoutes registers routes for all WebRoutable services
func (a *App) registerServiceRoutes() {
	router := a.webServer.Engine()
	serviceIDs := a.manager.ListServices()

	for _, serviceID := range serviceIDs {
		svc, err := a.manager.GetService(serviceID)
		if err != nil {
			a.deps.Logger.Warn("failed to get service for route registration",
				zap.String("service_id", serviceID),
				zap.Error(err))
			continue
		}

		// Check if service implements WebRoutable interface
		if routable, ok := svc.(interface {
			IsWebRoutable() bool
			GetRouteRegistrar() func(router interface{})
		}); ok && routable.IsWebRoutable() {
			a.deps.Logger.Info("registering routes for service", zap.String("service", svc.Name()))
			registrar := routable.GetRouteRegistrar()
			registrar(router)
		}
	}
}

// Start starts the application
func (a *App) Start(ctx context.Context) error {
	a.deps.Logger.Info("starting REST-NATS application")

	// Initialize the web server first
	if err := a.webServer.Init(ctx); err != nil {
		return fmt.Errorf("failed to initialize web server: %w", err)
	}

	// NOW register routes (after webServer.Init which sets up the Gin engine)
	a.registerServiceRoutes()

	// Start all services
	if err := a.manager.StartServices(ctx); err != nil {
		return fmt.Errorf("failed to start services: %w", err)
	}

	return a.webServer.Start(ctx)
}

// Stop stops the application
func (a *App) Stop(ctx context.Context) error {
	a.deps.Logger.Info("stopping REST-NATS application")

	// Stop all services
	if a.manager != nil {
		_ = a.manager.StopServices(ctx)
	}

	// Stop web server
	if a.webServer != nil {
		_ = a.webServer.Stop(ctx)
	}

	// Close NATS connection
	if a.messenger != nil {
		_ = a.messenger.Close()
	}

	return nil
}

// HealthCheck performs a health check
func (a *App) HealthCheck(ctx context.Context) error {
	if a.webServer == nil {
		return fmt.Errorf("web server not initialized")
	}

	if a.messenger == nil || !a.messenger.IsConnected() {
		return fmt.Errorf("messenger not connected")
	}

	return nil
}

// ID returns the service ID
func (a *App) ID() string { return "rest-nats-service" }

// Name returns the service name
func (a *App) Name() string { return "REST-NATS Service" }

// Status returns current health status
func (a *App) Status() manager.HealthStatus { return manager.StatusRunning }

// LastError returns the last error
func (a *App) LastError() error { return nil }

// Dependencies returns list of service IDs this depends on
func (a *App) Dependencies() []string { return nil }
