package nats

import (
	"crypto/tls"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

// Client wraps NATS connection
type Client struct {
	conn   *nats.Conn
	js     nats.JetStreamContext
	logger *zap.Logger
	source string
	config Config
}

// Config holds NATS client configuration
type Config struct {
	URL               string        `mapstructure:"url"`
	MaxReconnects     int           `mapstructure:"max_reconnects"`
	ReconnectWait     time.Duration `mapstructure:"reconnect_wait"`
	ConnectionTimeout time.Duration `mapstructure:"connection_timeout"`
	DrainTimeout      time.Duration `mapstructure:"drain_timeout"`
	// Authentication configuration
	Auth AuthConfig `mapstructure:"auth"`
	// TLS configuration
	TLS TLSConfig `mapstructure:"tls"`
	// Middleware configuration
	Middleware NATSMiddlewareConfig `mapstructure:"middleware"`
}

// AuthConfig holds authentication configuration
type AuthConfig struct {
	Token     string `mapstructure:"token"`
	Username  string `mapstructure:"username"`
	Password  string `mapstructure:"password"`
	CredsFile string `mapstructure:"creds_file"`
}

// TLSConfig holds TLS configuration
type TLSConfig struct {
	Enabled    bool   `mapstructure:"enabled"`
	SkipVerify bool   `mapstructure:"skip_verify"`
	CAFile     string `mapstructure:"ca_file"`
	CertFile   string `mapstructure:"cert_file"`
	KeyFile    string `mapstructure:"key_file"`
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

// RateLimitConfig holds configuration for rate limiting
type RateLimitConfig struct {
	Enabled           bool    `mapstructure:"enabled"`
	RequestsPerSecond float64 `mapstructure:"requests_per_second"`
	Burst             int     `mapstructure:"burst"`
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

// CircuitBreakerConfig holds configuration for circuit breaker
type CircuitBreakerConfig struct {
	Enabled       bool          `mapstructure:"enabled"`
	MaxRequests   uint32        `mapstructure:"max_requests"`
	Interval      time.Duration `mapstructure:"interval"`
	Timeout       time.Duration `mapstructure:"timeout"`
	TripThreshold uint32        `mapstructure:"trip_threshold"`
}

// MiddlewareState holds generic enabled state
type MiddlewareState struct {
	Enabled bool `mapstructure:"enabled"`
}

// MetricsConfig holds configuration for metrics
type MetricsConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Path    string `mapstructure:"path"`
}

// LoggingConfig holds configuration for logging
type LoggingConfig struct {
	Enabled bool `mapstructure:"enabled"`
}

// TracingConfig holds configuration for tracing
type TracingConfig struct {
	Enabled bool `mapstructure:"enabled"`
}

// NewNATSClient creates a new NATS client
func NewNATSClient(cfg Config, logger *zap.Logger) (*Client, error) {
	if logger == nil {
		return nil, fmt.Errorf("logger is required")
	}

	return &Client{
		config: cfg,
		logger: logger,
	}, nil
}

// Connect establishes connection to NATS server
func (c *Client) Connect() error {
	opts := c.buildConnectionOptions()

	conn, err := nats.Connect(c.config.URL, opts...)
	if err != nil {
		return fmt.Errorf("failed to connect to NATS: %w", err)
	}

	c.conn = conn

	// Create JetStream Context
	js, err := conn.JetStream()
	if err != nil {
		conn.Close()
		return fmt.Errorf("failed to create JetStream context: %w", err)
	}
	c.js = js

	return nil
}

// buildConnectionOptions creates the NATS connection options based on config
func (c *Client) buildConnectionOptions() []nats.Option {
	logger := c.logger
	opts := []nats.Option{
		nats.MaxReconnects(c.config.MaxReconnects),
		nats.ReconnectWait(c.config.ReconnectWait),
		nats.Timeout(c.config.ConnectionTimeout),
		nats.RetryOnFailedConnect(true),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			if err != nil {
				logger.Error("NATS disconnected", zap.Error(err))
			}
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			logger.Debug("NATS reconnected", zap.String("url", nc.ConnectedUrl()))
		}),
	}

	// Add authentication if provided
	if c.config.Auth.CredsFile != "" {
		opts = append(opts, nats.UserCredentials(c.config.Auth.CredsFile))
	} else if c.config.Auth.Token != "" {
		opts = append(opts, nats.Token(c.config.Auth.Token))
	} else if c.config.Auth.Username != "" && c.config.Auth.Password != "" {
		opts = append(opts, nats.UserInfo(c.config.Auth.Username, c.config.Auth.Password))
	}

	// Add TLS if enabled
	if c.config.TLS.Enabled {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: c.config.TLS.SkipVerify,
		}
		if c.config.TLS.CAFile != "" {
			opts = append(opts, nats.RootCAs(c.config.TLS.CAFile))
		}
		if c.config.TLS.CertFile != "" && c.config.TLS.KeyFile != "" {
			opts = append(opts, nats.ClientCert(c.config.TLS.CertFile, c.config.TLS.KeyFile))
		}

		// Apply the base secure config (handles SkipVerify)
		// If we use nats.Secure(tlsConfig), we provide the base config.
		// Then we can append RootCAs and ClientCert which will modify the connection's TLS state.
		opts = append(opts, nats.Secure(tlsConfig))
	}

	return opts
}

// Close gracefully closes the NATS connection
func (c *Client) Close() error {
	if c.conn != nil {
		c.conn.Drain()
		c.conn.Close()
		c.logger.Debug("NATS connection closed")
	}
	return nil
}

// IsConnected returns true if connected to NATS
func (c *Client) IsConnected() bool {
	return c.conn != nil && c.conn.IsConnected()
}

// Conn returns the underlying NATS connection
func (c *Client) Conn() *nats.Conn {
	return c.conn
}

// JetStream returns the JetStream context, initializing it if necessary
func (c *Client) JetStream() (nats.JetStreamContext, error) {
	if c.js != nil {
		return c.js, nil
	}

	if !c.IsConnected() {
		return nil, fmt.Errorf("not connected to NATS")
	}

	js, err := c.conn.JetStream()
	if err != nil {
		return nil, fmt.Errorf("failed to create JetStream context: %w", err)
	}

	c.js = js
	return js, nil
}

// Ping checks the connection health
func (c *Client) Ping() error {
	if !c.IsConnected() {
		return fmt.Errorf("nats client not connected")
	}
	return nil
}
