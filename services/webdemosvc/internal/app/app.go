package app

import (
	"context"
	"fmt"

	"github.com/go-viper/mapstructure/v2"
	"github.com/google/uuid"

	"grouter/pkg/manager"
	"grouter/services/webdemosvc/internal/pkg/webdemo"

	"go.uber.org/zap"
)

type App struct {
	manager *manager.ServiceManager
	AppId   string

	startChan chan struct{}
	stopChan  chan struct{}
}

func New() *App {
	return &App{
		manager:   manager.NewServiceManager(),
		startChan: make(chan struct{}),
		stopChan:  make(chan struct{}),
	}
}

func (a *App) Init() error {
	if err := a.manager.Init(); err != nil {
		return fmt.Errorf("failed to init manager: %w", err)
	}
	if err := a.manager.InitNATS(); err != nil {
		return fmt.Errorf("failed to init nats: %w", err)
	}
	if err := a.manager.InitWebServer(); err != nil {
		return fmt.Errorf("failed to init web server: %w", err)
	}
	// Generate unique AppId
	a.AppId = a.manager.Config().App.Name + "-" + uuid.New().String()
	a.manager.Logger().Info("App initialized", zap.String("AppId", a.AppId))

	// Register Health Service
	healthSvc := NewHealthService(a.manager.Health(), a.manager.Config().App.Name)
	if err := a.manager.RegisterService(healthSvc); err != nil {
		return fmt.Errorf("failed to register health service: %w", err)
	}

	return nil
}

func (a *App) GetAppName() string {
	return a.manager.Config().App.Name
}

func (a *App) RegisterBootstrap() error {
	logger := a.manager.Logger()
	bootstrap := NewBootstrapService(a.startChan)
	if err := a.manager.RegisterService(bootstrap); err != nil {
		logger.Error("Failed to register bootstrap service", zap.Error(err))
	}

	logger.Info("Registering Bootstrap Service (HTTP only)")
	return nil
}

func (a *App) GetManager() *manager.ServiceManager {
	return a.manager
}

func (a *App) RegisterStop() error {
	logger := a.manager.Logger()
	// Register Stop Service (HTTP only)
	stopSvc := NewStopService(a.stopChan, a.manager.WebServer())
	if err := a.manager.RegisterService(stopSvc); err != nil {
		logger.Error("Failed to register stop service", zap.Error(err))
	}

	logger.Info("Registering Stop Service (HTTP only)")
	return nil
}

func (a *App) Start(ctx context.Context) error {
	logger := a.manager.Logger()

	logger.Info("Starting " + a.GetAppName() + "...")

	if err := a.RegisterBootstrap(); err != nil {
		logger.Error("Failed to register bootstrap service", zap.Error(err))
		return err
	}
	if err := a.RegisterStop(); err != nil {
		logger.Error("Failed to register stop service", zap.Error(err))
		return err
	}

	logger.Info("Waiting for start signal...")
	logger.Info("Services: ", zap.Any("services", a.manager.ListServices()))

	// Start Manager (begins message listening and web server)
	if err := a.manager.Start(ctx); err != nil {
		return err
	}

	for {
		// Block until start message is received
		select {
		case <-a.startChan:
			logger.Info("Start signal received. Registering services...")
			// Register services via config
			if err := a.RegisterServices(); err != nil {
				logger.Error("Failed to register services", zap.Error(err))
			}

			// Reload web server to apply new services
			logger.Info("Reloading web server to apply new services...")
			// Note: This might cause a race with the /start response if not handled carefully,
			// but for now it's necessary to register routes on the running engine.
			// ResetEngine stops the server first.
			if err := a.manager.WebServer().ResetEngine(context.Background()); err != nil {
				logger.Error("Failed to reset engine", zap.Error(err))
			}
			a.manager.ReRegisterServices()
			if err := a.manager.WebServer().Start(); err != nil {
				logger.Error("Failed to start web server", zap.Error(err))
			}
		case <-ctx.Done():
			return ctx.Err()
		}

		// Block until stop message is received
		select {
		case <-a.stopChan:
			logger.Info("Stop signal received. Unregistering services...")
			// Unregister services via config
			if err := a.UnregisterServices(); err != nil {
				logger.Error("Failed to unregister services", zap.Error(err))
			}
			// Reload web server to apply new services
			logger.Info("Reloading web server to apply new services...")
			// Note: This might cause a race with the /start response if not handled carefully,
			// but for now it's necessary to register routes on the running engine.
			// ResetEngine stops the server first.
			if err := a.manager.WebServer().ResetEngine(context.Background()); err != nil {
				logger.Error("Failed to reset engine", zap.Error(err))
			}
			a.manager.ReRegisterServices()
			if err := a.manager.WebServer().Start(); err != nil {
				logger.Error("Failed to start web server", zap.Error(err))
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (a *App) Stop(ctx context.Context) error {
	return a.manager.Stop(ctx)
}

func (a *App) UnregisterServices() error {
	logger := a.manager.Logger()
	services := a.manager.ListServices()

	for _, service := range services {
		if service == "start" || service == "stop" || service == "health" {
			continue
		}
		logger.Info("Unregistering service: " + service)
		a.manager.UnregisterService(service)
	}

	logger.Info("Services: ", zap.Any("services", a.manager.ListServices()))
	return nil
}

func (a *App) ShutdownChan() <-chan struct{} {
	// Return a never-closed channel so main.go blocks until OS signal
	// or until we implement internal shutdown trigger (which we haven't yet, distinct from NATS stop).
	return make(chan struct{})
}

func (a *App) RegisterServices() error {
	logger := a.manager.Logger()
	cfg := a.manager.Config()

	// build services list from cfg.Services
	// build services list from cfg.Services

	for name, serviceCfg := range cfg.Services {
		if name == "webdemosvc" {
			var webConfig webdemo.WebDemoConfig
			if err := decodeConfig(serviceCfg, &webConfig); err != nil {
				logger.Error("Failed to decode WebDemo config", zap.Error(err))
				return err
			}
			// Register WebDemo Service
			webSvc := webdemo.NewService()
			if err := a.manager.RegisterService(webSvc); err != nil {
				return err
			}
			logger.Info("Registered WebDemo Service")
		}
	}

	// Currently no dynamic services to register for webdemo,
	// but this hook is available for future expansion.
	logger.Info("Services: ", zap.Any("services", a.manager.ListServices()))
	return nil
}

func (a *App) Logger() *zap.Logger {
	return a.manager.Logger()
}

func decodeConfig(input interface{}, output interface{}) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   output,
		TagName:  "mapstructure",
	})
	if err != nil {
		return err
	}
	return decoder.Decode(input)
}
