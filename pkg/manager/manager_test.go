package manager

import (
	"context"
	"fmt"
	"testing"
	"time"

	"grouter/pkg/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// MockService matches Phase 3 Service interface
type MockService struct {
	mock.Mock
	id   string
	deps []string
	// Allow setting status behavior
	status  HealthStatus
	lastErr error
}

func (m *MockService) ID() string {
	return m.id
}

func (m *MockService) Name() string {
	return m.id
}

func (m *MockService) Init(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockService) Start(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockService) Stop(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockService) Status() HealthStatus {
	args := m.Called()
	if args.Get(0) == nil {
		return m.status
	}
	return args.Get(0).(HealthStatus)
}

func (m *MockService) LastError() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockService) Dependencies() []string {
	return m.deps
}

// ConfigAware Mock
type MockConfigAwareService struct {
	MockService
}

func (m *MockConfigAwareService) InitConfig(cfg map[string]interface{}) error {
	args := m.Called(cfg)
	return args.Error(0)
}

func TestNewServiceManager(t *testing.T) {
	logger := zap.NewNop()
	deps := Deps{
		Logger: logger,
		Config: &config.Config{},
	}
	mgr := NewServiceManager(deps)

	assert.NotNil(t, mgr)
	assert.Equal(t, logger, mgr.Logger())
	assert.NotNil(t, mgr.store)
	assert.Equal(t, mgr.store, mgr.deps.Store) // Verify Store injection into Deps
}

func TestServiceManager_BuildFromConfig(t *testing.T) {
	// Setup
	logger := zap.NewNop()
	cfg := &config.Config{
		Services: map[string]interface{}{
			"mock_svc": map[string]interface{}{
				"enabled": true,
				"foo":     "bar",
			},
		},
	}
	deps := Deps{Logger: logger, Config: cfg}
	mgr := NewServiceManager(deps)

	// Register Factory
	calledFactory := false
	RegisterFactory("mock_svc", func(d Deps) (Service, error) {
		calledFactory = true
		assert.Equal(t, deps.Config, d.Config)
		assert.Equal(t, deps.Logger, d.Logger)
		assert.NotNil(t, d.Store) // Store is injected
		ms := &MockConfigAwareService{MockService{id: "mock_svc"}}
		ms.On("InitConfig", mock.Anything).Return(nil)
		return ms, nil
	})

	// Execute
	err := mgr.BuildFromConfig()

	// Verify
	assert.NoError(t, err)
	assert.True(t, calledFactory)
	assert.True(t, mgr.store.Exists("mock_svc"))
}

func TestServiceManager_DependencyResolution(t *testing.T) {
	// A -> B -> C
	sA := &MockService{id: "A", deps: []string{"B"}}
	sB := &MockService{id: "B", deps: []string{"C"}}
	sC := &MockService{id: "C", deps: []string{}}

	logger := zap.NewNop()
	mgr := NewServiceManager(Deps{Logger: logger})
	mgr.RegisterService(sA)
	mgr.RegisterService(sB)
	mgr.RegisterService(sC)

	sorted, err := mgr.resolveDependencies()
	assert.NoError(t, err)
	assert.Len(t, sorted, 3)
	assert.Equal(t, "C", sorted[0].Name())
	assert.Equal(t, "B", sorted[1].Name())
	assert.Equal(t, "A", sorted[2].Name())
}

func TestServiceManager_CycleDetection(t *testing.T) {
	// A -> B -> A
	sA := &MockService{id: "A", deps: []string{"B"}}
	sB := &MockService{id: "B", deps: []string{"A"}}

	logger := zap.NewNop()
	mgr := NewServiceManager(Deps{Logger: logger})
	mgr.RegisterService(sA)
	mgr.RegisterService(sB)

	_, err := mgr.resolveDependencies()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "circular dependency")
}

