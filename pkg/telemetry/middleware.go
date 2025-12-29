package telemetry

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

// Middleware returns a Gin middleware that handles both Tracing and Metrics
func Middleware(serviceName string) gin.HandlerFunc {
	tracer := otel.Tracer("grouter/pkg/telemetry")

	return func(c *gin.Context) {
		start := time.Now()
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		// 1. Tracing: Extract context and start span
		ctx := otel.GetTextMapPropagator().Extract(c.Request.Context(), propagation.HeaderCarrier(c.Request.Header))
		ctx, span := tracer.Start(ctx, path, trace.WithAttributes(
			semconv.HTTPMethod(c.Request.Method),
			semconv.HTTPURL(c.Request.URL.String()),
			semconv.HTTPRoute(path),
		))
		defer span.End()

		// Inject trace context into Gin context
		c.Request = c.Request.WithContext(ctx)

		// 2. Metrics: Increment active requests
		httpActiveRequests.WithLabelValues(serviceName).Inc()
		defer httpActiveRequests.WithLabelValues(serviceName).Dec()

		// Process request
		c.Next()

		// 3. Metrics: Record duration and count
		status := strconv.Itoa(c.Writer.Status())
		duration := time.Since(start).Seconds()

		httpRequestsTotal.WithLabelValues(serviceName, c.Request.Method, path, status).Inc()
		httpRequestDuration.WithLabelValues(serviceName, c.Request.Method, path, status).Observe(duration)

		// 4. Tracing: Update span with status
		span.SetAttributes(semconv.HTTPStatusCode(c.Writer.Status()))
		if c.Writer.Status() >= 500 {
			span.RecordError(nil) // Error detail could be extracted from c.Errors
		}
	}
}
