package web

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

// TracingMiddleware adds OpenTelemetry tracing
func TracingMiddleware(serviceName string) gin.HandlerFunc {
	return otelgin.Middleware(serviceName,
		otelgin.WithFilter(func(req *http.Request) bool {
			// Skip tracing for health checks and metrics
			path := req.URL.Path
			return !strings.HasPrefix(path, "/health") &&
				!strings.HasPrefix(path, "/metrics")
		}),
	)
}
