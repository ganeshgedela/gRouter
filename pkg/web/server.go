package web

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/secure"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.uber.org/zap"

	_ "grouter/docs" // Import generated docs
	"grouter/pkg/health"
)

// Server wraps the Gin engine and manages the HTTP server lifecycle
// @title gRouter API
// @version 1.0
// @description Generic Web Framework for gRouter
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
// @schemes http https
type Server struct {
	engine *gin.Engine
	server *http.Server
	cfg    Config
	logger *zap.Logger
	health *health.HealthService
}

func InitEngine(cfg Config, logger *zap.Logger) *gin.Engine {
	engine := gin.New()
	engine.Use(RequestIDMiddleware())
	engine.Use(gin.Recovery())
	engine.Use(LoggerMiddleware(logger))

	if cfg.Tracing.Enabled {
		engine.Use(otelgin.Middleware(cfg.Tracing.ServiceName))
	}

	if cfg.CORS.Enabled {
		corsConfig := cors.DefaultConfig()
		if len(cfg.CORS.AllowedOrigins) > 0 {
			corsConfig.AllowOrigins = cfg.CORS.AllowedOrigins
		} else {
			corsConfig.AllowAllOrigins = true
		}
		if len(cfg.CORS.AllowedMethods) > 0 {
			corsConfig.AllowMethods = cfg.CORS.AllowedMethods
		}
		if len(cfg.CORS.AllowedHeaders) > 0 {
			corsConfig.AllowHeaders = cfg.CORS.AllowedHeaders
		}
		if len(cfg.CORS.ExposedHeaders) > 0 {
			corsConfig.ExposeHeaders = cfg.CORS.ExposedHeaders
		}
		corsConfig.AllowCredentials = cfg.CORS.AllowCredentials
		if cfg.CORS.MaxAge > 0 {
			corsConfig.MaxAge = time.Duration(cfg.CORS.MaxAge) * time.Second
		}
		engine.Use(cors.New(corsConfig))
	}

	if cfg.Security.Enabled {
		secureConfig := secure.DefaultConfig()
		if cfg.Security.XSSProtection != "" {
			secureConfig.BrowserXssFilter = true
		}
		if cfg.Security.ContentTypeNosniff != "" {
			secureConfig.ContentTypeNosniff = true
		}
		if cfg.Security.XFrameOptions != "" {
			secureConfig.FrameDeny = cfg.Security.XFrameOptions == "DENY"
			secureConfig.CustomFrameOptionsValue = cfg.Security.XFrameOptions
		}
		if cfg.Security.HSTSMaxAge > 0 {
			secureConfig.STSSeconds = int64(cfg.Security.HSTSMaxAge)
			secureConfig.STSIncludeSubdomains = cfg.Security.HSTSExcludeSubdomains
		}
		if cfg.Security.ContentSecurityPolicy != "" {
			secureConfig.ContentSecurityPolicy = cfg.Security.ContentSecurityPolicy
		}

		if cfg.Security.ReferrerPolicy != "" {
			secureConfig.ReferrerPolicy = cfg.Security.ReferrerPolicy
		}

		// Disable SSL Redirect if TLS is not enabled
		if !cfg.TLS.Enabled {
			secureConfig.SSLRedirect = false
		} else {
			secureConfig.SSLRedirect = true
		}

		// If in development/debug mode, we might want to relax some security settings
		if cfg.Mode == "debug" {
			secureConfig.IsDevelopment = true
		}

		engine.Use(secure.New(secureConfig))
	}

	if cfg.RateLimit.Enabled {
		engine.Use(RateLimitMiddleware(cfg.RateLimit.RequestsPerSecond, cfg.RateLimit.Burst))
	}

	if cfg.Metrics.Enabled {
		engine.Use(MetricsMiddleware())
		// Register metrics handler
		path := cfg.Metrics.Path
		if path == "" {
			path = "/metrics"
		}
		engine.GET(path, gin.WrapH(promhttp.Handler()))
	}

	if cfg.Swagger.Enabled {
		path := cfg.Swagger.Path
		if path == "" {
			path = "/swagger"
		}
		engine.GET(path+"/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
	return engine
}

// NewWebServer creates a new Web Server instance
func NewWebServer(cfg Config, logger *zap.Logger, healthSvc *health.HealthService) *Server {
	// Set Gin mode
	gin.SetMode(cfg.Mode)

	engine := InitEngine(cfg, logger)

	server := &Server{
		engine: engine,
		cfg:    cfg,
		logger: logger,
		health: healthSvc,
	}

	// Register health handlers
	if healthSvc != nil {
		server.engine.GET("/health/live", healthSvc.LivenessHandler)
		server.engine.GET("/health/ready", healthSvc.ReadinessHandler)
	}
	return server
}

// RegisterService registers a service's routes with the server
func (s *Server) RegisterWebService(service WebService) {
	service.RegisterRoutes(s.engine.Group("/"))
}

// Health returns the underlying health service
func (s *Server) Health() *health.HealthService {
	return s.health
}

// Start starts the HTTP server
func (s *Server) Start() error {
	s.server = &http.Server{
		Addr:         fmt.Sprintf(":%d", s.cfg.Port),
		Handler:      s.engine,
		ReadTimeout:  s.cfg.ReadTimeout,
		WriteTimeout: s.cfg.WriteTimeout,
	}

	s.logger.Info("Starting web server", zap.Int("port", s.cfg.Port), zap.Bool("tls", s.cfg.TLS.Enabled))

	go func() {
		var err error
		if s.cfg.TLS.Enabled {
			if s.cfg.TLS.CertFile == "" || s.cfg.TLS.KeyFile == "" {
				s.logger.Fatal("TLS enabled but cert or key file missing")
			}
			err = s.server.ListenAndServeTLS(s.cfg.TLS.CertFile, s.cfg.TLS.KeyFile)
		} else {
			err = s.server.ListenAndServe()
		}

		if err != nil && err != http.ErrServerClosed {
			// In a restart scenario, we might just log error instead of Fatal if it's transient
			// But for now, sticking to Fatal for critical failures, except we can't Fatal in restart loop ideally.
			// Let's degrade to Error for robustness if it was a restart.
			s.logger.Error("Web server stopped", zap.Error(err))
		}
	}()

	return nil
}

// Stop gracefully shuts down the HTTP server
func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("Stopping web server")

	if s.server == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(ctx, s.cfg.ShutdownTimeout)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		// Attempt force close if shutdown fails
		s.server.Close()
		return fmt.Errorf("web server forced to shutdown: %w", err)
	}

	return nil
}

// Restart stops and starts the web server
func (s *Server) ResetEngine(ctx context.Context) error {
	s.logger.Info("Restarting web server...")

	if err := s.Stop(ctx); err != nil {
		s.logger.Error("Failed to stop server during restart", zap.Error(err))
		// Proceeding to start anyway
	}
	s.engine = nil
	// Small delay to allow port release
	time.Sleep(1 * time.Second)

	s.engine = InitEngine(s.cfg, s.logger)
	if s.health != nil {
		s.engine.GET("/health/live", s.health.LivenessHandler)
		s.engine.GET("/health/ready", s.health.ReadinessHandler)
	}
	return nil
}

// LoggerMiddleware logs HTTP requests using zap
func LoggerMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		if len(c.Errors) > 0 {
			for _, e := range c.Errors.Errors() {
				logger.Error(e)
			}
		} else {
			logger.Info("HTTP Request",
				zap.String("request_id", c.GetString("RequestID")),
				zap.Int("status", c.Writer.Status()),
				zap.String("method", c.Request.Method),
				zap.String("path", path),
				zap.String("query", query),
				zap.String("ip", c.ClientIP()),
				zap.String("user-agent", c.Request.UserAgent()),
				zap.Duration("latency", latency),
			)
		}
	}
}