func TestServiceManager_Lifecycle(t *testing.T) {
	sA := &MockService{id: "A"}
	sA.On("Init", mock.Anything).Return(nil)
	sA.On("Start", mock.Anything).Return(nil)
	sA.On("Stop", mock.Anything).Return(nil)

	logger := zap.NewNop()
	mgr := NewServiceManager(Deps{Logger: logger})
	mgr.RegisterService(sA)

	ctx := context.Background()
	assert.NoError(t, mgr.InitServices(ctx))
	assert.NoError(t, mgr.StartServices(ctx))
	assert.NoError(t, mgr.StopServices(ctx))

	sA.AssertExpectations(t)
}

// MockWebRoutable Service
type MockWebRoutableService struct {
	MockService
}

func (m *MockWebRoutableService) RegisterRoutes(router interface{}) error {
	args := m.Called(router)
	return args.Error(0)
}

// MockGRPCRoutable Service
type MockGRPCRoutableService struct {
	MockService
}

func (m *MockGRPCRoutableService) RegisterGRPC(server interface{}) error {
	args := m.Called(server)
	return args.Error(0)
}

func TestServiceManager_AutoWiring(t *testing.T) {
	// Setup Deps with Mock Router/Server
	mockWebRouter := "mockGinEngine"
	mockGRPCServer := "mockGRPCServer"
	deps := Deps{
		Logger:     zap.NewNop(),
		WebRouter:  mockWebRouter,
		GRPCServer: mockGRPCServer,
	}
	mgr := NewServiceManager(deps)

	// Web Service
	webSvc := &MockWebRoutableService{MockService{id: "web"}}
	webSvc.On("Init", mock.Anything).Return(nil)
	webSvc.On("RegisterRoutes", mockWebRouter).Return(nil)
	mgr.RegisterService(webSvc)

	// GRPC Service
	grpcSvc := &MockGRPCRoutableService{MockService{id: "grpc"}}
	grpcSvc.On("Init", mock.Anything).Return(nil)
	grpcSvc.On("RegisterGRPC", mockGRPCServer).Return(nil)
	mgr.RegisterService(grpcSvc)

	// Execute InitServices (where auto-wiring happens)
	err := mgr.InitServices(context.Background())

	// Verify
	assert.NoError(t, err)
	webSvc.AssertExpectations(t)
	webSvc.AssertExpectations(t)
	grpcSvc.AssertExpectations(t)
}

func TestServiceManager_Wrappers(t *testing.T) {
	mgr := NewServiceManager(Deps{Logger: zap.NewNop()})
	svc := &MockService{id: "test"}

	// Register
	mgr.RegisterService(svc)
	assert.Contains(t, mgr.ListServices(), "test")

	// Get
	got, err := mgr.GetService("test")
	assert.NoError(t, err)
	assert.Equal(t, svc, got)

	_, err = mgr.GetService("missing")
	assert.Error(t, err)

	// Unregister
	mgr.UnregisterService("test")
	assert.NotContains(t, mgr.ListServices(), "test")
}

func TestServiceManager_BuildFromConfig_Errors(t *testing.T) {
	logger := zap.NewNop()

	// Case 1: Service enabled but no factory
	cfg := &config.Config{
		Services: map[string]interface{}{
			"unknown": map[string]interface{}{"enabled": true},
		},
	}
	mgr := NewServiceManager(Deps{Logger: logger, Config: cfg})
	err := mgr.BuildFromConfig()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no factory registered")

	// Case 2: Factory returns error
	RegisterFactory("fail_svc", func(d Deps) (Service, error) {
		return nil, assert.AnError
	})
	cfg2 := &config.Config{
		Services: map[string]interface{}{
			"fail_svc": map[string]interface{}{"enabled": true},
		},
	}
	mgr2 := NewServiceManager(Deps{Logger: logger, Config: cfg2})
	err = mgr2.BuildFromConfig()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create service")
}

