package app

import (
	"bytes"
	"context"
	"fmt"
	"sync"

	"grouter/pkg/manager"
	messaging "grouter/pkg/messaging/nats"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/expfmt"
	"go.uber.org/zap"
)

// MetricService exposes Prometheus metrics via NATS.
type MetricService struct {
	messenger *messaging.Messenger
	logger    *zap.Logger
	mu        sync.RWMutex

	Collector *MetricsCollector

	status    manager.HealthStatus
	lastError error
}

// init registers the MetricService factory.
func init() {
	manager.RegisterFactory("metrics", func(deps manager.Deps) (manager.Service, error) {
		return NewMetricService(deps), nil
	})
}

// NewMetricService creates a new MetricService.
func NewMetricService(deps manager.Deps) *MetricService {
	// Note: We need manager instance for collector?
	// NewMetricsCollector(mgr) depends on manager to list services.
	// But deps doesn't have manager (avoid cycle).
	// However, Manager passes itself? No.
	// We can't pass Manager into Deps because Manager depends on Deps.
	// MetricsService needs Manager to ListServices.
	// This suggests MetricsService is tightly coupled to Manager.
	// If so, we might need a workaround or keep it manual?
	// "Standardizing" means decoupling.
	// Maybe MetricsCollector shouldn't need Manager?
	// It uses Manager.ListServices().
	// Can we use ServiceStore from Deps?
	// Deps has *ServiceStore!
	// Yes! We should update NewMetricsCollector to use ServiceStore instead of Manager.
	// init registers the MetricService factory.
func init() {
	manager.RegisterFactory("metrics", func(deps manager.Deps) (manager.Service, error) {
		return NewMetricService(deps), nil
	})
}

// NewMetricService creates a new MetricService.
func NewMetricService(deps manager.Deps) *MetricService {
	collector := NewMetricsCollector(deps.Store)
	prometheus.MustRegister(collector)

	return &MetricService{
		messenger: deps.Messenger,
		logger:    deps.Logger,
		Collector: collector,
	}
}

// Name returns the service name.
func (s *MetricService) Name() string {
	return "metrics"
}

func (s *MetricService) ID() string {
	return "metrics"
}

func (s *MetricService) Status() manager.HealthStatus {
	return s.status
}

func (s *MetricService) LastError() error {
	return s.lastError
}

// Dependencies returns the list of dependencies.
func (s *MetricService) Dependencies() []string {
	return nil
}

// Init initializes the service.
func (s *MetricService) Init(ctx context.Context) error {
	s.logger.Debug("Initializing Metrics Service")
	return nil
}

// Start starts the service and subscribes to the metrics subject.
func (s *MetricService) Start(ctx context.Context) error {
	s.logger.Debug("Starting Metrics Service")

	return s.Subscribe(ctx, &messaging.SubscribeOptions{
		QueueGroup: "metrics-readers",
	})
}

// Stop stops the service.
func (s *MetricService) Stop(ctx context.Context) error {
	s.logger.Debug("Stopping Metrics Service")
	return s.Unsubscribe(ctx)
}

// Subscribe subscribes to a NATS subject.
func (s *MetricService) Subscribe(ctx context.Context, opts *messaging.SubscribeOptions) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	topic := s.messenger.Source() + "." + s.Name()

	s.logger.Debug("Subscribing to metrics topic", zap.String("topic", topic))

	// Define handler
	handler := func(ctx context.Context, subject string, msg *messaging.MessageEnvelope) error {
		return s.HandleMetrics(ctx, msg)
	}

	_, err := s.messenger.Subscribe(ctx, topic, handler, opts)
	if err != nil {
		return fmt.Errorf("failed to subscribe to %s: %w", topic, err)
	}
	return nil
}

// Unsubscribe unsubscribes from all topics.
func (s *MetricService) Unsubscribe(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	subject := s.messenger.Source() + s.Name()
	err := s.messenger.Subscriber.UnsubscribeSubject(ctx, subject)
	if err != nil {
		return fmt.Errorf("failed to unsubscribe from %s: %w", subject, err)
	}
	return nil
}

// HandleMetrics gathers metrics and replies with the result.
func (s *MetricService) HandleMetrics(ctx context.Context, msg *messaging.MessageEnvelope) error {
	gatherer := prometheus.DefaultGatherer
	mfs, err := gatherer.Gather()
	if err != nil {
		s.logger.Error("Failed to gather metrics", zap.Error(err))
		return err
	}

	var buf bytes.Buffer
	enc := expfmt.NewEncoder(&buf, expfmt.FmtText)
	for _, mf := range mfs {
		if err := enc.Encode(mf); err != nil {
			s.logger.Error("Failed to encode metric", zap.Error(err))
			// Continue encoding others
		}
	}

	// Publish response
	if msg.Reply == "" {
		s.logger.Warn("Received metrics request with no Reply subject")
		return nil
	}

	// Create a map for the response data
	replyData := map[string]string{
		"metrics": buf.String(),
	}

	// Messenger.Publish takes data interface{}, it defines serialization (usually JSON).
	// We pass the map directly.
	return s.messenger.Publish(ctx, msg.Reply, "metrics.response", replyData, nil)
}

// MetricsCollector collects metrics about registered services and executes registered callbacks.
type MetricsCollector struct {
	store       *manager.ServiceStore
	serviceDesc *prometheus.Desc
	countDesc   *prometheus.Desc

	mu        sync.RWMutex
	callbacks map[string]func(chan<- prometheus.Metric)
}

// NewMetricsCollector creates a new MetricsCollector.
func NewMetricsCollector(store *manager.ServiceStore) *MetricsCollector {
	return &MetricsCollector{
		store: store,
		serviceDesc: prometheus.NewDesc(
			"grouter_service_info",
			"Information about registered services",
			[]string{"name"},
			nil,
		),
		countDesc: prometheus.NewDesc(
			"grouter_services_registered_total",
			"Total number of registered services",
			nil,
			nil,
		),
		callbacks: make(map[string]func(chan<- prometheus.Metric)),
	}
}

// Register registers a metrics callback for a service.
func (c *MetricsCollector) Register(name string, cb func(chan<- prometheus.Metric)) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.callbacks[name] = cb
}

// Describe implements prometheus.Collector.
func (c *MetricsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.serviceDesc
	ch <- c.countDesc
}

// Collect implements prometheus.Collector.
func (c *MetricsCollector) Collect(ch chan<- prometheus.Metric) {
	services := c.store.List()

	ch <- prometheus.MustNewConstMetric(
		c.countDesc,
		prometheus.GaugeValue,
		float64(len(services)),
	)

	for _, name := range services {
		ch <- prometheus.MustNewConstMetric(
			c.serviceDesc,
			prometheus.GaugeValue,
			1.0,
			name,
		)
	}

	// Execute callbacks
	c.mu.RLock()
	defer c.mu.RUnlock()
	for _, cb := range c.callbacks {
		cb(ch)
	}
}
