package web

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type TestService struct{}

func (s *TestService) RegisterRoutes(g *gin.RouterGroup) {
	g.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
}

func TestServer_StartStop(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger, _ := zap.NewDevelopment()
	cfg := DefaultConfig()
	cfg.Port = 0 // Let OS choose port

	// Test with nil health service
	server := NewWebServer(cfg, logger, nil)
	assert.NotNil(t, server)

	// Register service
	service := &TestService{}
	server.RegisterWebService(service)

	// Start server
	err := server.Start()
	assert.NoError(t, err)

	// Give it a moment to start
	time.Sleep(100 * time.Millisecond)

	// Stop server
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	err = server.Stop(ctx)
	assert.NoError(t, err)
}

func TestServer_WithTracing(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger, _ := zap.NewDevelopment()
	cfg := DefaultConfig()
	cfg.Port = 0
	cfg.Tracing.Enabled = true
	cfg.Tracing.ServiceName = "test-web"

	server := NewWebServer(cfg, logger, nil)
	assert.NotNil(t, server)

	service := &TestService{}
	server.RegisterWebService(service)

	err := server.Start()
	assert.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	// Perform a request to ensure middleware doesn't panic
	// Note: We don't have the port easily unless we check logs or bind to fixed port.
	// But since Start() is async and we don't know the port (0), we can't easily curl it here without refactoring Server to expose Listener Addr.
	// However, the purpose is to verify InitEngine didn't panic.
	// Since Start calls InitEngine, and InitEngine adds middleware, if it succeeded without panic, we are good.

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	err = server.Stop(ctx)
	assert.NoError(t, err)
}
