package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RBACMiddleware provides role-based access control
func RBACMiddleware(logger *zap.Logger, requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user roles from context (set by auth middleware)
		rolesInterface, exists := c.Get("user_roles")
		if !exists {
			if logger != nil {
				logger.Warn("RBAC: roles not found in context",
					zap.String("request_id", c.GetString("request_id")),
					zap.String("path", c.Request.URL.Path),
					zap.String("method", c.Request.Method),
				)
			}
			c.AbortWithStatusJSON(http.StatusForbidden, ErrorResponse{
				Error:   "NO_ROLES",
				Message: "user roles not found in context",
				Code:    "FORBIDDEN",
			})
			return
		}

		userRoles, ok := rolesInterface.([]string)
		if !ok {
			if logger != nil {
				logger.Warn("RBAC: invalid roles format",
					zap.String("request_id", c.GetString("request_id")),
					zap.Any("roles_value", rolesInterface),
				)
			}
			c.AbortWithStatusJSON(http.StatusForbidden, ErrorResponse{
				Error:   "INVALID_ROLES",
				Message: "invalid roles format",
				Code:    "FORBIDDEN",
			})
			return
		}

		// Check if user has any of the required roles
		hasRole := false
		for _, required := range requiredRoles {
			for _, userRole := range userRoles {
				if userRole == required {
					hasRole = true
					break
				}
			}
			if hasRole {
				break
			}
		}

		if !hasRole {
			if logger != nil {
				logger.Warn("RBAC: insufficient permissions",
					zap.String("request_id", c.GetString("request_id")),
					zap.String("user_id", c.GetString("user_id")),
					zap.Strings("user_roles", userRoles),
					zap.Strings("required_roles", requiredRoles),
					zap.String("path", c.Request.URL.Path),
					zap.String("method", c.Request.Method),
				)
			}
			c.AbortWithStatusJSON(http.StatusForbidden, ErrorResponse{
				Error:   "INSUFFICIENT_PERMISSIONS",
				Message: "you do not have the required role",
				Code:    "FORBIDDEN",
			})
			return
		}

		c.Next()
	}
}

// RequireRole is an alias for RBACMiddleware for single role
func RequireRole(logger *zap.Logger, role string) gin.HandlerFunc {
	return RBACMiddleware(logger, role)
}

// RequireAnyRole requires at least one of the specified roles
func RequireAnyRole(logger *zap.Logger, roles ...string) gin.HandlerFunc {
	return RBACMiddleware(logger, roles...)
}

// RequireAllRoles requires all of the specified roles
func RequireAllRoles(logger *zap.Logger, roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user roles from context
		rolesInterface, exists := c.Get("user_roles")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, ErrorResponse{
				Error:   "NO_ROLES",
				Message: "user roles not found in context",
				Code:    "FORBIDDEN",
			})
			return
		}

		userRoles, ok := rolesInterface.([]string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, ErrorResponse{
				Error:   "INVALID_ROLES",
				Message: "invalid roles format",
				Code:    "FORBIDDEN",
			})
			return
		}

		// Check if user has all required roles
		for _, required := range roles {
			found := false
			for _, userRole := range userRoles {
				if userRole == required {
					found = true
					break
				}
			}
			if !found {
				c.AbortWithStatusJSON(http.StatusForbidden, ErrorResponse{
					Error:   "INSUFFICIENT_PERMISSIONS",
					Message: "you do not have all required roles",
					Code:    "FORBIDDEN",
				})
				return
			}
		}

		c.Next()
	}
}
