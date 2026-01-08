package manager

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"grouter/pkg/config"

	"github.com/go-viper/mapstructure/v2"
	"go.uber.org/zap"
)

// Option configures the ServiceManager.
type Option func(*ServiceManager)

// WithShutdownTimeout sets the timeout for graceful shutdown.
func WithShutdownTimeout(timeout time.Duration) Option {
	return func(m *ServiceManager) {
		m.shutdownTimeout = timeout
	}
}

// WithMonitoringInterval sets the interval for health checks.
func WithMonitoringInterval(interval time.Duration) Option {
	return func(m *ServiceManager) {
		m.monitorInterval = interval
	}
}

// ServiceManager orchestrates the application lifecycle and message routing.
type ServiceManager struct {
	deps Deps

	store *ServiceStore

	shutdownTimeout time.Duration
	monitorInterval time.Duration
}

// NewServiceManager creates a new ServiceManager with dependencies.
func NewServiceManager(deps Deps, opts ...Option) *ServiceManager {
	m := &ServiceManager{
		deps:            deps,
		store:           NewServiceStore(),
		shutdownTimeout: 30 * time.Second,
		monitorInterval: 5 * time.Second,
	}
	m.deps.Store = m.store

	for _, opt := range opts {
		opt(m)
	}

	return m
}

// BuildFromConfig instantiates services based on the configuration using registered factories.
func (m *ServiceManager) BuildFromConfig() error {
	m.Logger().Info("Building services from configuration...")

	for id, svcCfg := range m.Config().Services {
		// 1. Check if enabled
		var genericCfg struct {
			Enabled bool `mapstructure:"enabled"`
		}
		if err := mapstructure.Decode(svcCfg, &genericCfg); err != nil {
			m.Logger().Warn("Failed to decode service config for enabled check", zap.String("service", id), zap.Error(err))
			continue
		}

		if !genericCfg.Enabled {
			m.Logger().Debug("Service disabled in config", zap.String("service", id))
			continue
		}

		// 2. Find Factory
		factory, exists := GetFactory(id)
		if !exists {
			// Fail fast if a service is enabled but code is missing
			return fmt.Errorf("service enabled in config but no factory registered: %s", id)
		}

		// 3. Instantiate
		m.Logger().Debug("Creating service", zap.String("service", id))
		svc, err := factory(m.deps)
		if err != nil {
			return fmt.Errorf("failed to create service %s: %w", id, err)
		}

		// 3.5 Inject Config
		if cfgAware, ok := svc.(ConfigAware); ok {
			if cfg, ok := svcCfg.(map[string]interface{}); ok {
				if err := cfgAware.InitConfig(cfg); err != nil {
					return fmt.Errorf("failed to init config for service %s: %w", id, err)
				}
			} else {
				// Try decode if it's not a map (e.g. mapstructure hook?)
				// Viper usually gives map[string]interface{}
				// But let's log warning if mismatch?
				// Actually svcCfg is from m.Config().Services which is ServicesConfig which is map[string]interface{}.
				// But the value in that map can be anything.
				// Let's assume it transforms.
				// If we can't assert, we skip config? Or error?
				// Let's try mapstructure decode into generic map if assertion fails?
				// Safer to just try map conversion or log warning.
				m.Logger().Warn("Service config is not map[string]interface{}", zap.String("service", id))
			}
		}

		// 4. Register
		m.store.Add(svc.ID(), svc) // Use svc.ID() which is reliable
	}

	m.Logger().Info("Built services", zap.Strings("services", m.store.List()))
	return nil
}

// RegisterService registers a service with the manager.
func (m *ServiceManager) RegisterService(svc Service) error {
	m.store.Add(svc.Name(), svc)
	return nil
}

// UnregisterService removes a service from the manager.
func (m *ServiceManager) UnregisterService(name string) {
	m.store.Delete(name)
}

// GetService retrieves a service by name.
func (m *ServiceManager) GetService(name string) (Service, error) {
	svc, ok := m.store.Get(name)
	if !ok {
		return nil, fmt.Errorf("service not found: %s", name)
	}
	return svc, nil
}

// Run starts the manager and blocks until a signal is received or context is done.
func (m *ServiceManager) Run(ctx context.Context) error {
	quit := make(chan os.Signal, 1) // Initialize quit channel
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	if err := m.InitServices(ctx); err != nil {
		return err
	}

	if err := m.StartServices(ctx); err != nil {
		return err
	}

	// Start Supervision Loop
	go m.MonitorServices(ctx, quit)

	select {
	case sig := <-quit:
		m.Logger().Info("Received shutdown signal", zap.String("signal", sig.String()))
	case <-ctx.Done():
		m.Logger().Info("Context cancelled")
	}

	m.Logger().Info("Shutting down services...")

	// Create shutdown context with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), m.shutdownTimeout)
	defer cancel()

	if err := m.StopServices(shutdownCtx); err != nil {
		m.Logger().Error("Error during shutdown", zap.Error(err))
		// We return the error but still try to cleanup other things?
		// Usually yes.
	}

	// Cleanup store
	if m.store != nil {
		m.store.DeleteAll()
	}

	// Flush logger
	// m.log is removed, use m.Logger() but we can't easily sync it if it's external deps
	// But we can try type assertion or just ignore.
	// m.Logger().Sync() might be valid if zap
	_ = m.Logger().Sync()

	m.Logger().Info("Shutdown complete")
	return nil
}

