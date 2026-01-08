package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"grouter/pkg/config"
	"grouter/pkg/logger"
	"grouter/pkg/manager"
	"grouter/templates/nats-service/internal/app"

	"go.uber.org/zap"
)

func main() {
	// 1. Initialize Logger
	log, err := logger.New(logger.Config{
		Level:      "info",
		Format:     "json",
		OutputPath: "stdout",
	})
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	log.Info("starting NATS worker service")

	// 2. Load Configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("failed to load configuration", zap.Error(err))
	}

	// 3. Setup Dependencies
	deps := manager.Deps{
		Config: cfg,
		Logger: log,
	}

	// 4. Create Service Manager
	mgr := manager.NewServiceManager(deps)

	// 5. Create and Register Application
	application := app.NewApp(deps)
	if err := mgr.RegisterService(application); err != nil {
		log.Fatal("failed to register application", zap.Error(err))
	}

	// 6. Lifecycle Management
	ctx := context.Background()

	if err := mgr.InitServices(ctx); err != nil {
		log.Fatal("failed to initialize services", zap.Error(err))
	}

	if err := mgr.StartServices(ctx); err != nil {
		log.Fatal("failed to start services", zap.Error(err))
	}

	log.Info("NATS worker service running")

	// 7. Wait for shutdown signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down gracefully")

	if err := mgr.StopServices(ctx); err != nil {
		log.Error("error stopping services", zap.Error(err))
	}
}
