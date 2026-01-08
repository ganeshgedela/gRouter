package nats

import (
	"context"
	"encoding/json"
	"time"

	"github.com/nats-io/nats.go"
)

type NATService interface {
	// Init initializes the service.
	Init(ctx context.Context) error
	Start(ctx context.Context) error
	Stop(ctx context.Context) error

	// Subscribe registers all topics handled by the service.
	Subscribe(ctx context.Context, opts *SubscribeOptions) error
	// Unsubscribe unregisters all topics handled by the service.
	Unsubscribe(ctx context.Context) error
}

// topic = <service_manager_identity>.<service>.<operation>

// MessageEnvelope wraps all messages with metadata. It implements the Envelope Pattern,
// providing a consistent structure for all messages while allowing for deferred
// parsing of the actual payload.
type MessageEnvelope struct {
	// ID is a unique identifier for the message, used for tracking and deduplication.
	ID string `json:"id"`
	// Type identifies the kind of message payload, used for type discrimination.
	Type string `json:"type"`
	// Timestamp is the time when the message was generated.
	Timestamp time.Time `json:"timestamp"`
	// Source identifies the service or component that generated the message.
	Source string `json:"source"`
	// Reply is an optional subject where responses should be sent.
	Reply string `json:"reply,omitempty"`
	// Data is the raw message payload, to be unmarshaled based on the Type.
	Data json.RawMessage `json:"data"`
	// Metadata contains optional key-value pairs for tracing, routing, or other purposes.
	Metadata map[string]string `json:"metadata,omitempty"`
}

// HandlerFunc is the function signature for message handlers
type HandlerFunc func(ctx context.Context, subject string, msg *MessageEnvelope) error

// Validator defines the interface for message schema validation.
type Validator interface {
	// Validate checks if the data matches the schema for the given message type.
	Validate(msgType string, data []byte) error
}

// Publisher defines the interface for publishing messages.
type Publisher interface {
	Publish(ctx context.Context, subject string, msgType string, data interface{}, opts *PublishOptions) error
	PublishError(ctx context.Context, subject string, errMsg string) error
	Request(ctx context.Context, subject string, msgType string, data interface{}, timeout time.Duration, opts *PublishOptions) (*MessageEnvelope, error)

	// JetStream methods
	PublishJS(ctx context.Context, subject string, msgType string, data interface{}, opts *PublishOptions) (*nats.PubAck, error)
	PublishAsyncJS(ctx context.Context, subject string, msgType string, data interface{}, opts *PublishOptions) (nats.PubAckFuture, error)

	Use(mw ...PublisherMiddleware)
	UseRequest(mw ...RequestMiddleware)
	UseJS(mw ...JSPublisherMiddleware)
	UseAsyncJS(mw ...JSAsyncPublisherMiddleware)
	SetValidator(v Validator)
}

// PublishOptions configures message publishing behavior.
type PublishOptions struct {
	// Async determines if the publish should be asynchronous.
	// If false, the publisher will flush the connection to ensure the message is sent.
	Async bool
	// Timeout specifies how long to wait for a response in request-response patterns.
	Timeout time.Duration
	// Metadata contains optional key-value pairs for tracing, routing, or other purposes.
	Metadata map[string]string
	// NatsOptions provides direct access to low-level NATS JS options
	NatsOptions []nats.PubOpt
	// MsgID is an optional unique identifier for JetStream deduplication.
	// If set, NATS will use this to prevent processing duplicate messages within the deduplication window.
	MsgID string
	// Retries specifies the number of times to retry a failed publish operation.
	// Deprecated: Internal retries are removed. Use RetryMiddleware instead.
	Retries int
	// RetryWait is the time to wait between retries.
	// Deprecated: Internal retries are removed. Use RetryMiddleware instead.
	RetryWait time.Duration
}

