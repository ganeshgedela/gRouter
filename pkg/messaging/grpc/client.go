package grpc

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

// Client wraps a gRPC client connection with middleware support.
type Client struct {
	conn   *grpc.ClientConn
	config ClientConfig
	logger *zap.Logger

	// Interceptors
	unaryInterceptors  []UnaryClientInterceptor
	streamInterceptors []StreamClientInterceptor
}

// ClientOption defines a functional option for configuring the Client.
type ClientOption func(*Client)

// WithClientInterceptors adds unary and stream interceptors to the client.
func WithClientInterceptors(unary []UnaryClientInterceptor, stream []StreamClientInterceptor) ClientOption {
	return func(c *Client) {
		c.unaryInterceptors = append(c.unaryInterceptors, unary...)
		c.streamInterceptors = append(c.streamInterceptors, stream...)
	}
}

// WithClientUnaryInterceptor adds a unary interceptor to the client.
func WithClientUnaryInterceptor(interceptor UnaryClientInterceptor) ClientOption {
	return func(c *Client) {
		c.unaryInterceptors = append(c.unaryInterceptors, interceptor)
	}
}

// WithClientStreamInterceptor adds a stream interceptor to the client.
func WithClientStreamInterceptor(interceptor StreamClientInterceptor) ClientOption {
	return func(c *Client) {
		c.streamInterceptors = append(c.streamInterceptors, interceptor)
	}
}

// NewClient creates a new gRPC client with the given configuration.
func NewClient(config ClientConfig, logger *zap.Logger, opts ...ClientOption) (*Client, error) {
	c := &Client{
		config: config,
		logger: logger,
	}

	// Apply options
	for _, opt := range opts {
		opt(c)
	}

	// Build dial options
	var dialOpts []grpc.DialOption

	// Add TLS credentials if enabled
	if config.TLS.Enabled {
		creds, err := configureClientTLS(config.TLS)
		if err != nil {
			return nil, fmt.Errorf("failed to configure TLS: %w", err)
		}
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(creds))
		logger.Info("TLS enabled for gRPC client",
			zap.String("target", config.Target),
			zap.String("server_name", config.TLS.ServerName),
		)
	} else {
		// Use insecure credentials
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	// Add keepalive parameters
	if config.KeepAliveTime > 0 {
		dialOpts = append(dialOpts, grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                config.KeepAliveTime,
			Timeout:             config.KeepAliveTimeout,
			PermitWithoutStream: true,
		}))
	}

	// Add interceptors
	if len(c.unaryInterceptors) > 0 {
		dialOpts = append(dialOpts, grpc.WithChainUnaryInterceptor(c.unaryInterceptors...))
	}
	if len(c.streamInterceptors) > 0 {
		dialOpts = append(dialOpts, grpc.WithChainStreamInterceptor(c.streamInterceptors...))
	}

	// Set default timeout if not specified
	timeout := config.Timeout
	if timeout == 0 {
		timeout = 10 * time.Second
	}

	// Create context with timeout for dial
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Dial the target
	conn, err := grpc.DialContext(ctx, config.Target, dialOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to dial %s: %w", config.Target, err)
	}

	c.conn = conn
	logger.Info("gRPC client connected", zap.String("target", config.Target))

	return c, nil
}

// GetConn returns the underlying gRPC connection.
// This can be used to create service clients.
func (c *Client) GetConn() *grpc.ClientConn {
	return c.conn
}

// Close closes the client connection.
func (c *Client) Close() error {
	if c.conn != nil {
		c.logger.Info("Closing gRPC client connection")
		return c.conn.Close()
	}
	return nil
}

// Invoke performs a unary RPC call.
// This is a convenience method for testing or direct invocation.
func (c *Client) Invoke(ctx context.Context, method string, args interface{}, reply interface{}, opts ...grpc.CallOption) error {
	return c.conn.Invoke(ctx, method, args, reply, opts...)
}

// NewStream creates a new stream.
// This is a convenience method for testing or direct stream creation.
func (c *Client) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	return c.conn.NewStream(ctx, desc, method, opts...)
}
