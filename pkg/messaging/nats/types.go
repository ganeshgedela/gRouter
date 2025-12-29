package nats

import (
	"context"
	"encoding/json"
	"time"

	"github.com/nats-io/nats.go"
)

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
	Request(ctx context.Context, subject string, msgType string, data interface{}, timeout time.Duration) (*MessageEnvelope, error)
	// JetStream methods
	PublishJS(ctx context.Context, subject string, msgType string, data interface{}, opts ...nats.PubOpt) (*nats.PubAck, error)
	PublishAsyncJS(ctx context.Context, subject string, msgType string, data interface{}, opts ...nats.PubOpt) (nats.PubAckFuture, error)
	Use(mw ...PublisherMiddleware)
	SetValidator(v Validator)
}

// PublishOptions configures message publishing behavior.
type PublishOptions struct {
	// Async determines if the publish should be asynchronous.
	// If false, the publisher will flush the connection to ensure the message is sent.
	Async bool
	// Timeout specifies how long to wait for a response in request-response patterns.
	Timeout time.Duration
}

// SubscribeOptions configures message subscription behavior.
type SubscribeOptions struct {
	// QueueGroup enables load balancing between multiple instances of a service.
	// When set, NATS will deliver each message to only one member of the group.
	QueueGroup string
	// MaxWorkers specifies the maximum number of concurrent workers for processing messages.
	MaxWorkers int
}

// PublisherMiddleware defines the middleware for publishing messages.
type PublisherMiddleware func(next PublisherFunc) PublisherFunc

// PublisherFunc is the function signature for publishing messages.
type PublisherFunc func(ctx context.Context, subject string, msgType string, data interface{}, opts *PublishOptions) error

// SubscriberMiddleware defines the middleware for subscribing to messages.
type SubscriberMiddleware func(next HandlerFunc) HandlerFunc

// Subscriber defines the interface for subscribing to messages.
type Subscriber interface {
	Subscribe(subject string, handler HandlerFunc, opts *SubscribeOptions) error
	SubscribePush(subject string, handler HandlerFunc, opts ...nats.SubOpt) error
	SubscribePull(subject, durable string, handler HandlerFunc, opts ...PullOption) error
	Unsubscribe() error
	Close() error

	Use(mw ...SubscriberMiddleware)
	SetValidator(v Validator)
}

// PullOptions configures behavior for pull consumers.
type PullOptions struct {
	BatchSize    int
	FetchTimeout time.Duration
}

// PullOption is a functional option for configuring pull consumers.
type PullOption func(*PullOptions)

// WithBatchSize sets the number of messages to fetch in each batch.
func WithBatchSize(size int) PullOption {
	return func(o *PullOptions) {
		o.BatchSize = size
	}
}

// WithFetchTimeout sets the maximum time to wait for a batch of messages.
func WithFetchTimeout(timeout time.Duration) PullOption {
	return func(o *PullOptions) {
		o.FetchTimeout = timeout
	}
}
