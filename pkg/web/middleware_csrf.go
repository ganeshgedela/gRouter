package web

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CSRFMiddleware provides CSRF protection
func CSRFMiddleware(logger *zap.Logger, tokenLength int) gin.HandlerFunc {
	if tokenLength == 0 {
		tokenLength = 32 // default
	}

	return func(c *gin.Context) {
		// Skip for safe methods
		if c.Request.Method == "GET" || c.Request.Method == "HEAD" ||
			c.Request.Method == "OPTIONS" || c.Request.Method == "TRACE" {
			// Generate and set token for safe methods
			token := generateToken(tokenLength)
			c.Set("csrf_token", token)
			c.Header("X-CSRF-Token", token)
			c.Next()
			return
		}

		// For unsafe methods, verify token
		requestToken := c.GetHeader("X-CSRF-Token")
		if requestToken == "" {
			// Try form value
			requestToken = c.PostForm("csrf_token")
		}

		// Get session token (in production, this would come from session store)
		sessionToken, exists := c.Get("csrf_token")
		if !exists || requestToken == "" || requestToken != sessionToken {
			if logger != nil {
				logger.Warn("CSRF validation failed",
					zap.String("request_id", c.GetString("request_id")),
					zap.String("method", c.Request.Method),
					zap.String("path", c.Request.URL.Path),
					zap.Bool("token_exists", exists),
					zap.Bool("request_token_empty", requestToken == ""),
					zap.String("client_ip", c.ClientIP()),
				)
			}
			c.AbortWithStatusJSON(http.StatusForbidden, ErrorResponse{
				Error:   "CSRF_VALIDATION_FAILED",
				Message: "CSRF token validation failed",
				Code:    "FORBIDDEN",
			})
			return
		}

		c.Next()
	}
}

// generateToken generates a random token
func generateToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}
