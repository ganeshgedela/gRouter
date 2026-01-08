package web

import (
	"context"
	"net/http"
	"testing"
	"time"

	"grouter/pkg/config"
	"grouter/pkg/manager"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func TestNewServer(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	tests := []struct {
		name    string
		id      string
		sname   string
		cfg     config.WebConfig
		logger  *zap.Logger
		wantErr bool
	}{
		{
			name:    "valid server",
			id:      "web-test",
			sname:   "Test Web Server",
			cfg:     validTestConfig(),
			logger:  logger,
			wantErr: false,
		},
		{
			name:    "missing logger",
			id:      "web-test",
			sname:   "Test",
			cfg:     validTestConfig(),
			logger:  nil,
			wantErr: true,
		},
		{
			name:    "missing ID",
			id:      "",
			sname:   "Test",
			cfg:     validTestConfig(),
			logger:  logger,
			wantErr: true,
		},
		{
			name:    "invalid port",
			id:      "web-test",
			sname:   "Test",
			cfg:     invalidPortConfig(),
			logger:  logger,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv, err := NewServer(tt.id, tt.sname, tt.cfg, tt.logger)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewServer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && srv == nil {
				t.Error("NewServer() returned nil server")
			}
		})
	}
}

func TestServer_ManagerInterface(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	cfg := validTestConfig()

	srv, err := NewServer("web-test", "Test Server", cfg, logger)
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	// Verify it implements manager.Service
	var _ manager.Service = srv

	// Test ID
	if srv.ID() != "web-test" {
		t.Errorf("ID() = %v, want %v", srv.ID(), "web-test")
	}

	// Test Name
	if srv.Name() != "Test Server" {
		t.Errorf("Name() = %v, want %v", srv.Name(), "Test Server")
	}

	// Test Status (should be Created initially)
	if srv.Status() != manager.StatusCreated {
		t.Errorf("Status() = %v, want %v", srv.Status(), manager.StatusCreated)
	}

	// Test LastError (should be nil initially)
	if srv.LastError() != nil {
		t.Errorf("LastError() = %v, want nil", srv.LastError())
	}

	// Test Dependencies
	deps := srv.Dependencies()
	if len(deps) != 0 {
		t.Errorf("Dependencies() = %v, want empty slice", deps)
	}
}

func TestServer_Lifecycle(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	cfg := testConfigWithPort(18080)

	srv, err := NewServer("web-lifecycle", "Lifecycle Test", cfg, logger)
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	ctx := context.Background()

	// Test Init
	if err := srv.Init(ctx); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	if srv.Status() != manager.StatusInitialized {
		t.Errorf("Status after Init = %v, want %v", srv.Status(), manager.StatusInitialized)
	}

	// Test Start
	if err := srv.Start(ctx); err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	if srv.Status() != manager.StatusRunning {
		t.Errorf("Status after Start = %v, want %v", srv.Status(), manager.StatusRunning)
	}

	// Give server time to fully start
	time.Sleep(100 * time.Millisecond)

	// Test health endpoints
	resp, err := http.Get("http://localhost:18080/health/live")
	if err != nil {
		t.Errorf("health check failed: %v", err)
	} else {
		resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("health check status = %v, want %v", resp.StatusCode, http.StatusOK)
		}
	}

	// Test Stop
	if err := srv.Stop(ctx); err != nil {
		t.Fatalf("Stop() error = %v", err)
	}

	if srv.Status() != manager.StatusStopped {
		t.Errorf("Status after Stop = %v, want %v", srv.Status(), manager.StatusStopped)
	}
}

func TestServer_DoubleStart(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	cfg := testConfigWithPort(18081)

	srv, err := NewServer("web-double-start", "Double Start Test", cfg, logger)
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	ctx := context.Background()

	if err := srv.Init(ctx); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	if err := srv.Start(ctx); err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	defer srv.Stop(ctx)

	// Try to start again
	err = srv.Start(ctx)
	if err == nil {
		t.Error("Start() should fail when server already running")
	}
}

func TestServer_StopWithoutStart(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	cfg := validTestConfig()

	srv, err := NewServer("web-stop", "Stop Test", cfg, logger)
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	ctx := context.Background()

	// Stop without start should not error
	if err := srv.Stop(ctx); err != nil {
		t.Errorf("Stop() without Start should not error, got: %v", err)
	}
}

func TestServer_RegisterService(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	cfg := validTestConfig()

	srv, err := NewServer("web-register", "Register Test", cfg, logger)
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	// Create mock web service
	mockSvc := &mockWebService{name: "mock-service"}
	srv.RegisterService(mockSvc)

	if !mockSvc.registered {
		t.Error("RegisterService() did not call RegisterRoutes")
	}
}

func TestServer_RegisterRoutes(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	cfg := testConfigWithPort(18082)

	srv, err := NewServer("web-routes", "Routes Test", cfg, logger)
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	handler := func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "test"})
	}

	srv.RegisterRoutes("/test", handler, "GET", "POST")

	ctx := context.Background()
	if err := srv.Init(ctx); err != nil {
		t.Fatalf("Init() error = %v", err)
	}
	if err := srv.Start(ctx); err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	defer srv.Stop(ctx)

	time.Sleep(100 * time.Millisecond)

	// Test GET
	resp, err := http.Get("http://localhost:18082/test")
	if err != nil {
		t.Errorf("GET request failed: %v", err)
	} else {
		resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("GET status = %v, want %v", resp.StatusCode, http.StatusOK)
		}
	}
}

// Helper functions

func validTestConfig() config.WebConfig {
	return config.WebConfig{
		Port:            18080,
		Mode:            "test",
		ReadTimeout:     5 * time.Second,
		WriteTimeout:    5 * time.Second,
		ShutdownTimeout: 2 * time.Second,
		Logging: config.LoggingConfig{
			Enabled: true,
		},
	}
}

func testConfigWithPort(port int) config.WebConfig {
	cfg := validTestConfig()
	cfg.Port = port
	return cfg
}

func invalidPortConfig() config.WebConfig {
	cfg := validTestConfig()
	cfg.Port = 99999 // Invalid port
	return cfg
}

// Mock WebService
type mockWebService struct {
	name       string
	registered bool
}

func (m *mockWebService) RegisterRoutes(router gin.IRouter) {
	m.registered = true
}

func (m *mockWebService) ServiceName() string {
	return m.name
}
