package grpc

import (
	"context"
	"time"

	"google.golang.org/grpc"
)

// TLSConfig holds TLS/mTLS configuration
type TLSConfig struct {
	// Enabled indicates if TLS should be used
	Enabled bool

	// Server-side certificates
	CertFile string // Path to server certificate file
	KeyFile  string // Path to server private key file

	// CA certificate for client verification (mTLS)
	CAFile string // Path to CA certificate file

	// Client authentication
	ClientAuth bool // Require client certificates (mTLS)

	// Client-side settings
	ServerName         string // Expected server name for verification
	InsecureSkipVerify bool   // Skip server certificate verification (NOT recommended for production)
}

// ServerConfig holds configuration for the gRPC server.
type ServerConfig struct {
	Port              int
	ReflectionEnabled bool
	GatewayEnabled    bool
	GatewayPort       int
	TLS               TLSConfig // TLS configuration
	KeepaliveTime     time.Duration
	KeepaliveTimeout  time.Duration

	// Interceptor configurations
	Interceptors InterceptorConfig
}

// ClientConfig holds configuration for the gRPC client.
type ClientConfig struct {
	Target           string
	DialTimeout      time.Duration
	Timeout          time.Duration // Request timeout
	TLS              TLSConfig     // TLS configuration
	MaxRetries       int
	RetryBackoff     time.Duration
	KeepAliveTime    time.Duration
	KeepAliveTimeout time.Duration

	// Interceptor configurations
	Interceptors InterceptorConfig
}

// InterceptorConfig holds configuration for all interceptors.
type InterceptorConfig struct {
	Logging        LoggingConfig
	Metrics        MetricsConfig
	Tracing        TracingConfig
	Recovery       RecoveryConfig
	Retry          RetryConfig
	Timeout        TimeoutConfig
	RateLimit      RateLimitConfig
	CircuitBreaker CircuitBreakerConfig
	Auth           AuthConfig
	Validator      ValidatorConfig
}

// LoggingConfig configures logging interceptor.
type LoggingConfig struct {
	Enabled bool
	Level   string // debug, info, warn, error
}

// MetricsConfig configures metrics interceptor.
type MetricsConfig struct {
	Enabled bool
}

// TracingConfig configures tracing interceptor.
type TracingConfig struct {
	Enabled bool
}

// RecoveryConfig configures panic recovery interceptor.
type RecoveryConfig struct {
	Enabled bool
}

// RetryConfig configures retry interceptor (client-side).
type RetryConfig struct {
	Enabled     bool
	MaxAttempts int
	InitialWait time.Duration
	MaxWait     time.Duration
	Multiplier  float64
}

// TimeoutConfig configures timeout interceptor.
type TimeoutConfig struct {
	Enabled  bool
	Duration time.Duration
}

// RateLimitConfig configures rate limiting interceptor.
type RateLimitConfig struct {
	Enabled bool
	Rate    int // requests per second
	Burst   int // burst size
}

// CircuitBreakerConfig configures circuit breaker interceptor (client-side).
type CircuitBreakerConfig struct {
	Enabled     bool
	Threshold   int           // failures before open
	Timeout     time.Duration // time before half-open
	MaxRequests int           // max requests in half-open
}

// AuthConfig configures authentication interceptor.
type AuthConfig struct {
	Enabled   bool
	JWTSecret string
	APIKeys   []string
}

// ValidatorConfig configures validation interceptor.
type ValidatorConfig struct {
	Enabled bool
}

// UnaryServerInterceptor is the type for unary server interceptors.
type UnaryServerInterceptor = grpc.UnaryServerInterceptor

// StreamServerInterceptor is the type for stream server interceptors.
type StreamServerInterceptor = grpc.StreamServerInterceptor

// UnaryClientInterceptor is the type for unary client interceptors.
type UnaryClientInterceptor = grpc.UnaryClientInterceptor

// StreamClientInterceptor is the type for stream client interceptors.
type StreamClientInterceptor = grpc.StreamClientInterceptor

// ServerInterceptors holds all server interceptors.
type ServerInterceptors struct {
	Unary  []UnaryServerInterceptor
	Stream []StreamServerInterceptor
}

// ClientInterceptors holds all client interceptors.
type ClientInterceptors struct {
	Unary  []UnaryClientInterceptor
	Stream []StreamClientInterceptor
}

// HandlerFunc is the signature for RPC handlers in recovery/middleware.
type HandlerFunc func(ctx context.Context, req interface{}) (interface{}, error)

// StreamHandlerFunc is the signature for stream handlers.
type StreamHandlerFunc func(srv interface{}, stream grpc.ServerStream) error
