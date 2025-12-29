package web

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"grouter/pkg/health"
)

type IntegrationTestService struct{}

func (s *IntegrationTestService) RegisterRoutes(g *gin.RouterGroup) {
	g.GET("/integration", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
}

func TestServerIntegration(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	logger, _ := zap.NewDevelopment()
	cfg := DefaultConfig()
	cfg.Port = 0 // Random port
	cfg.Metrics.Enabled = true
	cfg.RateLimit.Enabled = true
	cfg.RateLimit.RequestsPerSecond = 100
	cfg.RateLimit.Burst = 200
	cfg.Swagger.Enabled = false // Disable swagger for test to avoid dependency issues

	healthSvc := health.NewHealthService()
	server := NewWebServer(cfg, logger, healthSvc)
	server.RegisterWebService(&IntegrationTestService{})

	// Start server in goroutine
	go func() {
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			t.Errorf("Failed to start server: %v", err)
		}
	}()

	// Wait for server to start
	time.Sleep(200 * time.Millisecond)
	defer server.Stop(context.Background())

	// Get the actual port (since we used 0)
	// In a real scenario we'd need a way to get the listener address,
	// but for this test we can assume it bound to the port in config if it was non-zero.
	// Since we used 0, we need to find the port.
	// However, Server struct doesn't expose the listener.
	// For integration test, let's use a fixed port to simplify, or we'd need to modify Server to expose Addr.
	// Let's use a fixed port that is likely free.

	// RESTART with fixed port
	server.Stop(context.Background())
	cfg.Port = 18085
	server = NewWebServer(cfg, logger, healthSvc)
	server.RegisterWebService(&IntegrationTestService{})
	go func() {
		server.Start()
	}()
	time.Sleep(200 * time.Millisecond)
	defer server.Stop(context.Background())

	baseURL := "http://localhost:18085"

	// Test 1: Service Endpoint
	resp, err := http.Get(baseURL + "/integration")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotEmpty(t, resp.Header.Get("X-Request-ID"))
	resp.Body.Close()

	// Test 2: Health Check
	resp, err = http.Get(baseURL + "/health/live")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()

	// Test 3: Metrics
	resp, err = http.Get(baseURL + "/metrics")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()

	// Test 4: 404
	resp, err = http.Get(baseURL + "/notfound")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	resp.Body.Close()
}
