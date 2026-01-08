package app

import (
	"context"
	"fmt"

	"grouter/pkg/manager"
	"grouter/pkg/web"
	appconfig "grouter/templates/rest-service/internal/config"

	// Import service packages to register factories
	_ "grouter/templates/rest-service/internal/pkg/orders"
	_ "grouter/templates/rest-service/internal/pkg/users"

	"github.com/gin-gonic/gin"
)

// App manages the REST API application lifecycle
type App struct {
	deps      manager.Deps
	server    web.Server
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
	a.deps.Logger.Info("initializing REST API application")

	// Create Web Server
	srv, err := web.NewServer("rest-api", "REST API Service", a.deps.Config.Web, a.deps.Logger)
	if err != nil {
		return fmt.Errorf("failed to create web server: %w", err)
	}
	a.server = srv

	// Add router to deps for services
	a.deps.WebRouter = a.server.Engine()

	// Create ServiceManager
	a.manager = manager.NewServiceManager(a.deps)

	// Build services from config using factory pattern
	if err := a.manager.BuildFromConfig(); err != nil {
		return fmt.Errorf("failed to build services: %w", err)
	}

	// Initialize all services
	if err := a.manager.InitServices(ctx); err != nil {
		return fmt.Errorf("failed to initialize services: %w", err)
	}

	// Register routes with services implementing WebRoutable
	a.registerRoutes(a.server.Engine())

	a.deps.Logger.Info("REST API application initialized")
	return nil
}

// Start starts the application
func (a *App) Start(ctx context.Context) error {
	a.deps.Logger.Info("starting REST API application")

	// Initialize the web server
	if err := a.server.Init(ctx); err != nil {
		return fmt.Errorf("failed to initialize web server: %w", err)
	}

	// Start all services
	if err := a.manager.StartServices(ctx); err != nil {
		return fmt.Errorf("failed to start services: %w", err)
	}

	// Start the web server
	return a.server.Start(ctx)
}

// Stop stops the application
func (a *App) Stop(ctx context.Context) error {
	a.deps.Logger.Info("stopping REST API application")

	// Stop all services
	if a.manager != nil {
		_ = a.manager.StopServices(ctx)
	}

	if a.server != nil {
		return a.server.Stop(ctx)
	}
	return nil
}

// HealthCheck performs a health check
func (a *App) HealthCheck(ctx context.Context) error {
	if a.server == nil {
		return fmt.Errorf("server not initialized")
	}
	return nil
}

// registerRoutes registers all application routes
func (a *App) registerRoutes(router *gin.Engine) {
	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Health endpoint
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "healthy"})
		})

		// Services implementing WebRoutable will register their own routes
		// This is now handled automatically by the service manager during Init
	}
}

// ID returns the service ID
func (a *App) ID() string { return "rest-api" }

// Name returns the service name
func (a *App) Name() string { return "REST API" }

// Status returns current health status
func (a *App) Status() manager.HealthStatus { return manager.StatusRunning }

// LastError returns the last error
func (a *App) LastError() error { return nil }

// Dependencies returns list of service IDs this depends on
func (a *App) Dependencies() []string { return nil }
