package app

import (
	"context"
	"fmt"
	"net"
	"net/http"

	pb "grouter/templates/rest-grpc-service/api/proto"
	"grouter/templates/rest-grpc-service/internal/pkg/business"
	"grouter/templates/rest-grpc-service/internal/config"
	grpcserver "grouter/templates/rest-grpc-service/internal/grpc"
	"grouter/templates/rest-grpc-service/internal/pkg/rest"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// App manages both HTTP and gRPC servers
type App struct {
	logger       *zap.Logger
	config       *config.Config
	itemService  *business.ItemService
	httpServer   *http.Server
	grpcServer   *grpc.Server
	httpListener net.Listener
	grpcListener net.Listener
}

// NewApp creates a new application
func NewApp(logger *zap.Logger, cfg *config.Config) *App {
	return &App{
		logger:      logger,
		config:      cfg,
		itemService: business.NewItemService(logger),
	}
}

// Init initializes both servers
func (a *App) Init(ctx context.Context) error {
	a.logger.Info("initializing HTTP-gRPC service")

	// Initialize HTTP server
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	// Register REST routes
	restHandler := rest.NewHandler(a.logger, a.itemService)
	restHandler.RegisterRoutes(router)

	a.httpServer = &http.Server{
		Addr:    ":" + a.config.HTTPPort,
		Handler: router,
	}

	// Initialize gRPC server
	a.grpcServer = grpc.NewServer()
	grpcHandler := grpcserver.NewServer(a.logger, a.itemService)
	pb.RegisterItemServiceServer(a.grpcServer, grpcHandler)

	// Enable reflection for grpcurl
	reflection.Register(a.grpcServer)

	a.logger.Info("HTTP-gRPC service initialized",
		zap.String("http_port", a.config.HTTPPort),
		zap.String("grpc_port", a.config.GRPCPort))

	return nil
}

// Start starts both HTTP and gRPC servers
func (a *App) Start(ctx context.Context) error {
	a.logger.Info("starting HTTP-gRPC service")

	// Start HTTP server
	httpListener, err := net.Listen("tcp", ":"+a.config.HTTPPort)
	if err != nil {
		return fmt.Errorf("failed to create HTTP listener: %w", err)
	}
	a.httpListener = httpListener

	go func() {
		a.logger.Info("HTTP server listening", zap.String("port", a.config.HTTPPort))
		if err := a.httpServer.Serve(httpListener); err != nil && err != http.ErrServerClosed {
			a.logger.Error("HTTP server error", zap.Error(err))
		}
	}()

	// Start gRPC server
	grpcListener, err := net.Listen("tcp", ":"+a.config.GRPCPort)
	if err != nil {
		return fmt.Errorf("failed to create gRPC listener: %w", err)
	}
	a.grpcListener = grpcListener

	go func() {
		a.logger.Info("gRPC server listening", zap.String("port", a.config.GRPCPort))
		if err := a.grpcServer.Serve(grpcListener); err != nil {
			a.logger.Error("gRPC server error", zap.Error(err))
		}
	}()

	a.logger.Info("HTTP-gRPC service started",
		zap.String("http_endpoint", "http://localhost:"+a.config.HTTPPort),
		zap.String("grpc_endpoint", "localhost:"+a.config.GRPCPort))

	return nil
}

// Stop gracefully stops both servers
func (a *App) Stop(ctx context.Context) error {
	a.logger.Info("stopping HTTP-gRPC service")

	// Stop HTTP server
	if a.httpServer != nil {
		if err := a.httpServer.Shutdown(ctx); err != nil {
			a.logger.Error("HTTP server shutdown error", zap.Error(err))
		}
	}

	// Stop gRPC server
	if a.grpcServer != nil {
		a.grpcServer.GracefulStop()
	}

	a.logger.Info("HTTP-gRPC service stopped")
	return nil
}
