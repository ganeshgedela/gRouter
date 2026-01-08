package app

import (
	"context"
	"grouter/pkg/manager"
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

func (s *StopService) ID() string {
	return "stop"
}

// Name returns the service name "stop".
func (s *StopService) Name() string {
	return "stop"
}

func (s *StopService) Status() manager.HealthStatus {
	return manager.StatusRunning
}

func (s *StopService) LastError() error {
	return nil
}

// Dependencies returns the list of dependencies.
func (s *StopService) Dependencies() []string {
	return nil
}

func (s *StopService) Init(ctx context.Context) error {
	return nil
}

func (s *StopService) Start(ctx context.Context) error {
	return nil
}

func (s *StopService) Stop(ctx context.Context) error {
	return nil
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
