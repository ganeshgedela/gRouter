package config

import "time"

// Config represents the complete application configuration
type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Manager  ManagerConfig  `mapstructure:"manager"`
	NATS     NATSConfig     `mapstructure:"nats"`
	GRPC     GRPCConfig     `mapstructure:"grpc"`
	Log      LogConfig      `mapstructure:"log"`
	Web      WebConfig      `mapstructure:"web"`
	Tracing  TracingConfig  `mapstructure:"tracing"`
	Services ServicesConfig `mapstructure:"services"`
	Database DatabaseConfig `mapstructure:"database"`
}

// ManagerConfig holds configuration for the ServiceManager
type ManagerConfig struct {
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
	MonitorInterval time.Duration `mapstructure:"monitor_interval"`
}

// AppConfig holds application-level settings
type AppConfig struct {
	Name        string `mapstructure:"name"`
	Version     string `mapstructure:"version"`
	Environment string `mapstructure:"environment"`
}

// NATSConfig holds NATS connection settings
// NATSConfig holds NATS connection settings
type NATSConfig struct {
	Enabled           bool                 `mapstructure:"enabled"`
	URL               string               `mapstructure:"url"`
	MaxReconnects     int                  `mapstructure:"max_reconnects"`
	ReconnectWait     time.Duration        `mapstructure:"reconnect_wait"`
	ConnectionTimeout time.Duration        `mapstructure:"connection_timeout"`
	DrainTimeout      time.Duration        `mapstructure:"drain_timeout"`
	Auth              AuthConfig           `mapstructure:"auth"`
	TLS               TLSConfig            `mapstructure:"tls"`
	Middleware        NATSMiddlewareConfig `mapstructure:"middleware"`
}

// NATSMiddlewareConfig holds configuration for NATS middleware
type NATSMiddlewareConfig struct {
	Recovery       MiddlewareState      `mapstructure:"recovery"`
	Metrics        MetricsConfig        `mapstructure:"metrics"`
	Tracing        MiddlewareState      `mapstructure:"tracing"`
	Logging        LoggingConfig        `mapstructure:"logging"`
	CircuitBreaker CircuitBreakerConfig `mapstructure:"circuit_breaker"`
	Retry          RetryConfig          `mapstructure:"retry"`
	Timeout        TimeoutConfig        `mapstructure:"timeout"`
	RateLimit      RateLimitConfig      `mapstructure:"rate_limit"`
}

// TimeoutConfig holds configuration for timeout middleware
type TimeoutConfig struct {
	Enabled bool          `mapstructure:"enabled"`
	Default time.Duration `mapstructure:"default"`
}

// RetryConfig holds configuration for retry middleware
type RetryConfig struct {
	Enabled         bool          `mapstructure:"enabled"`
	MaxAttempts     int           `mapstructure:"max_attempts"`
	InitialInterval time.Duration `mapstructure:"initial_interval"`
	Multiplier      float64       `mapstructure:"multiplier"`
	MaxInterval     time.Duration `mapstructure:"max_interval"`
}

// CircuitBreakerConfig holds configuration for circuit breaker middleware
type CircuitBreakerConfig struct {
	Enabled       bool          `mapstructure:"enabled"`
	MaxRequests   uint32        `mapstructure:"max_requests"`
	Interval      time.Duration `mapstructure:"interval"`
	Timeout       time.Duration `mapstructure:"timeout"`
	TripThreshold uint32        `mapstructure:"trip_threshold"` // Consecutive failures to trip
}

// MiddlewareState holds generic enabled state
type MiddlewareState struct {
	Enabled bool `mapstructure:"enabled"`
}

// LoggingConfig holds configuration for logging middleware
type LoggingConfig struct {
	Enabled bool `mapstructure:"enabled"`
}

// LogConfig holds logging configuration
type LogConfig struct {
	Level      string `mapstructure:"level"`
	Format     string `mapstructure:"format"` // json or console
	OutputPath string `mapstructure:"output_path"`
}

// WebConfig holds web server configuration
type WebConfig struct {
	Enabled            bool            `mapstructure:"enabled"`
	Port               int             `mapstructure:"port"`
	ReadTimeout        time.Duration   `mapstructure:"read_timeout"`
	WriteTimeout       time.Duration   `mapstructure:"write_timeout"`
	ShutdownTimeout    time.Duration   `mapstructure:"shutdown_timeout"`
	Mode               string          `mapstructure:"mode"`
	MaxRequestBodySize int64           `mapstructure:"max_request_body_size"` // In bytes, 0 = 10MB default
	Metrics            MetricsConfig   `mapstructure:"metrics"`
	TLS                TLSConfig       `mapstructure:"tls"`
	CORS               CORSConfig      `mapstructure:"cors"`
	Security           SecurityConfig  `mapstructure:"security"`
	RateLimit          RateLimitConfig `mapstructure:"rate_limit"`
	Swagger            SwaggerConfig   `mapstructure:"swagger"`
	Logging            LoggingConfig   `mapstructure:"logging"`
	Auth               AuthConfig      `mapstructure:"auth"`
}

type AuthConfig struct {
	Enabled   bool   `mapstructure:"enabled"`
	Issuer    string `mapstructure:"issuer"`
	Audience  string `mapstructure:"audience"`
	Token     string `mapstructure:"token"`
	Username  string `mapstructure:"username"`
	Password  string `mapstructure:"password"`
	CredsFile string `mapstructure:"creds_file"`
}

// TLSConfig holds configuration for TLS
type TLSConfig struct {
	Enabled    bool   `mapstructure:"enabled"`
	CertFile   string `mapstructure:"cert_file"`
	KeyFile    string `mapstructure:"key_file"`
	CAFile     string `mapstructure:"ca_file"`
	SkipVerify bool   `mapstructure:"skip_verify"`
}

// GRPCConfig holds gRPC server configuration
type GRPCConfig struct {
	Enabled           bool          `mapstructure:"enabled"`
	Port              int           `mapstructure:"port"`
	KeepaliveTime     time.Duration `mapstructure:"keepalive_time"`
	KeepaliveTimeout  time.Duration `mapstructure:"keepalive_timeout"`
	TLS               TLSConfig     `mapstructure:"tls"`
	ReflectionEnabled bool          `mapstructure:"reflection_enabled"`
	GatewayEnabled    bool          `mapstructure:"gateway_enabled"`
	GatewayPort       int           `mapstructure:"gateway_port"`
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

// MetricsConfig holds configuration for metrics
type MetricsConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Path    string `mapstructure:"path"`
}

// ServicesConfig holds service-specific settings
type ServicesConfig map[string]interface{}

// TracingConfig holds OpenTelemetry tracing configuration
type TracingConfig struct {
	Enabled     bool   `mapstructure:"enabled"`
	ServiceName string `mapstructure:"service_name"`
	Exporter    string `mapstructure:"exporter"` // e.g., "jaeger", "stdout"
	Endpoint    string `mapstructure:"endpoint"` // e.g., "http://localhost:14268/api/traces"
}

// DatabaseConfig holds database connection settings
type DatabaseConfig struct {
	Driver          string        `mapstructure:"driver"` // postgres, sqlite, etc.
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	DBName          string        `mapstructure:"dbname"`
	SSLMode         string        `mapstructure:"ssl_mode"`
	LogLevel        string        `mapstructure:"log_level"` // silent, error, warn, info
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}
