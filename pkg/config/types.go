package config

import "time"

// Config represents the complete application configuration
type Config struct {
	App      AppConfig      `mapstructure:"app"`
	NATS     NATSConfig     `mapstructure:"nats"`
	Log      LogConfig      `mapstructure:"log"`
	Web      WebConfig      `mapstructure:"web"`
	Tracing  TracingConfig  `mapstructure:"tracing"`
	Services ServicesConfig `mapstructure:"services"`
	Database DatabaseConfig `mapstructure:"database"`
	Metrics  MetricsConfig  `mapstructure:"metrics"`
}

// AppConfig holds application-level settings
type AppConfig struct {
	Name        string `mapstructure:"name"`
	Version     string `mapstructure:"version"`
	Environment string `mapstructure:"environment"`
}

// NATSConfig holds NATS connection settings
type NATSConfig struct {
	Enabled           bool          `mapstructure:"enabled"`
	URL               string        `mapstructure:"url"`
	MaxReconnects     int           `mapstructure:"max_reconnects"`
	ReconnectWait     time.Duration `mapstructure:"reconnect_wait"`
	ConnectionTimeout time.Duration `mapstructure:"connection_timeout"`
	Token             string        `mapstructure:"token"`
	Username          string        `mapstructure:"username"`
	Password          string        `mapstructure:"password"`
	CredsFile         string        `mapstructure:"creds_file"`
	UseTLS            bool          `mapstructure:"use_tls"`
	SkipVerify        bool          `mapstructure:"skip_verify"`
	CAFile            string        `mapstructure:"ca_file"`
	CertFile          string        `mapstructure:"cert_file"`
	KeyFile           string        `mapstructure:"key_file"`
	Metrics           MetricsConfig `mapstructure:"metrics"`
	Logging           LoggingConfig `mapstructure:"logging"`
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
	Enabled         bool            `mapstructure:"enabled"`
	Port            int             `mapstructure:"port"`
	ReadTimeout     time.Duration   `mapstructure:"read_timeout"`
	WriteTimeout    time.Duration   `mapstructure:"write_timeout"`
	ShutdownTimeout time.Duration   `mapstructure:"shutdown_timeout"`
	Mode            string          `mapstructure:"mode"`
	Metrics         MetricsConfig   `mapstructure:"metrics"`
	TLS             TLSConfig       `mapstructure:"tls"`
	CORS            CORSConfig      `mapstructure:"cors"`
	Security        SecurityConfig  `mapstructure:"security"`
	RateLimit       RateLimitConfig `mapstructure:"rate_limit"`
	Swagger         SwaggerConfig   `mapstructure:"swagger"`
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
