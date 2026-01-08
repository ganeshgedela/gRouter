package app

import (
	"context"
	"fmt"
	"strings"

	"grouter/pkg/manager"
	messaging "grouter/pkg/messaging/nats"
	"grouter/pkg/telemetry"

	"github.com/google/uuid"

	"grouter/pkg/config"
	"grouter/pkg/logger"

	"time"

	"github.com/go-viper/mapstructure/v2"
	"go.uber.org/zap"
)

type App struct {
	manager *manager.ServiceManager
	AppId   string

	bootstrap *BootstrapService
	// health    *HealthService // Removed
	// metrics   *MetricService // Removed

	messenger *messaging.Messenger
}

func New() *App {
	return &App{}
}

const ShutdownTimeout = 15 * time.Second

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

	// 4. Initialize Dependencies (NATS, Tracing)
	// Initialize OpenTelemetry
	_, err = telemetry.InitTracer(cfg.Tracing)
	if err != nil {
		log.Error("Failed to initialize OpenTelemetry", zap.Error(err))
		return err
	}

	// Initialize Messenger
	if cfg.NATS.Enabled {
		messenger := &messaging.Messenger{}
		log.Debug("Initializing NATS Messenger", zap.String("URL", cfg.NATS.URL))
		if err := messenger.InitFromConfig(cfg, log); err != nil {
			log.Error("Failed to initialize NATS Messenger", zap.Error(err))
			return err
		}
		a.messenger = messenger
		deps.Messenger = messenger
	}

	// 5. Create Service Manager
	a.manager = manager.NewServiceManager(deps, manager.WithShutdownTimeout(ShutdownTimeout))

	// 6. Build Services from Config
	if err := a.manager.BuildFromConfig(); err != nil {
		return err
	}

	// 7. Initialize Internal Services (that are not in config or need special handling)
	// Generate unique AppId
	a.AppId = cfg.App.Name + "-" + strings.Split(uuid.New().String(), "-")[0]
	log.Debug("App initialized", zap.String("AppId", a.AppId))

	// Initialize Startup Service
	if err := a.initStartupService(context.Background()); err != nil {
		return err
	}
	return nil
}

func (a *App) GetAppName() string {
	return a.manager.Config().App.Name
}

func (a *App) initStartupService(ctx context.Context) error {

	// 2. start messenger
	if a.messenger != nil {
		if err := a.messenger.Start(); err != nil {
			return fmt.Errorf("failed to start messenger: %w", err)
		}

		// 3. register bootstrap service with servicemanager
		if err := a.RegisterBootstrap(ctx); err != nil {
			return err
		}
	}

	return nil
}

// RegisterHealth and RegisterMetrics removal

func (a *App) RegisterBootstrap(ctx context.Context) error {
	// Initialize Bootstrap Service
	a.bootstrap = NewBootstrapService(a.messenger, a.manager.Logger(), a.GetAppName())

	// 1. Init bootstrap service explicitly as it's not part of standard lifecycle until started
	if err := a.bootstrap.Init(ctx); err != nil {
		return err
	}
	// 2. subscribe to bootstrap service with nats
	if err := a.bootstrap.Start(ctx); err != nil {
		return err
	}
	return nil
}

func (a *App) Start(ctx context.Context) error {
	logger := a.manager.Logger()

	logger.Info("Send NATS message to " + a.GetAppName() + ".start to begin.")

	for {
		// Block until start message is received
		select {
		case <-a.bootstrap.WaitForStart():
			// Drain any stale stop signals received while backend was stopped
			select {
			case <-a.bootstrap.WaitForStop():
				logger.Debug("Drained stale stop signal")
			default:
			}

			logger.Debug("Start signal received. Registering configured services...")
			if err := a.LoadServices(ctx); err != nil {
				logger.Error("Failed to load services", zap.Error(err))
				return err
			}
			logger.Info("Application is ready to accept requests...")
		case <-ctx.Done():
			return ctx.Err()
		}

		// Block until stop message is received
		select {
		case <-a.bootstrap.WaitForStop():
			// Drain any stale start signals received while backend was running
			select {
			case <-a.bootstrap.WaitForStart():
				logger.Debug("Drained stale start signal")
			default:
			}

			logger.Debug("Stop signal received. Unregistering services...")
			if err := a.UnloadServices(ctx); err != nil {
				logger.Error("Failed to unload services", zap.Error(err))
			}
			logger.Info("Services stopped. Waiting for start signal...")
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (a *App) LoadServices(ctx context.Context) error {
	logger := a.manager.Logger()

	// 1. & 2. Health and Metrics are now loaded via BuildFromConfig (if enabled in config)
	// We no longer manually register them here.

	// 3. Active Orchestration: Init and Start
	if err := a.manager.InitServices(ctx); err != nil {
		logger.Error("Failed to initialize services", zap.Error(err))
		return err
	}
	// 4. Start Services internals
	if err := a.manager.StartServices(ctx); err != nil {
		logger.Error("Failed to start services", zap.Error(err))
		return err
	}
	return nil
}
func (a *App) UnloadServices(ctx context.Context) error {

	if err := a.manager.StopServices(ctx); err != nil {
		a.manager.Logger().Error("Failed to stop services", zap.Error(err))
		return err
	}

	if err := a.UnregisterServices(); err != nil {
		a.manager.Logger().Error("Failed to unregister services", zap.Error(err))
		return err
	}

	return nil
}

func (a *App) UnregisterServices() error {

	logger := a.manager.Logger()

	services := a.manager.ListServices()

	for _, service := range services {
		logger.Debug("Unregistering service: " + service)
		a.manager.UnregisterService(service)

	}

	logger.Debug("Services: ", zap.Any("services", a.manager.ListServices()))
	return nil
}

func (a *App) Stop(ctx context.Context) error {
	return a.UnloadServices(ctx)
}

func (a *App) RegisterServices() error {
	// Replaced by BuildFromConfig
	return nil
}

// createService removed, replaced by Factory pattern

func (a *App) Logger() *zap.Logger {
	return a.manager.Logger()
}

func (a *App) decodeConfig(input interface{}, output interface{}) error {
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
