package webdemo

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Service implements manager.Service and web.WebService interfaces
type Service struct{}

// NewService creates a new WebDemoService
func NewService() *Service {
	return &Service{}
}

// Name returns the service name
func (s *Service) Name() string {
	return "webdemo"
}

// Ready checks if the service is ready
func (s *Service) Ready(ctx context.Context) error {
	return nil
}

// Start starts the service
func (s *Service) Start(ctx context.Context) error {
	return nil
}

// Stop stops the service
func (s *Service) Stop(ctx context.Context) error {
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
