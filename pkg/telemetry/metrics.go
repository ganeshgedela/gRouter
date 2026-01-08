package telemetry

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests processed",
		},
		[]string{"service", "method", "path", "status"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"service", "method", "path", "status"},
	)

	httpActiveRequests = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "http_active_requests",
			Help: "Number of currently active HTTP requests",
		},
		[]string{"service"},
	)
)

// InitMetrics registers the standard metrics with Prometheus
func InitMetrics() {
	// Register metrics with the global prometheus registry
	// We use Register instead of MustRegister to avoid panics if called multiple times in tests
	_ = prometheus.Register(httpRequestsTotal)
	_ = prometheus.Register(httpRequestDuration)
	_ = prometheus.Register(httpActiveRequests)
}

// PrometheusHandler returns a Gin handler for the metrics endpoint
func PrometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
