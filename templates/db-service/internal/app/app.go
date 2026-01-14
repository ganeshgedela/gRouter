package app

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	productv1 "grouter/templates/db-service/api/v1"
	"grouter/templates/db-service/internal/config"
	grpcserver "grouter/templates/db-service/internal/grpc"
	"grouter/templates/db-service/internal/pkg/products"

	"grouter/pkg/database"
	"grouter/pkg/manager"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

// App manages the database service application lifecycle
type App struct {
	deps         manager.Deps
	db           *database.Database
	grpcServer   *grpc.Server
	httpServer   *http.Server
	productRepo  *products.Repository
	appConfig    *config.AppConfig
	grpcListener net.Listener
}

// NewApp creates a new application instance
func NewApp(deps manager.Deps) *App {
	return &App{
		deps:      deps,
		appConfig: config.DefaultConfig(),
	}
}

// Init initializes the application
func (a *App) Init(ctx context.Context) error {
	a.deps.Logger.Info("initializing database service application")

	// Initialize database connection
	if err := a.initDatabase(ctx); err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	// Initialize repository
	a.productRepo = products.NewRepository(a.db.DB)

	// Initialize gRPC server
	if err := a.initGRPCServer(); err != nil {
		return fmt.Errorf("failed to initialize gRPC server: %w", err)
	}

	// Initialize HTTP server for metrics
	a.initMetricsServer()

	a.deps.Logger.Info("database service application initialized")
	return nil
}

// Start starts the application
func (a *App) Start(ctx context.Context) error {
	a.deps.Logger.Info("starting database service application")

	// Start metrics server
	go func() {
		addr := ":8080"
		a.deps.Logger.Info("starting metrics server", zap.String("addr", addr))
		if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.deps.Logger.Error("metrics server failed", zap.Error(err))
		}
	}()

	// Start gRPC server
	go func() {
		a.deps.Logger.Info("starting gRPC server", zap.String("addr", a.grpcListener.Addr().String()))
		if err := a.grpcServer.Serve(a.grpcListener); err != nil {
			a.deps.Logger.Error("gRPC server failed", zap.Error(err))
		}
	}()

	a.deps.Logger.Info("database service started successfully")
	return nil
}

// Stop stops the application
func (a *App) Stop(ctx context.Context) error {
	a.deps.Logger.Info("stopping database service application")

	// Stop gRPC server gracefully
	if a.grpcServer != nil {
		stopped := make(chan struct{})
		go func() {
			a.grpcServer.GracefulStop()
			close(stopped)
		}()

		select {
		case <-stopped:
			a.deps.Logger.Info("gRPC server stopped gracefully")
		case <-time.After(30 * time.Second):
			a.deps.Logger.Warn("gRPC server force stopped after timeout")
			a.grpcServer.Stop()
		}
	}

	// Stop HTTP server
	if a.httpServer != nil {
		if err := a.httpServer.Shutdown(ctx); err != nil {
			a.deps.Logger.Error("error shutting down metrics server", zap.Error(err))
		}
	}

	// Close database connection
	if a.db != nil {
		sqlDB, err := a.db.DB.DB()
		if err == nil {
			sqlDB.Close()
		}
	}

	a.deps.Logger.Info("database service stopped")
	return nil
}

// HealthCheck performs a health check
func (a *App) HealthCheck(ctx context.Context) error {
	if a.db == nil {
		return fmt.Errorf("database not initialized")
	}
	return a.db.HealthCheck(ctx)
}

// ID returns the service ID
func (a *App) ID() string { return "db-service" }

// Name returns the service name
func (a *App) Name() string { return "Database Service" }

// Status returns current health status
func (a *App) Status() manager.HealthStatus { return manager.StatusRunning }

// LastError returns the last error
func (a *App) LastError() error { return nil }

// Dependencies returns list of service IDs this depends on
func (a *App) Dependencies() []string { return nil }

// Helper methods

func (a *App) initDatabase(ctx context.Context) error {
	a.deps.Logger.Info("initializing database connection")

	// Create database connection
	db, err := database.New(a.deps.Config.Database, a.deps.Logger)
	if err != nil {
		return fmt.Errorf("failed to create database connection: %w", err)
	}
	a.db = db

	// Auto-migrate schema
	if err := db.AutoMigrate(&products.Product{}); err != nil {
		return fmt.Errorf("failed to auto-migrate schema: %w", err)
	}

	// Enable metrics collection
	if err := db.EnableMetrics("products_db", 15*time.Second); err != nil {
		a.deps.Logger.Warn("failed to enable database metrics", zap.Error(err))
	}

	// Health check
	if err := db.HealthCheck(ctx); err != nil {
		return fmt.Errorf("database health check failed: %w", err)
	}

	a.deps.Logger.Info("database connection established")
	return nil
}

func (a *App) initGRPCServer() error {
	a.deps.Logger.Info("initializing gRPC server")

	// Create gRPC server with interceptors
	a.grpcServer = grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_prometheus.StreamServerInterceptor,
			grpc_zap.StreamServerInterceptor(a.deps.Logger),
			grpc_recovery.StreamServerInterceptor(),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_prometheus.UnaryServerInterceptor,
			grpc_zap.UnaryServerInterceptor(a.deps.Logger),
			grpc_recovery.UnaryServerInterceptor(),
		)),
	)

	// Register product service
	productServer := grpcserver.NewProductServer(a.productRepo, a.deps.Logger)
	productv1.RegisterProductServiceServer(a.grpcServer, productServer)

	// Register health service
	healthServer := health.NewServer()
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
	healthServer.SetServingStatus("product.v1.ProductService", grpc_health_v1.HealthCheckResponse_SERVING)
	grpc_health_v1.RegisterHealthServer(a.grpcServer, healthServer)

	// Enable reflection for grpcurl/grpcui
	reflection.Register(a.grpcServer)

	// Initialize Prometheus metrics
	grpc_prometheus.Register(a.grpcServer)

	// Create listener
	listener, err := net.Listen("tcp", ":9090")
	if err != nil {
		return fmt.Errorf("failed to create listener: %w", err)
	}
	a.grpcListener = listener

	a.deps.Logger.Info("gRPC server initialized", zap.String("addr", ":9090"))
	return nil
}

func (a *App) initMetricsServer() {
	// Create HTTP server for Prometheus metrics
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	a.httpServer = &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	a.deps.Logger.Info("metrics server initialized", zap.String("addr", ":8080"))
}
