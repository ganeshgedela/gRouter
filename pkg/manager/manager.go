package manager

import (
	"context"
	"fmt"
	"time"

	"grouter/pkg/config"
	"grouter/pkg/health"
	"grouter/pkg/logger"
	messaging "grouter/pkg/messaging/nats"
	"grouter/pkg/telemetry"
	"grouter/pkg/web"

	"go.uber.org/zap"
)

// ServiceManager orchestrates the application lifecycle and message routing.
type ServiceManager struct {
	cfg *config.Config
	log *zap.Logger

	router *ServiceRouter

	messenger *messaging.Messenger

	webServer *web.Server

	health  *health.HealthService
	timeout time.Duration

	// Cleanup for OpenTelemetry
	tracerShutdown func(context.Context) error
}

// NewServiceManager creates a new ServiceManager with default settings.
func NewServiceManager() *ServiceManager {
	return &ServiceManager{
		router:  NewServiceRouter(),
		timeout: 10 * time.Second,
	}
}

// Init initializes configuration, logger, NATS, and registers services.
func (m *ServiceManager) Init() error {
	if err := m.initConfig(); err != nil {
		return err
	}
	if err := m.initLogger(); err != nil {
		return err
	}

	// Initialize OpenTelemetry
	shutdown, err := telemetry.InitTracer(m.cfg.Tracing)
	if err != nil {
		return fmt.Errorf("failed to initialize tracer: %w", err)
	}
	m.tracerShutdown = shutdown

	m.log.Info("Initializing gRouter service",
		zap.String("name", m.cfg.App.Name),
		zap.String("version", m.cfg.App.Version),
		zap.String("environment", m.cfg.App.Environment),
	)

	// Register health service
	m.health = health.NewHealthService()

	return nil
}

func (m *ServiceManager) initConfig() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}
	m.cfg = cfg
	return nil
}

func (m *ServiceManager) initLogger() error {
	if m.cfg == nil {
		return fmt.Errorf("init logger: config is nil")
	}
	log, err := logger.New(logger.Config{
		Level:      m.cfg.Log.Level,
		Format:     m.cfg.Log.Format,
		OutputPath: m.cfg.Log.OutputPath,
	})
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}
	m.log = log
	return nil
}

func (m *ServiceManager) InitNATS() error {
	if m.cfg == nil || m.log == nil {
		return fmt.Errorf("init nats: config or logger is nil")
	}

	if !m.cfg.NATS.Enabled {
		//m.log.Info("NATS disabled")
		return nil
	}

	// Initialize Messenger
	m.messenger = &messaging.Messenger{}
	if err := m.messenger.Init(messaging.Config{
		URL:               m.cfg.NATS.URL,
		MaxReconnects:     m.cfg.NATS.MaxReconnects,
		ReconnectWait:     m.cfg.NATS.ReconnectWait,
		ConnectionTimeout: m.cfg.NATS.ConnectionTimeout,
		Token:             m.cfg.NATS.Token,
		Username:          m.cfg.NATS.Username,
		Password:          m.cfg.NATS.Password,
		CredsFile:         m.cfg.NATS.CredsFile,
		UseTLS:            m.cfg.NATS.UseTLS,
		SkipVerify:        m.cfg.NATS.SkipVerify,
		CAFile:            m.cfg.NATS.CAFile,
		CertFile:          m.cfg.NATS.CertFile,
		KeyFile:           m.cfg.NATS.KeyFile,
		Metrics: messaging.MetricsConfig{
			Enabled: m.cfg.NATS.Metrics.Enabled,
			Path:    m.cfg.NATS.Metrics.Path,
		},
		Logging: messaging.LoggingConfig{
			Enabled: m.cfg.NATS.Logging.Enabled,
		},
		Tracing: messaging.TracingConfig{
			Enabled: m.cfg.Tracing.Enabled,
		},
	}, m.log, m.cfg.App.Name); err != nil {
		return fmt.Errorf("failed to initialize messenger: %w", err)
	}

	m.log.Info("NATS initialized via Messenger",
		zap.String("url", m.cfg.NATS.URL),
		zap.String("app", m.cfg.App.Name),
	)

	return nil
}

