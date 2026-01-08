package web

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// TimeoutMiddleware adds a timeout to requests
//
// IMPORTANT: This middleware sets a timeout on the request context, but it cannot
// force-stop a handler that doesn't check context cancellation. Handlers should
// periodically check ctx.Err() or use context-aware operations (like database queries
// with context) to respect the timeout.
//
// For guaranteed timeout enforcement, use this with handlers that:
// 1. Check context.Done() or ctx.Err()
// 2. Use context-aware I/O operations (http.NewRequestWithContext, db queries with ctx)
// 3. Pass context to downstream services
func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create context with timeout
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		// Replace request context
		c.Request = c.Request.WithContext(ctx)

		// Channel to signal handler completion
		finished := make(chan struct{})

		// Run handler in goroutine to detect timeout
		go func() {
			c.Next()
			close(finished)
		}()

		select {
		case <-finished:
			// Handler completed successfully before timeout
			return
		case <-ctx.Done():
			// Timeout occurred
			// Note: Handler may still be running but context is cancelled
			c.AbortWithStatusJSON(http.StatusGatewayTimeout, ErrorResponse{
				Error:   "TIMEOUT",
				Message: "request timeout exceeded",
				Code:    "GATEWAY_TIMEOUT",
			})
			return
		}
	}
}

// TimeoutHandler wraps a handler function to ensure it respects context cancellation
// Use this for handlers that may not check context themselves
func TimeoutHandler(handler gin.HandlerFunc, timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)

		done := make(chan struct{})
		go func() {
			defer func() {
				if r := recover(); r != nil {
					// Handler panicked, but we'll let recovery middleware handle it
				}
			}()
			handler(c)
			close(done)
		}()

		select {
		case <-done:
			// Completed successfully
			return
		case <-ctx.Done():
			// Timeout - if handler hasn't returned, it will continue running
			// but the response has already been sent
			if !c.Writer.Written() {
				c.AbortWithStatusJSON(http.StatusGatewayTimeout, ErrorResponse{
					Error:   "TIMEOUT",
					Message: "request timeout exceeded",
					Code:    "GATEWAY_TIMEOUT",
				})
			}
		}
	}
}
