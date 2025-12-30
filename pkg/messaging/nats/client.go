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
	config Config
}

// Config holds NATS client configuration
type Config struct {
	URL               string        `mapstructure:"url"`
	MaxReconnects     int           `mapstructure:"max_reconnects"`
	ReconnectWait     time.Duration `mapstructure:"reconnect_wait"`
	ConnectionTimeout time.Duration `mapstructure:"connection_timeout"`
	Token             string        `mapstructure:"token"`
	Username          string        `mapstructure:"username"`
	Password          string        `mapstructure:"password"`
	// TLS configuration
	UseTLS     bool   `mapstructure:"use_tls"`
	SkipVerify bool   `mapstructure:"skip_verify"`
	CAFile     string `mapstructure:"ca_file"`
	CertFile   string `mapstructure:"cert_file"`
	KeyFile    string `mapstructure:"key_file"`
	// NATS 2.0+ Credentials
	CredsFile string `mapstructure:"creds_file"`
	// Metrics configuration
	Metrics MetricsConfig `mapstructure:"metrics"`
	// Logging configuration
	Logging LoggingConfig `mapstructure:"logging"`
	// Tracing configuration
	Tracing TracingConfig `mapstructure:"tracing"`
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
	opts := []nats.Option{
		nats.MaxReconnects(c.config.MaxReconnects),
		nats.ReconnectWait(c.config.ReconnectWait),
		nats.Timeout(c.config.ConnectionTimeout),
		nats.RetryOnFailedConnect(true),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			if err != nil {
				c.logger.Error("NATS disconnected", zap.Error(err))
			}
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			c.logger.Info("NATS reconnected", zap.String("url", nc.ConnectedUrl()))
		}),
		nats.ClosedHandler(func(nc *nats.Conn) {
			c.logger.Warn("NATS connection closed")
		}),
	}

	// Add authentication if provided
	if c.config.CredsFile != "" {
		opts = append(opts, nats.UserCredentials(c.config.CredsFile))
	} else if c.config.Token != "" {
		opts = append(opts, nats.Token(c.config.Token))
	} else if c.config.Username != "" && c.config.Password != "" {
		opts = append(opts, nats.UserInfo(c.config.Username, c.config.Password))
	}

	// Add TLS if enabled
	if c.config.UseTLS {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: c.config.SkipVerify,
		}
		if c.config.CAFile != "" {
			opts = append(opts, nats.RootCAs(c.config.CAFile))
		}
		if c.config.CertFile != "" && c.config.KeyFile != "" {
			opts = append(opts, nats.ClientCert(c.config.CertFile, c.config.KeyFile))
		}

		// If custom TLS config is needed beyond just files (e.g. SkipVerify is already handled)
		// We can still use nats.Secure(tlsConfig) but RootCAs and ClientCert helper options
		// read the files directly which is often safer/easier.
		// However, nats.Secure overwrites the TLS config, so we should be careful mixing them.
		// The helper options modify the internal TLS config.
		// If SkipVerify is set, we still need to ensure that applies.
		// Let's rely on the helper options for certs, and manual Secure() for SkipVerify if needed,
		// but nats.Secure() takes a *tls.Config.

		// Better approach:
		// If we use nats.Secure(tlsConfig), we provide the base config.
		// Then we can append RootCAs and ClientCert which will modify the connection's TLS state.
		opts = append(opts, nats.Secure(tlsConfig))
	}

	conn, err := nats.Connect(c.config.URL, opts...)
	if err != nil {
		return fmt.Errorf("failed to connect to NATS: %w", err)
	}

	c.conn = conn
	if c.conn.IsConnected() {
		c.logger.Info("Connected to NATS", zap.String("url", c.config.URL))
	} else {
		c.logger.Warn("NATS connection established but not yet connected (reconnecting mode)", zap.String("url", c.config.URL))
	}
	return nil
}

// Close gracefully closes the NATS connection
func (c *Client) Close() error {
	if c.conn != nil {
		c.conn.Drain()
		c.conn.Close()
		c.logger.Info("NATS connection closed")
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
