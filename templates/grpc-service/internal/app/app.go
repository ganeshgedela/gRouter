package app

import (
	"context"
	"fmt"

	"grouter/pkg/manager"
	"grouter/pkg/messaging/grpc"
	pb "grouter/templates/grpc-service/api/proto"
	appconfig "grouter/templates/grpc-service/internal/config"
	internal "grouter/templates/grpc-service/internal/grpc"

	googlegrpc "google.golang.org/grpc"
)

// App manages the gRPC application lifecycle
type App struct {
	deps      manager.Deps
	server    *grpc.Server
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
	a.deps.Logger.Info("initializing gRPC application")

	// Create gRPC Server with options
	serverOpts := []grpc.Option{
		grpc.WithPort(a.deps.Config.GRPC.Port),
	}

	// Reflection disabled for now to avoid Bazel rules_proto dependency
	// if a.deps.Config.GRPC.ReflectionEnabled {
	// 	serverOpts = append(serverOpts, grpc.WithReflection())
	// }

	srv := grpc.NewServer(a.deps.Logger, serverOpts...)

	// Register service implementation
	srv.RegisterService(func(s googlegrpc.ServiceRegistrar) {
		helloServer := internal.NewHelloServer(a.deps.Logger)
		pb.RegisterHelloServiceServer(s, helloServer)
	})

	a.server = srv
	a.deps.Logger.Info("gRPC application initialized")
	return nil
}

// Start starts the application
func (a *App) Start(ctx context.Context) error {
	a.deps.Logger.Info("starting gRPC application")
	return a.server.Start()
}

// Stop stops the application
func (a *App) Stop(ctx context.Context) error {
	a.deps.Logger.Info("stopping gRPC application")
	if a.server != nil {
		a.server.Stop(ctx)
	}
	return nil
}

// HealthCheck performs a health check
func (a *App) HealthCheck(ctx context.Context) error {
	if a.server == nil {
		return fmt.Errorf("server not initialized")
	}
	// gRPC server health check logic
	return nil
}

// ID returns the service ID
func (a *App) ID() string { return "grpc-service" }

// Name returns the service name
func (a *App) Name() string { return "gRPC Service" }

// Status returns current health status
func (a *App) Status() manager.HealthStatus { return manager.StatusRunning }

// LastError returns the last error
func (a *App) LastError() error { return nil }

// Dependencies returns list of service IDs this depends on
func (a *App) Dependencies() []string { return nil }
