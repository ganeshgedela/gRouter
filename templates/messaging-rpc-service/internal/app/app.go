package app

import (
	"context"
	"fmt"

	"grouter/pkg/manager"
	"grouter/pkg/messaging/grpc"
	"grouter/pkg/messaging/nats"
	pb "grouter/templates/grpc-service/api/proto"
	appconfig "grouter/templates/messaging-rpc-service/internal/config"
	internal "grouter/templates/messaging-rpc-service/internal/grpc"

	googlegrpc "google.golang.org/grpc"

	// Import service packages to register factories
	_ "grouter/templates/messaging-rpc-service/internal/pkg/api"
)

// App manages the Messaging+RPC application lifecycle
type App struct {
	deps       manager.Deps
	grpcServer *grpc.Server
	messenger  *nats.Messenger
	manager    *manager.ServiceManager
	appConfig  *appconfig.AppConfig
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
	a.deps.Logger.Info("initializing Messaging+RPC application")

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

	// Create gRPC Server
	serverOpts := []grpc.Option{
		grpc.WithPort(a.deps.Config.GRPC.Port),
	}

	srv := grpc.NewServer(a.deps.Logger, serverOpts...)

	// Register gRPC service implementation
	srv.RegisterService(func(s googlegrpc.ServiceRegistrar) {
		helloServer := internal.NewHelloServer(a.deps.Logger)
		pb.RegisterHelloServiceServer(s, helloServer)
	})

	a.grpcServer = srv

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

	a.deps.Logger.Info("Messaging+RPC application initialized")
	return nil
}

// Start starts the application
func (a *App) Start(ctx context.Context) error {
	a.deps.Logger.Info("starting Messaging+RPC application")

	// Start all services (NATS handlers)
	if err := a.manager.StartServices(ctx); err != nil {
		return fmt.Errorf("failed to start services: %w", err)
	}

	// Start gRPC server
	if err := a.grpcServer.Start(); err != nil {
		return fmt.Errorf("failed to start gRPC server: %w", err)
	}

	a.deps.Logger.Info("Messaging+RPC application started successfully")
	return nil
}

// Stop stops the application
func (a *App) Stop(ctx context.Context) error {
	a.deps.Logger.Info("stopping Messaging+RPC application")

	// Stop all services
	if a.manager != nil {
		_ = a.manager.StopServices(ctx)
	}

	// Stop gRPC server
	if a.grpcServer != nil {
		a.grpcServer.Stop(ctx)
	}

	// Close NATS connection
	if a.messenger != nil {
		_ = a.messenger.Close()
	}

	return nil
}

// HealthCheck performs a health check
func (a *App) HealthCheck(ctx context.Context) error {
	if a.grpcServer == nil {
		return fmt.Errorf("gRPC server not initialized")
	}

	if a.messenger == nil || !a.messenger.IsConnected() {
		return fmt.Errorf("messenger not connected")
	}

	return nil
}

// ID returns the service ID
func (a *App) ID() string { return "messaging-rpc-service" }

// Name returns the service name
func (a *App) Name() string { return "Messaging+RPC Service" }

// Status returns current health status
func (a *App) Status() manager.HealthStatus { return manager.StatusRunning }

// LastError returns the last error
func (a *App) LastError() error { return nil }

// Dependencies returns list of service IDs this depends on
func (a *App) Dependencies() []string { return nil }