// SubscribeOptions configures behavior for subscriptions.
type SubscribeOptions struct {
	// QueueGroup enables load balancing between multiple instances of a service.
	// When set, NATS will deliver each message to only one member of the group.
	QueueGroup string
	// MaxWorkers specifies the maximum number of concurrent workers for processing messages.
	MaxWorkers int
	// ServiceName is the name of the service for metrics
	ServiceName string
	// Metadata contains optional key-value pairs for tracing, routing, or other purposes.
	Metadata map[string]string
	// Durable is the durable consumer name. If set, the consumer survives disconnects.
	Durable string
	// AckWait is the time to wait for acknowledgement before redelivery.
	AckWait time.Duration
	// MaxDeliveries is the maximum number of times a message will be delivered.
	MaxDeliveries int
	// FilterSubject allows filtering messages on the stream (wildcards supported).
	FilterSubject string
	// DLQSubject is the subject to send messages to after MaxDeliveries is reached.
	DLQSubject string
	// BatchSize specifies the number of messages to fetch in each batch (Pull Consumers only).
	BatchSize int
	// FetchTimeout specifies the maximum time to wait for a batch of messages (Pull Consumers only).
	FetchTimeout time.Duration
	// NatsOptions provides direct access to low-level NATS JS options
	NatsOptions []nats.SubOpt
}

// PublisherMiddleware defines the middleware for publishing messages.
type PublisherMiddleware func(next PublisherFunc) PublisherFunc

// PublisherFunc is the function signature for publishing messages.
type PublisherFunc func(ctx context.Context, subject string, msgType string, data interface{}, opts *PublishOptions) error

// RequestMiddleware defines the middleware for request-reply messaging.
type RequestMiddleware func(next RequestFunc) RequestFunc

// RequestFunc is the function signature for request-reply messaging.
type RequestFunc func(ctx context.Context, subject string, msgType string, data interface{}, timeout time.Duration, opts *PublishOptions) (*MessageEnvelope, error)

// JSPublisherMiddleware defines the middleware for JetStream publishing.
type JSPublisherMiddleware func(next JSPublisherFunc) JSPublisherFunc

// JSPublisherFunc is the function signature for JetStream publishing.
type JSPublisherFunc func(ctx context.Context, subject string, msgType string, data interface{}, opts *PublishOptions) (*nats.PubAck, error)

// JSAsyncPublisherMiddleware defines the middleware for asynchronous JetStream publishing.
type JSAsyncPublisherMiddleware func(next JSAsyncPublisherFunc) JSAsyncPublisherFunc

// JSAsyncPublisherFunc is the function signature for asynchronous JetStream publishing.
type JSAsyncPublisherFunc func(ctx context.Context, subject string, msgType string, data interface{}, opts *PublishOptions) (nats.PubAckFuture, error)

// SubscriberMiddleware defines the middleware for subscribing to messages.
type SubscriberMiddleware func(next HandlerFunc) HandlerFunc

// Subscriber defines the interface for subscribing to messages.
type Subscriber interface {
	//Subscribe subscribes to a topic.
	Subscribe(ctx context.Context, subject string, handler HandlerFunc, opts *SubscribeOptions) (*nats.Subscription, error)

	//SubscribePush subscribes to a topic using push-based delivery.
	SubscribePush(ctx context.Context, subject string, handler HandlerFunc, opts *SubscribeOptions) (*nats.Subscription, error)
	//SubscribePull subscribes to a topic using pull-based delivery.
	SubscribePull(ctx context.Context, subject string, handler HandlerFunc, opts *SubscribeOptions) error

	//Unsubscribe unsubscribes from a topic.
	Unsubscribe(ctx context.Context) error

	//UnsubscribeSubject unsubscribes from a specific subject.
	UnsubscribeSubject(ctx context.Context, subject string) error

	//Close closes the subscriber.
	Close() error

	//Use registers middleware for the subscri	ber.
	Use(mw ...SubscriberMiddleware)
	//SetValidator sets the validator for the subscriber.
	SetValidator(v Validator)
}
