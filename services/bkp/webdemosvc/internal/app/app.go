package app

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"grouter/pkg/health"
	"grouter/pkg/manager"
	messaging "grouter/pkg/messaging/nats"
	"grouter/pkg/telemetry"
	"grouter/pkg/web"

	"grouter/pkg/config"
	"grouter/pkg/logger"

	"go.uber.org/zap"
)

type App struct {
	manager *manager.ServiceManager
	AppId   string

	startChan chan struct{}
	stopChan  chan struct{}

	messenger      *messaging.Messenger
	webServer      *web.Server
	healthSvc      *health.HealthService
	tracerShutdown func(context.Context) error
}

func New() *App {
	return &App{
		startChan: make(chan struct{}),
		stopChan:  make(chan struct{}),
	}
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

	// 4. Initialize Dependencies
	// Initialize OpenTelemetry
	shutdown, err := telemetry.InitTracer(cfg.Tracing)
	if err != nil {
		return fmt.Errorf("failed to init tracer: %w", err)
	}
	a.tracerShutdown = shutdown

	// Initialize Health Service
	a.healthSvc = health.NewHealthService()

	// Initialize Messenger
	if cfg.NATS.Enabled {
		a.messenger = &messaging.Messenger{}
		// Fix: pass the config pointer directly
		if err := a.messenger.InitFromConfig(cfg, log); err != nil {
			return fmt.Errorf("failed to init nats: %w", err)
		}
		deps.Messenger = a.messenger
	}

	// 5. Create Service Manager
	a.manager = manager.NewServiceManager(deps)

	// Initialize Web Server
	if a.manager.Config().Web.Enabled {
		webConfig := web.Config{
			Port:            a.manager.Config().Web.Port,
			ReadTimeout:     a.manager.Config().Web.ReadTimeout,
			WriteTimeout:    a.manager.Config().Web.WriteTimeout,
			ShutdownTimeout: a.manager.Config().Web.ShutdownTimeout,
			Mode:            a.manager.Config().Web.Mode,
			Metrics: web.MetricsConfig{
				Enabled: a.manager.Config().Web.Metrics.Enabled,
				Path:    a.manager.Config().Web.Metrics.Path,
			},
			Tracing: web.TracingConfig{
				Enabled:     a.manager.Config().Tracing.Enabled,
				ServiceName: a.manager.Config().Tracing.ServiceName,
			},
			TLS: web.TLSConfig{
				Enabled:  a.manager.Config().Web.TLS.Enabled,
				CertFile: a.manager.Config().Web.TLS.CertFile,
				KeyFile:  a.manager.Config().Web.TLS.KeyFile,
			},
			CORS: web.CORSConfig{
				Enabled:          a.manager.Config().Web.CORS.Enabled,
				AllowedOrigins:   a.manager.Config().Web.CORS.AllowedOrigins,
				AllowedMethods:   a.manager.Config().Web.CORS.AllowedMethods,
				AllowedHeaders:   a.manager.Config().Web.CORS.AllowedHeaders,
				ExposedHeaders:   a.manager.Config().Web.CORS.ExposedHeaders,
				AllowCredentials: a.manager.Config().Web.CORS.AllowCredentials,
				MaxAge:           a.manager.Config().Web.CORS.MaxAge,
			},
			Security: web.SecurityConfig{
				Enabled:               a.manager.Config().Web.Security.Enabled,
				XSSProtection:         a.manager.Config().Web.Security.XSSProtection,
				ContentTypeNosniff:    a.manager.Config().Web.Security.ContentTypeNosniff,
				XFrameOptions:         a.manager.Config().Web.Security.XFrameOptions,
				HSTSMaxAge:            a.manager.Config().Web.Security.HSTSMaxAge,
				HSTSExcludeSubdomains: a.manager.Config().Web.Security.HSTSExcludeSubdomains,
				ContentSecurityPolicy: a.manager.Config().Web.Security.ContentSecurityPolicy,
				ReferrerPolicy:        a.manager.Config().Web.Security.ReferrerPolicy,
				CustomHeaders:         a.manager.Config().Web.Security.CustomHeaders,
			},
			RateLimit: web.RateLimitConfig{
				Enabled:           a.manager.Config().Web.RateLimit.Enabled,
				RequestsPerSecond: a.manager.Config().Web.RateLimit.RequestsPerSecond,
				Burst:             a.manager.Config().Web.RateLimit.Burst,
			},
			Swagger: web.SwaggerConfig{
				Enabled: a.manager.Config().Web.Swagger.Enabled,
				Path:    a.manager.Config().Web.Swagger.Path,
			},
			Logging: web.LoggingConfig{
				Enabled: a.manager.Config().Web.Logging.Enabled,
			},
			Auth: web.AuthConfig{
				Enabled:  a.manager.Config().Web.Auth.Enabled,
				Issuer:   a.manager.Config().Web.Auth.Issuer,
				Audience: a.manager.Config().Web.Auth.Audience,
			},
		}
		a.webServer = web.NewWebServer(webConfig, log, a.healthSvc)
		if err := a.webServer.Start(); err != nil {
			return fmt.Errorf("failed to init web server: %w", err)
		}
	}

	// 7. Build Services from Config
	if err := a.manager.BuildFromConfig(); err != nil {
		return err
	}

	// 8. Register Web Services to Web Server (since BuildFromConfig registered them to Manager)
	a.ReRegisterServices()

	// Generate unique AppId
	a.AppId = cfg.App.Name + "-" + uuid.New().String()
	log.Info("App initialized", zap.String("AppId", a.AppId))

	// Register Health Service (Legacy internal registration, ideally would be factory too)
	healthSvc := NewHealthService(a.healthSvc, cfg.App.Name)
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
	stopSvc := NewStopService(a.stopChan, a.webServer)
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

	for {
		// Block until start message is received
		select {
		case <-a.startChan:
			logger.Info("Start signal received. Registering services...")
			// Register services via config (replaced by BuildFromConfig, but maybe we want dynamic reload?)
			// For now, assume services are already built. Or rebuild?
			// The original code re-registered services on start signal.
			// Is config reloading expected?
			// If so, we'd need to reload config and call BuildFromConfig again?
			// For now, just re-register routes is handled below.
			if err := a.RegisterServices(); err != nil {
				logger.Error("Failed to register services", zap.Error(err))
			}

			// Reload web server to apply new services
			logger.Info("Reloading web server to apply new services...")
			// Note: This might cause a race with the /start response if not handled carefully,
			// but for now it's necessary to register routes on the running engine.
			// ResetEngine stops the server first.
			if err := a.webServer.ResetEngine(context.Background()); err != nil {
				logger.Error("Failed to reset engine", zap.Error(err))
			}
			a.ReRegisterServices()
			if err := a.webServer.Start(); err != nil {
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
			if err := a.webServer.ResetEngine(context.Background()); err != nil {
				logger.Error("Failed to reset engine", zap.Error(err))
			}
			a.ReRegisterServices()
			if err := a.webServer.Start(); err != nil {
				logger.Error("Failed to start web server", zap.Error(err))
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (a *App) Stop(ctx context.Context) error {
	m := a.manager
	m.Logger().Info("Stopping gRouter service")

	if a.messenger != nil {
		if err := a.messenger.Close(); err != nil {
			m.Logger().Error("Failed to close messenger", zap.Error(err))
		}
	}
	if a.webServer != nil {
		if err := a.webServer.Stop(ctx); err != nil {
			m.Logger().Error("Failed to stop web server", zap.Error(err))
		}
	}

	if a.tracerShutdown != nil {
		if err := a.tracerShutdown(ctx); err != nil {
			m.Logger().Warn("Failed to shutdown tracer", zap.Error(err))
		}
	}
	return a.manager.Stop()
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
	// Replaced by BuildFromConfig.
	// If dynamic reload is needed, we might call a.manager.BuildFromConfig() here again?
	// But duplicate registration might error.
	// Manager should handle re-registration or we clear first.
	// Original code cleared services on Stop.
	// So on Start, we need to Re-Register.
	// BuildFromConfig adds to store.
	// If store was cleared, BuildFromConfig is correct.
	return a.manager.BuildFromConfig()
}

func (a *App) ReRegisterServices() {
	for _, name := range a.manager.ListServices() {
		svc, err := a.manager.GetService(name)
		if err != nil {
			continue
		}
		if webSvc, ok := svc.(web.WebService); ok && a.webServer != nil {
			a.webServer.RegisterWebService(webSvc)
		}
	}
}

func (a *App) Logger() *zap.Logger {
	return a.manager.Logger()
}

// decodeConfig removed
