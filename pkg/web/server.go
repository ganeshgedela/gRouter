package web

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"runtime/debug"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"grouter/pkg/config"
	"grouter/pkg/manager"
)

// server implements the Server interface and manager.Service
type server struct {
	id     string
	name   string
	engine *gin.Engine
	srv    *http.Server
	cfg    config.WebConfig
	logger *zap.Logger

	// State management
	mu       sync.RWMutex
	status   manager.HealthStatus
	lastErr  error
	listener net.Listener

	// Dependencies
	deps []string
}

// NewServer creates a new HTTP server compatible with manager
func NewServer(id, name string, cfg config.WebConfig, logger *zap.Logger) (Server, error) {
	if logger == nil {
		return nil, fmt.Errorf("logger is required")
	}

	if id == "" {
		return nil, fmt.Errorf("service ID is required")
	}

	// Validate configuration
	if err := ValidateConfig(cfg); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	// Set Gin mode
	switch cfg.Mode {
	case "release", "production":
		gin.SetMode(gin.ReleaseMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}

	// Create Gin engine
	engine := gin.New()

	s := &server{
		id:     id,
		name:   name,
		engine: engine,
		cfg:    cfg,
		logger: logger,
		status: manager.StatusCreated,
		deps:   []string{}, // No dependencies by default
	}

	// Setup core middleware
	s.setupCoreMiddleware()

	return s, nil
}

// setupCoreMiddleware sets up essential middleware
func (s *server) setupCoreMiddleware() {
	// Version header (first, so it's on all responses)
	s.engine.Use(func(c *gin.Context) {
		c.Header("X-Server-Version", "1.0.0")
		c.Next()
	})

	// Default security headers (always on, even if Security config disabled)
	s.engine.Use(func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Next()
	})

	// Request body size limit (to prevent DoS)
	maxBodySize := int64(10 << 20) // 10MB default
	if s.cfg.MaxRequestBodySize > 0 {
		maxBodySize = s.cfg.MaxRequestBodySize
	}
	s.engine.Use(func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBodySize)
		c.Next()
	})

	// Recovery middleware (after body limit)
	s.engine.Use(RecoveryMiddleware(s.logger))

	// Request ID
	s.engine.Use(RequestIDMiddleware())

	// Logging middleware
	if s.cfg.Logging.Enabled {
		s.engine.Use(LoggingMiddleware(s.logger))
	}
}

// --- Implement manager.Service interface ---

// ID returns the unique service identifier
func (s *server) ID() string {
	return s.id
}

// Name returns the service name
func (s *server) Name() string {
	return s.name
}

// Init initializes the web server
func (s *server) Init(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.status != manager.StatusCreated {
		return fmt.Errorf("service already initialized")
	}

	s.logger.Info("initializing web server",
		zap.String("id", s.id),
		zap.String("name", s.name),
		zap.Int("port", s.cfg.Port),
	)

	// Register health check endpoints
	s.engine.GET("/health/live", s.livenessHandler)
	s.engine.GET("/health/ready", s.readinessHandler)

	// Register metrics if enabled
	if s.cfg.Metrics.Enabled {
		metricsPath := s.cfg.Metrics.Path
		if metricsPath == "" {
			metricsPath = "/metrics"
		}
		s.engine.GET(metricsPath, gin.WrapH(promhttp.Handler()))
		s.logger.Info("metrics enabled", zap.String("path", metricsPath))
	}

	s.status = manager.StatusInitialized
	s.lastErr = nil

	s.logger.Info("web server initialized successfully")
	return nil
}