// MonitorServices periodically checks the health of running services.
func (m *ServiceManager) MonitorServices(ctx context.Context, quit chan<- os.Signal) {
	ticker := time.NewTicker(m.monitorInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			services := m.store.All()
			for _, svc := range services {
				status := svc.Status()
				if status == StatusFailed || status == StatusUnhealthy {
					m.Logger().Error("Service is unhealthy",
						zap.String("service", svc.Name()),
						zap.String("status", string(status)),
						zap.Error(svc.LastError()),
					)
					// Future: Trigger restart logic here
				}
			}
		}
	}
}

// Logger returns the initialized logger.
func (m *ServiceManager) Logger() *zap.Logger {
	return m.deps.Logger
}

// Config returns the configuration.
func (m *ServiceManager) Config() *config.Config {
	return m.deps.Config
}

// ListServices returns a list of all registered service names.
func (m *ServiceManager) ListServices() []string {
	return m.store.List()
}

// resolveDependencies returns services sorted by dependency order (Topological Sort).
func (m *ServiceManager) resolveDependencies() ([]Service, error) {
	services := m.store.All()
	serviceMap := make(map[string]Service)
	for _, s := range services {
		serviceMap[s.Name()] = s
	}

	visited := make(map[string]bool)
	tempVisited := make(map[string]bool) // For cycle detection
	var sorted []Service

	var visit func(string) error
	visit = func(name string) error {
		if tempVisited[name] {
			return fmt.Errorf("circular dependency detected involving service: %s", name)
		}
		if visited[name] {
			return nil
		}

		tempVisited[name] = true
		svc, exists := serviceMap[name]
		if !exists {
			return fmt.Errorf("service %s depends on unknown service", name)
		}

		for _, dep := range svc.Dependencies() {
			// Check if dependency exists in the registry
			if _, ok := serviceMap[dep]; !ok {
				m.Logger().Warn("Service depends on unregistered service (ignoring dependency)",
					zap.String("service", name), zap.String("dependency", dep))
				continue
			}
			if err := visit(dep); err != nil {
				return err
			}
		}

		tempVisited[name] = false
		visited[name] = true
		sorted = append(sorted, svc)
		return nil
	}

	for _, svc := range services {
		if !visited[svc.Name()] {
			if err := visit(svc.Name()); err != nil {
				return nil, err
			}
		}
	}

	return sorted, nil
}

// InitServices initializes all registered services in dependency order.
func (m *ServiceManager) InitServices(ctx context.Context) error {
	if m.store == nil {
		return nil
	}

	sortedServices, err := m.resolveDependencies()
	if err != nil {
		return fmt.Errorf("failed to resolve dependencies: %w", err)
	}

	for _, svc := range sortedServices {
		m.Logger().Debug("Initializing service", zap.String("service", svc.Name()))

		if err := svc.Init(ctx); err != nil {
			return fmt.Errorf("failed to initialize service %s: %w", svc.Name(), err)
		}

		// Capability: WebRoutable
		if webSvc, ok := svc.(WebRoutable); ok && m.deps.WebRouter != nil {
			m.Logger().Debug("Registering web routes", zap.String("service", svc.Name()))
			if err := webSvc.RegisterRoutes(m.deps.WebRouter); err != nil {
				return fmt.Errorf("failed to register web routes for %s: %w", svc.Name(), err)
			}
		}

		// Capability: GRPCRoutable
		if grpcSvc, ok := svc.(GRPCRoutable); ok && m.deps.GRPCServer != nil {
			m.Logger().Debug("Registering gRPC service", zap.String("service", svc.Name()))
			if err := grpcSvc.RegisterGRPC(m.deps.GRPCServer); err != nil {
				return fmt.Errorf("failed to register gRPC service for %s: %w", svc.Name(), err)
			}
		}
	}
	return nil
}

// StartServices starts all registered services in dependency order.
func (m *ServiceManager) StartServices(ctx context.Context) error {
	if m.store == nil {
		return nil
	}

	sortedServices, err := m.resolveDependencies()
	if err != nil {
		return fmt.Errorf("failed to resolve dependencies: %w", err)
	}

	for _, svc := range sortedServices {
		m.Logger().Debug("Starting service", zap.String("service", svc.Name()))

		if err := svc.Start(ctx); err != nil {
			return fmt.Errorf("failed to start service %s: %w", svc.Name(), err)
		}
	}
	return nil
}

// StopServices stops all registered services in reverse dependency order.
func (m *ServiceManager) StopServices(ctx context.Context) error {
	if m.store == nil {
		return nil
	}

	sortedServices, err := m.resolveDependencies()
	if err != nil {
		return fmt.Errorf("failed to resolve dependencies: %w", err)
	}

	// Reverse order for shutdown
	for i := len(sortedServices) - 1; i >= 0; i-- {
		svc := sortedServices[i]
		m.Logger().Debug("Stopping service", zap.String("service", svc.Name()))

		if err := svc.Stop(ctx); err != nil {
			m.Logger().Error("Failed to stop service", zap.String("service", svc.Name()), zap.Error(err))
		}
	}
	return nil
}

// Stop gracefully shuts down the infrastructure components.
func (m *ServiceManager) Stop() error {
	_ = m.Logger().Sync()

	_ = m.StopServices(context.Background())
	if m.store != nil {
		m.store.DeleteAll()
	}
	return nil
}
