package natdemo

import (
	"context"

	"grouter/pkg/manager"
	messaging "grouter/pkg/messaging/nats"

	"github.com/go-viper/mapstructure/v2"
	"go.uber.org/zap"
)

func init() {
	manager.RegisterFactory("natdemo", func(deps manager.Deps) (manager.Service, error) {
		return NewNATDemo(deps), nil
	})
}

type NATDemo struct {
	messenger *messaging.Messenger
	logger    *zap.Logger
	config    NATDemoConfig
	//topic registration for nats
	mappings map[string]messaging.HandlerFunc
	metrics  *NATDemoMetrics

	natsSvc   *NATSService
	status    manager.HealthStatus
	lastError error
}

func NewNATDemo(deps manager.Deps) *NATDemo {

	natDemo := &NATDemo{
		messenger: deps.Messenger,
		natsSvc:   NewNATSService(),
		logger:    deps.Logger,
		status:    manager.StatusCreated,
	}

	return natDemo
}

func (e *NATDemo) ID() string {
	return "natdemo"
}

func (e *NATDemo) Name() string {
	return e.config.Name
}

// Dependencies returns the list of dependencies.
// Dependencies returns the list of dependencies.
func (e *NATDemo) Dependencies() []string {
	return nil
}

func (e *NATDemo) Status() manager.HealthStatus {
	return e.status
}

func (e *NATDemo) LastError() error {
	return e.lastError
}

func (e *NATDemo) InitConfig(cfg map[string]interface{}) error {
	e.config = NATDemoConfig{}
	if err := mapstructure.Decode(cfg, &e.config); err != nil {
		e.logger.Error("Failed to decode config", zap.Error(err))
		return err
	}

	return nil
}

// Init initializes the subject to handler mapping.
func (e *NATDemo) Init(ctx context.Context) error {
	e.mappings = make(map[string]messaging.HandlerFunc)
	e.mappings["create"] = e.HandleCreate
	e.mappings["update"] = e.HandleUpdate
	e.mappings["delete"] = e.HandleDelete

	e.metrics = NewNATDemoMetrics(e.messenger.Source(), e.Name())

	e.status = manager.StatusInitialized
	return nil
}

// Start starts the service.
func (e *NATDemo) Start(ctx context.Context) error {
	e.status = manager.StatusStarting
	opts := &messaging.SubscribeOptions{
		QueueGroup: e.config.QueueGroup,
	}
	// Use configured Subject as prefix
	err := e.Subscribe(ctx, opts)
	if err != nil {
		e.status = manager.StatusFailed
		e.lastError = err
		return err
	}
	e.status = manager.StatusRunning
	return nil
}

// Stop stops the service.
func (e *NATDemo) Stop(ctx context.Context) error {
	e.status = manager.StatusStopping
	err := e.Unsubscribe(ctx)
	if err != nil {
		e.lastError = err
		// Do not fail status completely on stop error?
	}
	e.status = manager.StatusStopped
	return err
}

// Subscribe registers all handled topics with the messenger.
func (e *NATDemo) Subscribe(ctx context.Context, opts *messaging.SubscribeOptions) error {
	e.logger.Debug("NATDemo subscribing to topics", zap.String("prefix", e.Name()), zap.String("queueGroup", opts.QueueGroup))

	for subject, handler := range e.mappings {
		fullSubject := e.messenger.Source() + "." + e.Name() + "." + subject

		e.logger.Debug("Registering topic", zap.String("subject", fullSubject))
		_, err := e.messenger.Subscriber.Subscribe(ctx, fullSubject, handler, opts)
		if err != nil {
			return err
		}
	}
	return nil
}

// Unsubscribe unregisters all topics from NATS.
func (e *NATDemo) Unsubscribe(ctx context.Context) error {
	e.logger.Debug("NATDemo unsubscribing from all topics")

	for subject := range e.mappings {
		fullSubject := e.messenger.Source() + "." + e.Name() + "." + subject
		if err := e.messenger.Subscriber.UnsubscribeSubject(ctx, fullSubject); err != nil {
			e.logger.Error("Failed to unsubscribe", zap.Error(err), zap.String("subject", fullSubject))
		}
	}
	return nil
}

func (e *NATDemo) HandleCreate(ctx context.Context, subject string, msg *messaging.MessageEnvelope) error {
	e.logger.Debug("HandleCreate", zap.String("topic", subject), zap.String("id", msg.ID))
	e.metrics.RequestsTotal.WithLabelValues("create").Inc()
	return e.natsSvc.Create(ctx)
}

func (e *NATDemo) HandleUpdate(ctx context.Context, subject string, msg *messaging.MessageEnvelope) error {
	e.logger.Debug("HandleUpdate", zap.String("topic", subject), zap.String("id", msg.ID))
	e.metrics.RequestsTotal.WithLabelValues("update").Inc()
	return nil
}

func (e *NATDemo) HandleDelete(ctx context.Context, subject string, msg *messaging.MessageEnvelope) error {
	e.logger.Debug("HandleDelete", zap.String("topic", subject), zap.String("id", msg.ID))
	e.metrics.RequestsTotal.WithLabelValues("delete").Inc()
	return nil
}
