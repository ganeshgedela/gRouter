package telemetry

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"grouter/pkg/config"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

func TestTelemetry(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	cfg := config.Config{
		Tracing: config.TracingConfig{
			Enabled:     false, // Disable exporter for test
			ServiceName: "test-service",
		},
		Metrics: config.MetricsConfig{
			Enabled: true,
		},
	}

	// 1. Test Init
	shutdown, err := Init(cfg)
	assert.NoError(t, err)
	defer shutdown(context.Background())

	// 2. Setup Router with Middleware
	r := gin.New()
	r.Use(Middleware("test-service"))
	r.GET("/ping", func(c *gin.Context) {
		time.Sleep(10 * time.Millisecond) // Simulate work
		c.String(200, "pong")
	})

	// 3. Test Request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "pong", w.Body.String())

	// 4. Verify Metrics
	// We check if the metric families are registered and have values
	mfs, err := prometheus.DefaultGatherer.Gather()
	assert.NoError(t, err)

	var requestCountFound, durationFound bool
	for _, mf := range mfs {
		if mf.GetName() == "http_requests_total" {
			requestCountFound = true
			for _, m := range mf.GetMetric() {
				// Verify labels
				found := false
				for _, label := range m.GetLabel() {
					if label.GetName() == "service" && label.GetValue() == "test-service" {
						found = true
					}
				}
				if found {
					assert.Equal(t, float64(1), m.GetCounter().GetValue())
				}
			}
		}
		if mf.GetName() == "http_request_duration_seconds" {
			durationFound = true
			for _, m := range mf.GetMetric() {
				assert.Greater(t, m.GetHistogram().GetSampleCount(), uint64(0))
			}
		}
	}

	assert.True(t, requestCountFound, "http_requests_total should be recorded")
	assert.True(t, durationFound, "http_request_duration_seconds should be recorded")

	// 5. Test Prometheus Handler
	reqMetrics, _ := http.NewRequest("GET", "/metrics", nil)
	wMetrics := httptest.NewRecorder()

	// Create a separate handler for metrics to avoid interference with the router used above
	metricsHandler := PrometheusHandler()
	metricsContext, _ := gin.CreateTestContext(wMetrics)
	metricsContext.Request = reqMetrics

	metricsHandler(metricsContext)

	assert.Equal(t, 200, wMetrics.Code)
	assert.Contains(t, wMetrics.Body.String(), "http_requests_total")
}
