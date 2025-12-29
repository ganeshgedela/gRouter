package app

import (
	"grouter/pkg/web"
	"net/http"

	"github.com/gin-gonic/gin"
)

// StopService waits for a start signal.
type StopService struct {
	trigger   chan struct{}
	webServer *web.Server
}

// NewStopService creates a new StopService.
func NewStopService(trigger chan struct{}, webServer *web.Server) *StopService {
	return &StopService{
		trigger:   trigger,
		webServer: webServer,
	}
}

// Name returns the service name "stop".
func (s *StopService) Name() string {
	return "stop"
}

// RegisterRoutes registers the handlers for this service
func (s *StopService) RegisterRoutes(g *gin.RouterGroup) {
	g.GET("/stop", s.StopHandler)
}

func (s *StopService) StopHandler(c *gin.Context) {
	select {
	case s.trigger <- struct{}{}:
		c.JSON(http.StatusOK, gin.H{"status": "stopping"})
	default:
		c.JSON(http.StatusOK, gin.H{"status": "already stopping"})
	}
}