// Start starts the HTTP server
func (s *server) Start(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.status == manager.StatusRunning {
		return fmt.Errorf("server already running")
	}

	if s.status != manager.StatusInitialized {
		return fmt.Errorf("server not initialized, current status: %s", s.status)
	}

	s.status = manager.StatusStarting

	addr := fmt.Sprintf(":%d", s.cfg.Port)
	s.srv = &http.Server{
		Addr:         addr,
		Handler:      s.engine,
		ReadTimeout:  s.cfg.ReadTimeout,
		WriteTimeout: s.cfg.WriteTimeout,
		IdleTimeout:  90 * time.Second,
	}

	s.logger.Info("starting web server",
		zap.String("id", s.id),
		zap.String("addr", addr),
		zap.Bool("tls", s.cfg.TLS.Enabled),
	)

	// Start server in goroutine
	errChan := make(chan error, 1)
	go func() {
		var err error
		if s.cfg.TLS.Enabled {
			s.logger.Info("starting TLS server",
				zap.String("cert", s.cfg.TLS.CertFile),
				zap.String("key", s.cfg.TLS.KeyFile),
			)
			err = s.srv.ListenAndServeTLS(s.cfg.TLS.CertFile, s.cfg.TLS.KeyFile)
		} else {
			err = s.srv.ListenAndServe()
		}

		if err != nil && err != http.ErrServerClosed {
			s.logger.Error("server error", zap.Error(err))
			s.mu.Lock()
			s.status = manager.StatusFailed
			s.lastErr = err
			s.mu.Unlock()
			errChan <- err
		}
	}()

	// Wait a bit to check for immediate startup errors
	select {
	case err := <-errChan:
		s.status = manager.StatusFailed
		s.lastErr = err
		return fmt.Errorf("failed to start server: %w", err)
	case <-time.After(100 * time.Millisecond):
		s.status = manager.StatusRunning
		s.lastErr = nil
		s.logger.Info("web server started successfully",
			zap.String("id", s.id),
			zap.String("addr", addr),
		)
		return nil
	}
}

// Stop gracefully stops the HTTP server
func (s *server) Stop(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.status != manager.StatusRunning {
		s.logger.Warn("server not running, skipping stop",
			zap.String("status", string(s.status)),
		)
		return nil
	}

	s.status = manager.StatusStopping
	s.logger.Info("stopping web server", zap.String("id", s.id))

	// Create shutdown context with timeout
	shutdownCtx, cancel := context.WithTimeout(ctx, s.cfg.ShutdownTimeout)
	defer cancel()

	if err := s.srv.Shutdown(shutdownCtx); err != nil {
		s.logger.Error("error during shutdown, forcing close", zap.Error(err))
		closeErr := s.srv.Close()
		s.status = manager.StatusStopped
		s.lastErr = err
		return fmt.Errorf("shutdown error: %w, close error: %v", err, closeErr)
	}

	s.status = manager.StatusStopped
	s.lastErr = nil

	s.logger.Info("web server stopped successfully")
	return nil
}

// Status returns the current service status
func (s *server) Status() manager.HealthStatus {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.status
}

// LastError returns the last error encountered
func (s *server) LastError() error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.lastErr
}

// Dependencies returns service dependencies
func (s *server) Dependencies() []string {
	return s.deps
}

// SetDependencies sets service dependencies
func (s *server) SetDependencies(deps []string) {
	s.deps = deps
}

// --- Additional Server interface methods ---

// Engine returns the underlying Gin engine
func (s *server) Engine() *gin.Engine {
	return s.engine
}

// RegisterService registers a WebService
func (s *server) RegisterService(service WebService) {
	s.logger.Info("registering web service",
		zap.String("server", s.id),
		zap.String("service", service.ServiceName()),
	)
	service.RegisterRoutes(s.engine.Group("/"))
}

// Use adds middleware to the engine
func (s *server) Use(middleware ...gin.HandlerFunc) {
	s.engine.Use(middleware...)
}

// RegisterRoutes registers a route with specific HTTP methods
func (s *server) RegisterRoutes(path string, handler gin.HandlerFunc, methods ...string) {
	for _, method := range methods {
		switch method {
		case "GET":
			s.engine.GET(path, handler)
		case "POST":
			s.engine.POST(path, handler)
		case "PUT":
			s.engine.PUT(path, handler)
		case "DELETE":
			s.engine.DELETE(path, handler)
		case "PATCH":
			s.engine.PATCH(path, handler)
		case "OPTIONS":
			s.engine.OPTIONS(path, handler)
		case "HEAD":
			s.engine.HEAD(path, handler)
		default:
			s.logger.Warn("unsupported HTTP method",
				zap.String("method", method),
				zap.String("path", path),
			)
		}
	}
}

// --- Health check handlers ---

func (s *server) livenessHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "alive",
		"id":     s.id,
		"name":   s.name,
	})
}

func (s *server) readinessHandler(c *gin.Context) {
	s.mu.RLock()
	status := s.status
	lastErr := s.lastErr
	s.mu.RUnlock()

	httpStatus := http.StatusOK
	if status != manager.StatusRunning {
		httpStatus = http.StatusServiceUnavailable
	}

	response := gin.H{
		"status": string(status),
		"id":     s.id,
		"name":   s.name,
	}

	if lastErr != nil {
		response["error"] = lastErr.Error()
	}

	c.JSON(httpStatus, response)
}

// GetStackTrace returns the current stack trace (useful for debugging)
func GetStackTrace() string {
	return string(debug.Stack())
}
