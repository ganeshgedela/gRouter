package web

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

func TestMetricsMiddleware(t *testing.T) {
	// Reset metrics before test
	ResetMetrics()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(MetricsMiddleware())

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	// Make a request
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	// Verify metrics were recorded
	metrics, err := prometheus.DefaultGatherer.Gather()
	if err != nil {
		t.Fatalf("failed to gather metrics: %v", err)
	}

	// Check if our metrics exist
	foundRequestsTotal := false
	foundRequestDuration := false

	for _, mf := range metrics {
		if mf.GetName() == "http_requests_total" {
			foundRequestsTotal = true
		}
		if mf.GetName() == "http_request_duration_seconds" {
			foundRequestDuration = true
		}
	}

	if !foundRequestsTotal {
		t.Error("http_requests_total metric not found")
	}
	if !foundRequestDuration {
		t.Error("http_request_duration_seconds metric not found")
	}
}

func TestTimeoutMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		timeout        time.Duration
		handlerDelay   time.Duration
		expectedStatus int
	}{
		{
			name:           "request completes before timeout",
			timeout:        100 * time.Millisecond,
			handlerDelay:   10 * time.Millisecond,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "request times out",
			timeout:        50 * time.Millisecond,
			handlerDelay:   200 * time.Millisecond,
			expectedStatus: http.StatusGatewayTimeout,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.Use(TimeoutMiddleware(tt.timeout))

			router.GET("/test", func(c *gin.Context) {
				time.Sleep(tt.handlerDelay)
				c.JSON(http.StatusOK, gin.H{"message": "ok"})
			})

			req := httptest.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestRateLimitMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Allow 2 requests per second with burst of 2
	router.Use(RateLimitMiddleware(2, 2))

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	// First 2 requests should succeed
	for i := 0; i < 2; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("request %d: expected status 200, got %d", i+1, w.Code)
		}
	}

	// 3rd request should be rate limited
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Errorf("expected status 429, got %d", w.Code)
	}
}

func TestCompressionMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(CompressionMiddleware())

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "this is a long response that should be compressed",
			"data":    make([]int, 100),
		})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	// Check if response was compressed
	if w.Header().Get("Content-Encoding") != "gzip" {
		t.Error("expected gzip encoding, but not found")
	}
}

func TestCompressionMiddleware_ExcludedPaths(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(CompressionMiddleware())

	router.GET("/metrics", func(c *gin.Context) {
		c.String(http.StatusOK, "metrics data")
	})

	req := httptest.NewRequest("GET", "/metrics", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Metrics should NOT be compressed
	if w.Header().Get("Content-Encoding") == "gzip" {
		t.Error("metrics endpoint should not be compressed")
	}
}

func getMetricValue(name string, labels map[string]string) (float64, error) {
	metrics, err := prometheus.DefaultGatherer.Gather()
	if err != nil {
		return 0, err
	}

	for _, mf := range metrics {
		if mf.GetName() == name {
			for _, m := range mf.GetMetric() {
				if matchesLabels(m, labels) {
					if m.Counter != nil {
						return m.Counter.GetValue(), nil
					}
					if m.Gauge != nil {
						return m.Gauge.GetValue(), nil
					}
				}
			}
		}
	}

	return 0, nil
}

func matchesLabels(m *dto.Metric, labels map[string]string) bool {
	for key, val := range labels {
		found := false
		for _, lp := range m.Label {
			if lp.GetName() == key && lp.GetValue() == val {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
