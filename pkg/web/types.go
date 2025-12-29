package web

import "github.com/gin-gonic/gin"

// Service defines a component that exposes HTTP endpoints.
// Services implementing this interface can be registered with the Web Server.
type WebService interface {
	// RegisterRoutes registers the service's routes on the provided RouterGroup.
	// The router group will be scoped to the service's base path if configured,
	// or the root if not.
	RegisterRoutes(router *gin.RouterGroup)
}
