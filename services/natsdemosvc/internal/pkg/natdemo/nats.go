package natdemo

import (
	"context"

	messaging "grouter/pkg/messaging/nats"

	"go.uber.org/zap"
)

type NATDemo struct {
	publisher messaging.Publisher
	natsSvc   *NATSService
	logger    *zap.Logger
	config    NATDemoConfig
}

func NewNATDemo(pub messaging.Publisher, logger *zap.Logger, config NATDemoConfig) *NATDemo {
	return &NATDemo{publisher: pub, natsSvc: NewNATSService(), logger: logger, config: config}
}

func (e *NATDemo) Name() string {
	return e.natsSvc.Name()
}

func (e *NATDemo) Ready(ctx context.Context) error {
	return nil
}

func (e *NATDemo) Start(ctx context.Context) error {
	return nil
}

func (e *NATDemo) Stop(ctx context.Context) error {
	return nil
}

func (e *NATDemo) Handle(ctx context.Context, topic string, msg *messaging.MessageEnvelope) error {

	topic = msg.Type

	switch topic {
	case "natdemo.create":
		e.logger.Info("Creating NATS")
		return e.natsSvc.Create(ctx)
	default:
		e.logger.Info("Unknown topic", zap.String("topic", topic))
		return nil
	}
}
