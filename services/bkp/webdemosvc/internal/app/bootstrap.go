package app

import (
	"context"
	"grouter/pkg/manager"
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

func (s *BootstrapService) ID() string {
	return "start"
}

// Name returns the service name "start".
func (s *BootstrapService) Name() string {
	return "start"
}

func (s *BootstrapService) Status() manager.HealthStatus {
	return manager.StatusRunning
}

func (s *BootstrapService) LastError() error {
	return nil
}

// Dependencies returns the list of dependencies.
func (s *BootstrapService) Dependencies() []string {
	return nil
}

func (s *BootstrapService) Init(ctx context.Context) error {
	return nil
}

func (s *BootstrapService) Start(ctx context.Context) error {
	return nil
}

func (s *BootstrapService) Stop(ctx context.Context) error {
	return nil
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
