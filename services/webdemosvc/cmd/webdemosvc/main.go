package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"grouter/services/webdemosvc/internal/app"

	"go.uber.org/zap"
)

func main() {
	// Create application instance
	application := app.New()

	// Initialize application
	if err := application.Init(); err != nil {
		l, _ := zap.NewProduction()
		l.Fatal("Failed to initialize application", zap.Error(err))
	}

	// Create context that listens for signals
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start App in a goroutine so it doesn't block signal handling
	go func() {
		if err := application.Start(ctx); err != nil {
			// If Start returns an error (other than context canceled), log it
			if err != context.Canceled {
				application.Logger().Fatal("Failed to start app", zap.Error(err))
			}
		}
	}()

	// Handle OS signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Wait for shutdown signal or app stop request
	select {
	case <-sigChan:
		application.Logger().Info("Received OS signal")
	case <-application.ShutdownChan():
		application.Logger().Info("Received API stop signal")
	}

	// Graceful shutdown logic
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()

	if err := application.Stop(shutdownCtx); err != nil {
		application.Logger().Error("Error during shutdown", zap.Error(err))
	}
}
