package app

import (
	"context"
	"strings"

	"grouter/pkg/manager"
	"grouter/services/natsdemosvc/internal/pkg/natdemo"

	"github.com/google/uuid"

	"github.com/go-viper/mapstructure/v2"
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
		startChan: make(chan struct{}, 1),
		stopChan:  make(chan struct{}, 1),
	}
}

func (a *App) Init() error {
	if err := a.manager.Init(); err != nil {
		return err
	}
	if err := a.manager.InitNATS(); err != nil {
		return err
	}
	if err := a.manager.InitWebServer(); err != nil {
		return err
	}
	// Generate unique AppId
	a.AppId = a.manager.Config().App.Name + "-" + strings.Split(uuid.New().String(), "-")[0]
	a.manager.Logger().Info("App initialized", zap.String("AppId", a.AppId))

	if err := a.InitAppStartupServices(); err != nil {
		a.manager.Logger().Error("Failed to initialize startup services", zap.Error(err))
		return err
	}
	return nil
}

func (a *App) GetAppName() string {
	return a.manager.Config().App.Name
}

func (a *App) RegisterBootstrap() error {
	logger := a.manager.Logger()
	// Register Bootstrap Service to listen for start signal
	bootstrap := NewBootstrapService(a.startChan)
	if err := a.manager.RegisterService(bootstrap); err != nil {
		return err
	}
	subject := a.GetAppName() + ".start"
	logger.Info("Registering Bootstrap Service to listen for start signal on topic " + subject)
	if err := a.manager.SubscribeToTopics(subject, ""); err != nil {
		return err
	}
	return nil
}

func (a *App) RegisterStop() error {
	logger := a.manager.Logger()
	// Register Stop Service to listen for stop signal
	stopSvc := NewStopService(a.stopChan)
	if err := a.manager.RegisterService(stopSvc); err != nil {
		return err
	}
	subject := a.GetAppName() + ".stop"
	logger.Info("Registering Stop Service to listen for stop signal on topic " + subject)
	if err := a.manager.SubscribeToTopics(subject, ""); err != nil {
		return err
	}
	return nil
}

func (a *App) RegisterHealth() error {
	logger := a.manager.Logger()
	// Register Health Service to listen for health signal
	healthSvc := NewHealthService(a)
	if err := a.manager.RegisterService(healthSvc); err != nil {
		return err
	}
	subject := a.GetAppName() + ".health.>"
	logger.Info("Registering Health Service to listen for health signal on topic " + subject)
	if err := a.manager.SubscribeToTopics(subject, ""); err != nil {
		return err
	}
	return nil
}

func (a *App) InitAppStartupServices() error {
	logger := a.manager.Logger()

	if err := a.RegisterBootstrap(); err != nil {
		logger.Error("Failed to register bootstrap service", zap.Error(err))
		return err
	}
	if err := a.RegisterStop(); err != nil {
		logger.Error("Failed to register stop service", zap.Error(err))
		return err
	}
	if err := a.RegisterHealth(); err != nil {
		logger.Error("Failed to register health service", zap.Error(err))
		return err
	}
	return nil
}

func (a *App) Start(ctx context.Context) error {
	logger := a.manager.Logger()
	appName := a.GetAppName()

	logger.Info("Starting " + appName + "...")
	logger.Info("Send NATS message to " + appName + ".start to begin.")

	// Start Manager (begins message listening)
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
			logger.Info("Send NATS message to " + appName + ".stop to stop.")
			logger.Info("Services registered. Application is ready to accept requests...")
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
			logger.Info("Services stopped. Waiting for start signal...")
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
	return make(chan struct{})
}

func (a *App) RegisterServices() error {

	cfg := a.manager.Config()
	logger := a.manager.Logger()

	topic := a.GetAppName() + ".>"
	if err := a.manager.SubscribeToTopics(topic, cfg.App.Name); err != nil {
		return err
	}

	logger.Info("Registering service: " + a.GetAppName() + " to topic: " + topic)

	// build services list from cfg.Services
	for name, serviceCfg := range cfg.Services {

		if name == "natdemo" {
			var natConfig natdemo.NATDemoConfig
			if err := decodeConfig(serviceCfg, &natConfig); err != nil {
				logger.Error("Failed to decode NATDemo config", zap.Error(err))
				return err
			}

			if natConfig.Enabled {
				natModule := natdemo.NewNATDemo(a.manager.Publisher(), logger, natConfig)
				if err := a.manager.RegisterService(natModule); err != nil {
					logger.Error("Failed to register NATDemo Module", zap.Error(err))
					return err
				}
			}
		}
	}
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
