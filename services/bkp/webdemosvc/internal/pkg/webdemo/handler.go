package webdemo

import (
	"context"
	"net/http"

	"grouter/pkg/manager"

	"github.com/gin-gonic/gin"
)

func init() {
	manager.RegisterFactory("webdemo", func(deps manager.Deps) (manager.Service, error) {
		return NewService(deps), nil
	})
}

// Service implements manager.Service and manager.WebRoutable
type Service struct {
	status    manager.HealthStatus
	lastError error
}

// NewService creates a new WebDemoService
func NewService(deps manager.Deps) *Service {
	return &Service{
		status: manager.StatusCreated,
	}
}

// Name returns the service name
func (s *Service) Name() string {
	return "webdemo"
}

func (s *Service) ID() string {
	return "webdemo"
}

// Dependencies returns the list of dependencies.
func (s *Service) Dependencies() []string {
	return nil
}

// Ready checks if the service is ready
func (s *Service) Ready(ctx context.Context) error {
	return nil
}

func (s *Service) Status() manager.HealthStatus {
	return s.status
}

func (s *Service) LastError() error {
	return s.lastError
}

// Init initializes the service
func (s *Service) Init(ctx context.Context) error {
	s.status = manager.StatusInitialized
	return nil
}

// Start starts the service
func (s *Service) Start(ctx context.Context) error {
	s.status = manager.StatusRunning
	return nil
}

// Stop stops the service
func (s *Service) Stop(ctx context.Context) error {
	s.status = manager.StatusStopped
	return nil
}

// RegisterRoutes registers the handlers for this service
func (s *Service) RegisterRoutes(g *gin.RouterGroup) {
	// Debug print (using stdlib as logger isn't passed here easily, or assumes zap capture)
	g.GET("/hello", s.HelloHandler)
	g.GET("/echo", s.EchoHandler)
}

// HelloHandler says hello
func (s *Service) HelloHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello from WebDemoSvc!",
	})
}

// EchoHandler echoes the query param
func (s *Service) EchoHandler(c *gin.Context) {
	msg := c.DefaultQuery("msg", "nothing")
	c.JSON(http.StatusOK, gin.H{
		"echo": msg,
	})
}
