package web

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func noopHandler(c *gin.Context) {
	c.Status(http.StatusOK)
}

func setupBenchmarkServer(enableTracing bool) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	logger := zap.NewNop()

	// Replicate NewWebServer logic roughly or just use it if we can access engine.
	// We can't access engine. So we replicate middleware stack.

	engine := gin.New()
	engine.Use(gin.Recovery())

	// Add our middleware
	engine.Use(RequestIDMiddleware())
	engine.Use(LoggerMiddleware(logger))

	// Tracing
	if enableTracing {
		// Mock config or just assume it works.
		// For benchmark we skip actual otel init (complex), but we want to measure overhead if possible.
		// Without InitTracer, otelgin might be no-op or default.
		// We will test core middleware overhead.
	}

	engine.GET("/bench", noopHandler)
	return engine
}

func BenchmarkMiddlewareLoop(b *testing.B) {
	engine := setupBenchmarkServer(false)
	req := httptest.NewRequest("GET", "/bench", nil)
	w := httptest.NewRecorder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Create new recorder per iteration to avoid buffer growth?
		// Gin's ServeHTTP writes to writer.
		// Reusing recorder is faster but we should reset it.
		// For pure middleware bench:
		engine.ServeHTTP(w, req)
		// Reset recorder
		w.Body.Reset()
	}
}

func BenchmarkMiddlewareLoop_WithTracingDisabled(b *testing.B) {
	// This measures our standard middleware stack (Logger, RequestID)
	// Tracing is disabled in config
	engine := setupBenchmarkServer(false)
	req := httptest.NewRequest("GET", "/bench", nil)
	w := httptest.NewRecorder()
	w.Body.Reset()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		engine.ServeHTTP(w, req)
		w.Body.Reset()
	}
}
