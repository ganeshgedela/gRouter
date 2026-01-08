package web

import (
	"grouter/pkg/manager"

	"github.com/gin-gonic/gin"
)

// Server interface for HTTP server
type Server interface {
	manager.Service // Implements manager.Service interface

	// Additional web-specific methods
	Engine() *gin.Engine
	RegisterService(service WebService)
	Use(middleware ...gin.HandlerFunc)
	RegisterRoutes(path string, handler gin.HandlerFunc, methods ...string)
}

// WebService interface for HTTP services
type WebService interface {
	RegisterRoutes(router gin.IRouter)
	ServiceName() string
}

// Middleware represents a Gin middleware function
type Middleware func(gin.HandlerFunc) gin.HandlerFunc
