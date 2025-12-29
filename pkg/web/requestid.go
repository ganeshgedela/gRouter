package web

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	// HeaderXRequestID is the header name for request ID
	HeaderXRequestID = "X-Request-ID"
)

// RequestIDMiddleware adds a unique ID to every request
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := c.GetHeader(HeaderXRequestID)
		if rid == "" {
			rid = uuid.New().String()
		}

		// Set the header for the response
		c.Header(HeaderXRequestID, rid)

		// Set the ID in the context for other middleware/handlers to use
		c.Set("RequestID", rid)

		c.Next()
	}
}
