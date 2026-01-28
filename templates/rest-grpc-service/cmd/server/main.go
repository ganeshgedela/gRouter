package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"grouter/pkg/logger"
	"grouter/templates/rest-grpc-service/internal/app"
	"grouter/templates/rest-grpc-service/internal/config"

	"go.uber.org/zap"
)

func main() {
	// Initialize Logger
	log, err := logger.New(logger.Config{
		Level:      "info",
		Format:     "json",
		OutputPath: "stdout",
	})
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	log.Info("starting HTTP-gRPC Service")

	// Load Configuration
	cfg := config.DefaultConfig()

	// Create Application
	application := app.NewApp(log, cfg)

	// Lifecycle Management
	ctx := context.Background()

	if err := application.Init(ctx); err != nil {
		log.Fatal("failed to initialize application", zap.Error(err))
	}

	if err := application.Start(ctx); err != nil {
		log.Fatal("failed to start application", zap.Error(err))
	}

	log.Info("HTTP-gRPC Service running",
		zap.String("http", "http://localhost:8080"),
		zap.String("grpc", "localhost:50051"))

	// Wait for shutdown signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down gracefully")

	if err := application.Stop(ctx); err != nil {
		log.Error("error stopping application", zap.Error(err))
	}

	log.Info("shutdown complete")
}
