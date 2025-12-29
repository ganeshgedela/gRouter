package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// BootstrapService waits for a start signal.
type BootstrapService struct {
	trigger chan struct{}
}

// NewBootstrapService creates a new BootstrapService.
func NewBootstrapService(trigger chan struct{}) *BootstrapService {
	return &BootstrapService{
		trigger: trigger,
	}
}

// Name returns the service name "start".
func (s *BootstrapService) Name() string {
	return "start"
}

// RegisterRoutes registers the handlers for this service
func (s *BootstrapService) RegisterRoutes(g *gin.RouterGroup) {
	g.GET("/start", s.StartHandler)
}

func (s *BootstrapService) StartHandler(c *gin.Context) {
	select {
	case s.trigger <- struct{}{}:
		c.JSON(http.StatusOK, gin.H{"status": "starting"})
	default:
		c.JSON(http.StatusOK, gin.H{"status": "already started"})
	}
}
