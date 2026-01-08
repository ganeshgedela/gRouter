package web

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	httpRequestSize = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_size_bytes",
			Help:    "Size of HTTP requests in bytes",
			Buckets: prometheus.ExponentialBuckets(100, 10, 8), // 100B to ~10MB
		},
		[]string{"method", "path"},
	)

	httpResponseSize = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_size_bytes",
			Help:    "Size of HTTP responses in bytes",
			Buckets: prometheus.ExponentialBuckets(100, 10, 8),
		},
		[]string{"method", "path", "status"},
	)
)

// MetricsMiddleware records HTTP metrics
func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Get full path template (e.g., /users/:id instead of /users/123)
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		// Record request size
		if c.Request.ContentLength > 0 {
			httpRequestSize.WithLabelValues(
				c.Request.Method,
				path,
			).Observe(float64(c.Request.ContentLength))
		}

		c.Next()

		// Record metrics after request
		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())

		httpRequestsTotal.WithLabelValues(
			c.Request.Method,
			path,
			status,
		).Inc()

		httpRequestDuration.WithLabelValues(
			c.Request.Method,
			path,
		).Observe(duration)

		httpResponseSize.WithLabelValues(
			c.Request.Method,
			path,
			status,
		).Observe(float64(c.Writer.Size()))
	}
}

// ResetMetrics resets all metrics (useful for testing)
func ResetMetrics() {
	httpRequestsTotal.Reset()
	httpRequestDuration.Reset()
	httpRequestSize.Reset()
	httpResponseSize.Reset()
}
