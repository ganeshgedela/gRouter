package nats

import (
	"fmt"

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

// NewMessenger creates a new Messenger.
func NewMessenger(client *Client, pub Publisher, sub Subscriber) *Messenger {
	return &Messenger{
		Client:     client,
		Publisher:  pub,
		Subscriber: sub,
	}
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

	m.Client = client
	m.Publisher = NewPublisher(client, source)
	m.Subscriber = NewSubscriber(client, source)

	// Enable metrics middleware if configured
	if cfg.Metrics.Enabled {
		m.Publisher.Use(PublisherMetricsMiddleware())
		m.Subscriber.Use(MetricsMiddleware())
		logger.Info("Metrics middleware enabled for NATS")
	}

	return nil
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
