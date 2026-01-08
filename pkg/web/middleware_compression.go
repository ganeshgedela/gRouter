package web

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

// CompressionMiddleware adds gzip compression to responses
func CompressionMiddleware() gin.HandlerFunc {
	return gzip.Gzip(gzip.DefaultCompression, gzip.WithExcludedPaths([]string{
		"/metrics",     // Don't compress metrics
		"/health/live", // Don't compress health checks
		"/health/ready",
	}))
}