func (m *ServiceManager) InitWebServer() error {
	if m.cfg == nil || m.log == nil {
		return fmt.Errorf("init web server: config or logger is nil")
	}

	if !m.cfg.Web.Enabled {
		m.log.Info("Web server disabled")
		return nil
	}

	webConfig := web.Config{
		Port:            m.cfg.Web.Port,
		ReadTimeout:     m.cfg.Web.ReadTimeout,
		WriteTimeout:    m.cfg.Web.WriteTimeout,
		ShutdownTimeout: m.cfg.Web.ShutdownTimeout,
		Mode:            m.cfg.Web.Mode,
		Metrics: web.MetricsConfig{
			Enabled: m.cfg.Web.Metrics.Enabled,
			Path:    m.cfg.Web.Metrics.Path,
		},
		Tracing: web.TracingConfig{
			Enabled:     m.cfg.Tracing.Enabled,
			ServiceName: m.cfg.Tracing.ServiceName,
		},
		TLS: web.TLSConfig{
			Enabled:  m.cfg.Web.TLS.Enabled,
			CertFile: m.cfg.Web.TLS.CertFile,
			KeyFile:  m.cfg.Web.TLS.KeyFile,
		},
		CORS: web.CORSConfig{
			Enabled:          m.cfg.Web.CORS.Enabled,
			AllowedOrigins:   m.cfg.Web.CORS.AllowedOrigins,
			AllowedMethods:   m.cfg.Web.CORS.AllowedMethods,
			AllowedHeaders:   m.cfg.Web.CORS.AllowedHeaders,
			ExposedHeaders:   m.cfg.Web.CORS.ExposedHeaders,
			AllowCredentials: m.cfg.Web.CORS.AllowCredentials,
			MaxAge:           m.cfg.Web.CORS.MaxAge,
		},
		Security: web.SecurityConfig{
			Enabled:               m.cfg.Web.Security.Enabled,
			XSSProtection:         m.cfg.Web.Security.XSSProtection,
			ContentTypeNosniff:    m.cfg.Web.Security.ContentTypeNosniff,
			XFrameOptions:         m.cfg.Web.Security.XFrameOptions,
			HSTSMaxAge:            m.cfg.Web.Security.HSTSMaxAge,
			HSTSExcludeSubdomains: m.cfg.Web.Security.HSTSExcludeSubdomains,
			ContentSecurityPolicy: m.cfg.Web.Security.ContentSecurityPolicy,
			ReferrerPolicy:        m.cfg.Web.Security.ReferrerPolicy,
			CustomHeaders:         m.cfg.Web.Security.CustomHeaders,
		},
		RateLimit: web.RateLimitConfig{
			Enabled:           m.cfg.Web.RateLimit.Enabled,
			RequestsPerSecond: m.cfg.Web.RateLimit.RequestsPerSecond,
			Burst:             m.cfg.Web.RateLimit.Burst,
		},
		Swagger: web.SwaggerConfig{
			Enabled: m.cfg.Web.Swagger.Enabled,
			Path:    m.cfg.Web.Swagger.Path,
		},
	}
	m.webServer = web.NewWebServer(webConfig, m.log, m.health)

	// Start web server
	if err := m.webServer.Start(); err != nil {
		return fmt.Errorf("failed to start web server: %w", err)
	}

	return nil
}

// RegisterService registers a service with the manager.
// It automatically detects and registers capabilities (Web, NATS).
func (m *ServiceManager) RegisterService(svc Service) error {
	if svc == nil {
		return nil
	}
	m.router.Register(svc.Name(), svc)

	// Check for Web Capability
	if m.webServer != nil {
		if webSvc, ok := svc.(web.WebService); ok {
			m.webServer.RegisterWebService(webSvc)
		}
	}

	return nil
}

