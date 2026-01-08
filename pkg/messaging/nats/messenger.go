package nats

import (
	"context"
	"fmt"
	"grouter/pkg/config"
	"time"

	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

// Messenger wraps Client, Publisher, and Subscriber into a single unit.
type Messenger struct {
	Client     *Client
	Publisher  Publisher
	Subscriber Subscriber
}

func (m *Messenger) IsConnected() bool {
	return m.Client.IsConnected()
}

func (m *Messenger) Source() string {
	if m.Client == nil {
		return ""
	}
	return m.Client.source
}

// NewMessenger creates a new Messenger.
func NewMessenger(client *Client, pub Publisher, sub Subscriber) *Messenger {
	return &Messenger{
		Client:     client,
		Publisher:  pub,
		Subscriber: sub,
	}
}

// InitFromConfig initializes the messenger from configuration
func (m *Messenger) InitFromConfig(cfg *config.Config, logger *zap.Logger) error {
	if cfg == nil {
		return fmt.Errorf("config cannot be nil")
	}

	if err := m.validateConfig(cfg); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	// Transform global config to NATS client config
	natsCfg := Config{
		URL:               cfg.NATS.URL,
		MaxReconnects:     cfg.NATS.MaxReconnects,
		ReconnectWait:     cfg.NATS.ReconnectWait,
		ConnectionTimeout: cfg.NATS.ConnectionTimeout,
		DrainTimeout:      cfg.NATS.DrainTimeout,
		Auth: AuthConfig{
			Token:     cfg.NATS.Auth.Token,
			Username:  cfg.NATS.Auth.Username,
			Password:  cfg.NATS.Auth.Password,
			CredsFile: cfg.NATS.Auth.CredsFile,
		},
		TLS: TLSConfig{
			Enabled:    cfg.NATS.TLS.Enabled,
			SkipVerify: cfg.NATS.TLS.SkipVerify,
			CAFile:     cfg.NATS.TLS.CAFile,
			CertFile:   cfg.NATS.TLS.CertFile,
			KeyFile:    cfg.NATS.TLS.KeyFile,
		},
		Middleware: NATSMiddlewareConfig{
			Recovery: MiddlewareState{
				Enabled: cfg.NATS.Middleware.Recovery.Enabled,
			},
			Metrics: MetricsConfig{
				Enabled: cfg.NATS.Middleware.Metrics.Enabled,
				Path:    cfg.NATS.Middleware.Metrics.Path,
			},
			Tracing: MiddlewareState{
				Enabled: cfg.NATS.Middleware.Tracing.Enabled,
			},
			Logging: LoggingConfig{
				Enabled: cfg.NATS.Middleware.Logging.Enabled,
			},
			CircuitBreaker: CircuitBreakerConfig{
				Enabled:       cfg.NATS.Middleware.CircuitBreaker.Enabled,
				MaxRequests:   cfg.NATS.Middleware.CircuitBreaker.MaxRequests,
				Interval:      cfg.NATS.Middleware.CircuitBreaker.Interval,
				Timeout:       cfg.NATS.Middleware.CircuitBreaker.Timeout,
				TripThreshold: cfg.NATS.Middleware.CircuitBreaker.TripThreshold,
			},
			Retry: RetryConfig{
				Enabled:         cfg.NATS.Middleware.Retry.Enabled,
				MaxAttempts:     cfg.NATS.Middleware.Retry.MaxAttempts,
				InitialInterval: cfg.NATS.Middleware.Retry.InitialInterval,
				Multiplier:      cfg.NATS.Middleware.Retry.Multiplier,
				MaxInterval:     cfg.NATS.Middleware.Retry.MaxInterval,
			},
			Timeout: TimeoutConfig{
				Enabled: cfg.NATS.Middleware.Timeout.Enabled,
				Default: cfg.NATS.Middleware.Timeout.Default,
			},
			RateLimit: RateLimitConfig{
				Enabled:           cfg.NATS.Middleware.RateLimit.Enabled,
				RequestsPerSecond: cfg.NATS.Middleware.RateLimit.RequestsPerSecond,
				Burst:             cfg.NATS.Middleware.RateLimit.Burst,
			},
		},
	}

	return m.Init(natsCfg, logger, cfg.App.Name)
}

func (m *Messenger) validateConfig(cfg *config.Config) error {
	if cfg.NATS.URL == "" {
		return fmt.Errorf("NATS URL cannot be empty")
	}
	if cfg.NATS.ConnectionTimeout < 0 {
		return fmt.Errorf("connection timeout cannot be negative")
	}
	if cfg.NATS.DrainTimeout < 0 {
		return fmt.Errorf("drain timeout cannot be negative")
	}
	if cfg.NATS.Middleware.CircuitBreaker.Enabled {
		if cfg.NATS.Middleware.CircuitBreaker.TripThreshold <= 0 {
			return fmt.Errorf("circuit breaker trip threshold must be > 0")
		}
	}
	if cfg.NATS.Middleware.Retry.Enabled {
		if cfg.NATS.Middleware.Retry.MaxAttempts <= 0 {
			return fmt.Errorf("retry max attempts must be > 0")
		}
	}
	if cfg.NATS.Middleware.RateLimit.Enabled {
		if cfg.NATS.Middleware.RateLimit.RequestsPerSecond <= 0 {
			return fmt.Errorf("rate limit requests per second must be > 0")
		}
	}
	return nil
}

// Init initializes the Messenger with configuration, connecting to NATS and setting up pub/sub.
func (m *Messenger) Init(cfg Config, logger *zap.Logger, source string) error {
	client, err := NewNATSClient(cfg, logger)
	if err != nil {
		return fmt.Errorf("failed to create NATS client: %w", err)
	}

	if err := client.Connect(); err != nil {
		_ = client.Close()
		return fmt.Errorf("failed to connect to NATS: %w", err)
	}

	client.source = source
	m.Client = client

	return nil
}

func (m *Messenger) Start() error {

	m.Publisher = NewPublisher(m.Client)
	m.Subscriber = NewSubscriber(m.Client)

	// --- Middleware Registration ---
	// Order matters: Outer -> Inner
	// Current Implementation applies middlewares in order of valid array index
	// Publisher.Publish uses: for i := len(p.middleware) - 1; i >= 0; i--
	// This means the LAST added middleware wraps the PREVIOUS ones.
	// WAIT. If loop is reverse:
	// List: [A, B, C]
	// i=2: C(func)
	// i=1: B(C(func))
	// i=0: A(B(C(func)))
	// So A (first added) is OUTERMOST.
	// So my previous derivation was correct.
	// Recovery is added first -> Outermost. Correct.

	// 1. Recovery (Panic Safety) - Should be outermost
	if m.Client.config.Middleware.Recovery.Enabled {
		m.Publisher.Use(PublisherRecoveryMiddleware(m.Client.logger))
		m.Publisher.UseRequest(RequestRecoveryMiddleware(m.Client.logger))
		m.Publisher.UseJS(JSPublisherRecoveryMiddleware(m.Client.logger))
		m.Publisher.UseAsyncJS(JSAsyncPublisherRecoveryMiddleware(m.Client.logger))
		m.Subscriber.Use(RecoveryMiddleware(m.Client.logger))
		m.Client.logger.Debug("Recovery middleware enabled")
	}

	// 2. Metrics
	if m.Client.config.Middleware.Metrics.Enabled {
		m.Publisher.Use(PublisherMetricsMiddleware())
		m.Publisher.UseRequest(RequestMetricsMiddleware())
		m.Publisher.UseJS(JSPublisherMetricsMiddleware())
		m.Publisher.UseAsyncJS(JSAsyncPublisherMetricsMiddleware())
		m.Subscriber.Use(MetricsMiddleware())
		m.Client.logger.Debug("Metrics middleware enabled")
	}

	// 3. Logging
	if m.Client.config.Middleware.Logging.Enabled {
		m.Publisher.Use(PublisherLoggingMiddleware(m.Client.logger))
		m.Publisher.UseRequest(RequestLoggingMiddleware(m.Client.logger))
		m.Publisher.UseJS(JSPublisherLoggingMiddleware(m.Client.logger))
		m.Publisher.UseAsyncJS(JSAsyncPublisherLoggingMiddleware(m.Client.logger))
		m.Subscriber.Use(LoggingMiddleware(m.Client.logger))
		m.Client.logger.Debug("Logging middleware enabled")
	}

	// 4. Tracing (OpenTelemetry)
	if m.Client.config.Middleware.Tracing.Enabled {
		tracer := otel.Tracer("nats")
		m.Publisher.Use(PublisherTracingMiddleware(tracer))
		m.Publisher.UseRequest(RequestTracingMiddleware(tracer))
		m.Publisher.UseJS(JSPublisherTracingMiddleware(tracer))
		m.Publisher.UseAsyncJS(JSAsyncPublisherTracingMiddleware(tracer))
		m.Subscriber.Use(TracingMiddleware(tracer))
		m.Client.logger.Debug("Tracing middleware enabled")
	}

	// 5. Timeout (Global Operation Timeout)
	if m.Client.config.Middleware.Timeout.Enabled {
		m.Publisher.Use(TimeoutMiddleware(m.Client.config.Middleware.Timeout))
		m.Publisher.UseRequest(RequestTimeoutMiddleware(m.Client.config.Middleware.Timeout))
		m.Publisher.UseJS(JSPublisherTimeoutMiddleware(m.Client.config.Middleware.Timeout))
		m.Publisher.UseAsyncJS(JSAsyncPublisherTimeoutMiddleware(m.Client.config.Middleware.Timeout))
		m.Client.logger.Debug("Timeout middleware enabled")
	}

	// 6. Rate Limit (Traffic Shaping)
	if m.Client.config.Middleware.RateLimit.Enabled {
		limiter := NewRateLimiter(m.Client.config.Middleware.RateLimit)
		m.Publisher.Use(RateLimitMiddleware(limiter))
		m.Publisher.UseRequest(RequestRateLimitMiddleware(limiter))
		m.Publisher.UseJS(JSPublisherRateLimitMiddleware(limiter))
		m.Publisher.UseAsyncJS(JSAsyncPublisherRateLimitMiddleware(limiter))
		m.Client.logger.Debug("Rate Limit middleware enabled")
	}

	// 7. Validator (Schema Check) - Before Logic/CB
	// Validation errors shouldn't trip CB.
	// Since middlewares run Outer->Inner based on registration order (0..N where 0 is Outer),
	// We want Validator to be OUTER to CB.
	// So we register it BEFORE CB?
	// WAIT. Re-verify loop:
	// for i := len-1; i >= 0; i-- { func = mw[i](func) }
	// i=Last: func = Last(Original)
	// i=0:    func = First(Last(Original))
	// So First is OUTERMOST. Last is INNERMOST.
	// To make Validator OUTER to CB, Validator must have LOWER index than CB.
	// So Validator must be registered BEFORE CB.
	// Current Order: RateLimit -> Validator -> CB -> Retry.
	m.Publisher.Use(ValidatorMiddleware(m.Publisher.(*NATSPublisher)))
	m.Publisher.UseRequest(RequestValidatorMiddleware(m.Publisher.(*NATSPublisher)))
	m.Publisher.UseJS(JSPublisherValidatorMiddleware(m.Publisher.(*NATSPublisher)))
	m.Publisher.UseAsyncJS(JSAsyncPublisherValidatorMiddleware(m.Publisher.(*NATSPublisher)))
	m.Subscriber.Use(SubscriberValidatorMiddleware(m.Subscriber.(*NATSSubscriber)))
	m.Client.logger.Debug("Validator middleware enabled")

	// 8. Circuit Breaker (Fail Fast)
	if m.Client.config.Middleware.CircuitBreaker.Enabled {
		cb := NewCircuitBreaker(m.Client.config.Middleware.CircuitBreaker, m.Client.logger)
		m.Publisher.Use(CircuitBreakerMiddleware(cb, m.Client.logger))
		m.Publisher.UseRequest(RequestCircuitBreakerMiddleware(cb))
		m.Publisher.UseJS(JSPublisherCircuitBreakerMiddleware(cb))
		m.Publisher.UseAsyncJS(JSAsyncPublisherCircuitBreakerMiddleware(cb))
		// Subscriber doesn't typically need CB for incoming messages
		m.Client.logger.Debug("Circuit Breaker middleware enabled")
	}

	// 8. Retry (Transient Failures) - Innermost
	if m.Client.config.Middleware.Retry.Enabled {
		m.Publisher.Use(RetryMiddleware(m.Client.config.Middleware.Retry))
		m.Publisher.UseRequest(RequestRetryMiddleware(m.Client.config.Middleware.Retry))
		m.Publisher.UseJS(JSPublisherRetryMiddleware(m.Client.config.Middleware.Retry))
		m.Publisher.UseAsyncJS(JSAsyncPublisherRetryMiddleware(m.Client.config.Middleware.Retry))
		m.Client.logger.Debug("Retry middleware enabled")
	}

	return nil
}

// Publish sends a message to a topic.
func (m *Messenger) Publish(ctx context.Context, subject string, msgType string, data interface{}, opts *PublishOptions) error {
	return m.Publisher.Publish(ctx, subject, msgType, data, opts)
}

// PublishRequest sends a request to a topic and waits for a response.
func (m *Messenger) Request(ctx context.Context, subject string, msgType string, data interface{}, timeout time.Duration, opts *PublishOptions) (*MessageEnvelope, error) {
	return m.Publisher.Request(ctx, subject, msgType, data, timeout, opts)
}

// PublishError sends an error message to a topic.
func (m *Messenger) PublishError(ctx context.Context, subject string, errMsg string) error {
	return m.Publisher.PublishError(ctx, subject, errMsg)
}

// PublishJetStream sends a message to a JetStream topic.
func (m *Messenger) PublishJetStream(ctx context.Context, subject string, msgType string, data interface{}, opts *PublishOptions) (*nats.PubAck, error) {
	return m.Publisher.PublishJS(ctx, subject, msgType, data, opts)
}

func (m *Messenger) PublishAsyncJetStream(ctx context.Context, subject string, msgType string, data interface{}, opts *PublishOptions) (nats.PubAckFuture, error) {
	return m.Publisher.PublishAsyncJS(ctx, subject, msgType, data, opts)
}

// Subscribe registers all topics handled by the service.
func (m *Messenger) Subscribe(ctx context.Context, topic string, handler HandlerFunc, opts *SubscribeOptions) (*nats.Subscription, error) {
	return m.Subscriber.Subscribe(ctx, topic, handler, opts)
}

func (m *Messenger) SubscribePush(ctx context.Context, topic string, handler HandlerFunc, opts *SubscribeOptions) (*nats.Subscription, error) {
	return m.Subscriber.SubscribePush(ctx, topic, handler, opts)
}

func (m *Messenger) SubscribePull(ctx context.Context, subject string, handler HandlerFunc, opts *SubscribeOptions) error {
	return m.Subscriber.SubscribePull(ctx, subject, handler, opts)
}

// Unsubscribe unregisters all topics handled by the service.
func (m *Messenger) Unsubscribe(ctx context.Context) error {
	return m.Subscriber.Unsubscribe(ctx)
}

// Close closes the underlying client and subscriber.
func (m *Messenger) Close() error {
	if m.Subscriber != nil {
		_ = m.Subscriber.Close()
	}
	if m.Client != nil {
		return m.Client.Close()
	}
	return nil
}

// Ping checks the NATS connection health.
func (m *Messenger) Ping() error {
	if m.Client == nil {
		return fmt.Errorf("nats client not initialized")
	}
	return m.Client.Ping()
}