func TestServiceManager_Run(t *testing.T) {
	svc := &MockService{id: "run_svc"}
	svc.On("Init", mock.Anything).Return(nil)
	svc.On("Start", mock.Anything).Return(nil)
	svc.On("Stop", mock.Anything).Return(nil)
	svc.On("Status").Return(StatusUnhealthy)
	svc.On("LastError").Return(fmt.Errorf("random glitch"))

	logger := zap.NewNop()
	mgr := NewServiceManager(Deps{Logger: logger, Config: &config.Config{}}, WithMonitoringInterval(10*time.Millisecond))
	mgr.RegisterService(svc)

	ctx, cancel := context.WithCancel(context.Background())

	// Run in goroutine
	errCh := make(chan error)
	go func() {
		errCh <- mgr.Run(ctx)
	}()

	// Allow some time for startup
	time.Sleep(50 * time.Millisecond)

	// Cancel context to trigger shutdown
	cancel()

	err := <-errCh
	assert.NoError(t, err)
	svc.AssertExpectations(t)
}

func TestServiceManager_Stop_Timeout(t *testing.T) {
	// Setup manager with short timeout
	mgr := NewServiceManager(Deps{Logger: zap.NewNop()}, WithShutdownTimeout(1*time.Millisecond))

	// We can't easily verify the timeout logic without a blocking service stop,
	// but we can verify the option sets the field if we could access it or rely on Run usage.
	// For coverage, we just exercise Stop()
	err := mgr.Stop()
	assert.NoError(t, err)
}

func TestServiceManager_Lifecycle_Errors(t *testing.T) {
	logger := zap.NewNop()

	// Init Fails
	sInitFail := &MockService{id: "init_fail"}
	sInitFail.On("Init", mock.Anything).Return(assert.AnError)
	m1 := NewServiceManager(Deps{Logger: logger})
	m1.RegisterService(sInitFail)
	assert.Error(t, m1.InitServices(context.Background()))

	// Start Fails
	sStartFail := &MockService{id: "start_fail"}
	sStartFail.On("Init", mock.Anything).Return(nil)
	sStartFail.On("Start", mock.Anything).Return(assert.AnError)
	m2 := NewServiceManager(Deps{Logger: logger})
	m2.RegisterService(sStartFail)
	m2.InitServices(context.Background())
	assert.Error(t, m2.StartServices(context.Background()))

	// Stop Fails (should log but not return error)
	sStopFail := &MockService{id: "stop_fail"}
	sStopFail.On("Init", mock.Anything).Return(nil)
	sStopFail.On("Start", mock.Anything).Return(nil)
	sStopFail.On("Stop", mock.Anything).Return(assert.AnError)
	m3 := NewServiceManager(Deps{Logger: logger})
	m3.RegisterService(sStopFail)
	m3.InitServices(context.Background())
	m3.StartServices(context.Background())
	err := m3.StopServices(context.Background())
	assert.NoError(t, err) // StopServices swallows errors
	sStopFail.AssertExpectations(t)
}

func TestServiceManager_BuildFromConfig_ConfigInitError(t *testing.T) {
	// Case: ConfigAware InitConfig fails
	logger := zap.NewNop()
	cfg := &config.Config{
		Services: map[string]interface{}{
			"bad_cfg_svc": map[string]interface{}{"enabled": true, "key": "val"},
		},
	}

	RegisterFactory("bad_cfg_svc", func(d Deps) (Service, error) {
		ms := &MockConfigAwareService{MockService{id: "bad_cfg_svc"}}
		ms.On("InitConfig", mock.Anything).Return(assert.AnError)
		return ms, nil
	})

	mgr := NewServiceManager(Deps{Logger: logger, Config: cfg})
	err := mgr.BuildFromConfig()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to init config")
}

func TestServiceManager_AutoWiring_Errors(t *testing.T) {
	mockRouter := "mockRouter"
	deps := Deps{Logger: zap.NewNop(), WebRouter: mockRouter}
	mgr := NewServiceManager(deps)

	// Web Register Fails
	webSvc := &MockWebRoutableService{MockService{id: "web_fail"}}
	webSvc.On("Init", mock.Anything).Return(nil)
	webSvc.On("RegisterRoutes", mockRouter).Return(assert.AnError)

	mgr.RegisterService(webSvc)
	err := mgr.InitServices(context.Background())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to register web routes")
}
