package app

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"grouter/pkg/config"
	"grouter/pkg/logger"
	"grouter/pkg/manager"
)

type App struct {
	manager *manager.ServiceManager
	AppId   string
}

func New() *App {
	return &App{}
}

func (a *App) Init() error {
	// 1. Load Configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// 2. Initialize Logger
	log, err := logger.New(logger.Config{
		Level:      cfg.Log.Level,
		Format:     cfg.Log.Format,
		OutputPath: cfg.Log.OutputPath,
	})
	if err != nil {
		return fmt.Errorf("failed to init logger: %w", err)
	}

	// 3. Initialize Deps
	deps := manager.Deps{
		Config: cfg,
		Logger: log,
	}

	// 4. Create Service Manager
	a.manager = manager.NewServiceManager(deps)

	// 5. Build Services from Config
	if err := a.manager.BuildFromConfig(); err != nil {
		return err
	}

	// Generate unique AppId
	a.AppId = cfg.App.Name + "-" + strings.Split(uuid.New().String(), "-")[0]
	log.Debug("App initialized", zap.String("AppId", a.AppId))

	return nil
}

func (a *App) GetAppName() string {
	return a.manager.Config().App.Name
}

// Start starts the application
func (a *App) Start(ctx context.Context) error {
	// Directly load and start services since we don't have NATS bootstrap
	if err := a.LoadServices(ctx); err != nil {
		return err
	}

	// Block until context is done (simulate running forever until signal)
	<-ctx.Done()
	return nil
}

func (a *App) LoadServices(ctx context.Context) error {
	// 1. Register services via config with servicemanager
	// (Done via BuildFromConfig in Init)
	// if err := a.RegisterServices(); err != nil {
	// 	a.manager.Logger().Error("Failed to register services", zap.Error(err))
	// 	return err
	// }
	// 2. Init Services
	if err := a.manager.InitServices(ctx); err != nil {
		a.manager.Logger().Error("Failed to initialize services", zap.Error(err))
		return err
	}
	// 3. Start Services
	if err := a.manager.StartServices(ctx); err != nil {
		a.manager.Logger().Error("Failed to start services", zap.Error(err))
		return err
	}
	return nil
}

func (a *App) UnloadServices(ctx context.Context) error {
	if err := a.manager.StopServices(ctx); err != nil {
		a.manager.Logger().Error("Failed to stop services", zap.Error(err))
		return err
	}
	return nil
}

func (a *App) RegisterServices() error {
	// Replaced by BuildFromConfig
	return nil
}

// createService removed

func (a *App) Logger() *zap.Logger {
	return a.manager.Logger()
}

// decodeConfig removed