// ReRegisterServices iterates over all currently defined services and re-registers them.
// This is useful during a restart to ensure all services are active.
func (m *ServiceManager) ReRegisterServices() {
	for _, serviceName := range m.ListServices() {
		if svc, ok := m.GetService(serviceName); ok {
			m.RegisterService(svc)
		}
	}
}

// UnregisterService removes a service from the manager.
func (m *ServiceManager) UnregisterService(name string) {
	m.router.Unregister(name)
}

// Logger returns the initialized logger.
func (m *ServiceManager) Logger() *zap.Logger {
	return m.log
}

// Publisher returns the initialized NATS publisher.
func (m *ServiceManager) Publisher() messaging.Publisher {
	return m.messenger.Publisher
}

// Messenger returns the initialized Messenger instance
func (m *ServiceManager) Messenger() *messaging.Messenger {
	return m.messenger
}

func (m *ServiceManager) Config() *config.Config {
	return m.cfg
}

// Health returns the shared HealthService instance
func (m *ServiceManager) Health() *health.HealthService {
	return m.health
}

func (m *ServiceManager) WebServer() *web.Server {
	return m.webServer
}

func (m *ServiceManager) ListServices() []string {
	return m.router.store.List()
}

func (m *ServiceManager) GetService(name string) (Service, bool) {
	return m.router.store.Get(name)
}

// Start begins listening for messages on the configured topics.
func (m *ServiceManager) Start(ctx context.Context) error {
	m.log.Debug("ServiceManager started successfully")
	return nil
}

func (m *ServiceManager) onNATSMessage(ctx context.Context, subject string, env *messaging.MessageEnvelope) error {
	m.log.Debug("Received message",
		zap.String("subject", subject),
		zap.String("type", env.Type),
		zap.String("id", env.ID),
	)
	//topic := strings.TrimPrefix(subject, m.cfg.App.Name+".")
	topic := env.Type
	err := m.router.HandleMessage(ctx, topic, env)
	if err != nil {
		m.log.Error("HandleMessage failed",
			zap.Error(err),
			zap.String("topic", topic),
			zap.String("id", env.ID),
		)
		if env.Reply != "" && m.messenger != nil && m.messenger.Publisher != nil {
			return m.messenger.Publisher.PublishError(ctx, env.Reply, err.Error())
		}
		return nil
	}

	return nil
}

// replyError is deprecated. Use m.messenger.Publisher.PublishError instead.
// Keeping it removed.

// Stop gracefully shuts down the manager and its components.
func (m *ServiceManager) Stop(ctx context.Context) error {
	m.log.Info("Stopping gRouter service")

	if m.messenger != nil {
		if err := m.messenger.Close(); err != nil {
			m.log.Error("Failed to close messenger", zap.Error(err))
		}
	}
	if m.webServer != nil {
		if err := m.webServer.Stop(ctx); err != nil {
			m.log.Error("Failed to stop web server", zap.Error(err))
		}
	}
	if m.log != nil {
		_ = m.log.Sync()
	}

	if m.tracerShutdown != nil {
		if err := m.tracerShutdown(ctx); err != nil {
			m.log.Warn("Failed to shutdown tracer", zap.Error(err))
		}
	}
	return nil
}

func (m *ServiceManager) SubscribeToTopics(topic string, queueGroup string) error {
	m.log.Info("Subscribing to topics", zap.String("topic", topic))

	if m.messenger == nil {
		m.log.Warn("NATS disabled or messenger not initialized, skipping subscription", zap.String("topic", topic))
		return nil
	}

	if err := m.messenger.Subscriber.Subscribe(
		topic,
		m.onNATSMessage,
		&messaging.SubscribeOptions{
			QueueGroup: queueGroup,
		}); err != nil {
		return fmt.Errorf("failed to subscribe: %w", err)
	}

	return nil
}
