package app

import (
	"net/http"

	"grouter/pkg/health"

	"github.com/gin-gonic/gin"
)

// HealthService implements the Service interface for Web health checks
// It wraps the shared *health.HealthService and exposes HTTP handlers
type HealthService struct {
	svc     *health.HealthService
	appName string
}

// NewHealthService creates a new HealthService
func NewHealthService(svc *health.HealthService, appName string) *HealthService {
	return &HealthService{
		svc:     svc,
		appName: appName,
	}
}

// Name returns the service name
func (s *HealthService) Name() string {
	return "health"
}

// RegisterRoutes registers the handlers for this service
func (s *HealthService) RegisterRoutes(g *gin.RouterGroup) {
	// The manager already registers /health/live and /health/ready globally on the root.
	// We can add them here if we wanted them under a /health prefix group specifically for this service module,
	// or we can just leave it empty if the global ones are sufficient.
	// For now, let's explicitly add them to be sure, although it might duplicate if the group is root.
	// Actually, manager uses: server.engine.GET("/health/live", ...)
	// If we register here, it depends on what 'g' is.
	// Usually 'g' is the root group passed from manager.
	// Let's not duplicate routes to avoid conflicts, but we can add an "info" route.

	// Example: adding a specific check for this service
	// s.svc.AddLivenessCheck(s.appName, func() error { return nil })
	g.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
			"app":    s.appName,
		})
	})
}
