package web

import "time"

// Config holds configuration for the Web Server
type Config struct {
	// Port is the TCP port to listen on
	Port int `mapstructure:"port"`

	// ReadTimeout is the maximum duration for reading the entire request, including the body
	ReadTimeout time.Duration `mapstructure:"read_timeout"`

	// WriteTimeout is the maximum duration before timing out writes of the response
	WriteTimeout time.Duration `mapstructure:"write_timeout"`

	// ShutdownTimeout is the duration to wait for active connections to close during shutdown
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`

	// Mode is the Gin mode (debug, release, test)
	Mode string `mapstructure:"mode"`

	// Metrics configuration
	Metrics MetricsConfig `mapstructure:"metrics"`

	// Tracing configuration
	Tracing TracingConfig `mapstructure:"tracing"`

	// TLS configuration
	TLS TLSConfig `mapstructure:"tls"`

	// CORS configuration
	CORS CORSConfig `mapstructure:"cors"`

	// Security configuration
	Security SecurityConfig `mapstructure:"security"`

	// RateLimit configuration
	RateLimit RateLimitConfig `mapstructure:"rate_limit"`

	// Swagger configuration
	Swagger SwaggerConfig `mapstructure:"swagger"`
}

// MetricsConfig holds configuration for metrics
type MetricsConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Path    string `mapstructure:"path"`
}

// TracingConfig holds configuration for tracing
type TracingConfig struct {
	Enabled     bool   `mapstructure:"enabled"`
	ServiceName string `mapstructure:"service_name"`
}

// TLSConfig holds configuration for TLS
type TLSConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	CertFile string `mapstructure:"cert_file"`
	KeyFile  string `mapstructure:"key_file"`
}

// CORSConfig holds configuration for CORS
type CORSConfig struct {
	Enabled          bool     `mapstructure:"enabled"`
	AllowedOrigins   []string `mapstructure:"allowed_origins"`
	AllowedMethods   []string `mapstructure:"allowed_methods"`
	AllowedHeaders   []string `mapstructure:"allowed_headers"`
	ExposedHeaders   []string `mapstructure:"exposed_headers"`
	AllowCredentials bool     `mapstructure:"allow_credentials"`
	MaxAge           int      `mapstructure:"max_age"`
}

// SecurityConfig holds configuration for security headers
type SecurityConfig struct {
	Enabled               bool              `mapstructure:"enabled"`
	XSSProtection         string            `mapstructure:"xss_protection"`
	ContentTypeNosniff    string            `mapstructure:"content_type_nosniff"`
	XFrameOptions         string            `mapstructure:"x_frame_options"`
	HSTSMaxAge            int               `mapstructure:"hsts_max_age"`
	HSTSExcludeSubdomains bool              `mapstructure:"hsts_exclude_subdomains"`
	ContentSecurityPolicy string            `mapstructure:"content_security_policy"`
	ReferrerPolicy        string            `mapstructure:"referrer_policy"`
	CustomHeaders         map[string]string `mapstructure:"custom_headers"`
}

// RateLimitConfig holds configuration for rate limiting
type RateLimitConfig struct {
	Enabled           bool    `mapstructure:"enabled"`
	RequestsPerSecond float64 `mapstructure:"requests_per_second"`
	Burst             int     `mapstructure:"burst"`
}

// SwaggerConfig holds configuration for Swagger documentation
type SwaggerConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Path    string `mapstructure:"path"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() Config {
	return Config{
		Port:            8080,
		ReadTimeout:     10 * time.Second,
		WriteTimeout:    10 * time.Second,
		ShutdownTimeout: 5 * time.Second,
		Mode:            "release",
		Metrics: MetricsConfig{
			Enabled: true,
			Path:    "/metrics",
		},
		Tracing: TracingConfig{
			Enabled:     true,
			ServiceName: "grouter-web",
		},
		RateLimit: RateLimitConfig{
			Enabled:           true,
			RequestsPerSecond: 100,
			Burst:             200,
		},
		Swagger: SwaggerConfig{
			Enabled: true,
			Path:    "/swagger",
		},
	}
}
