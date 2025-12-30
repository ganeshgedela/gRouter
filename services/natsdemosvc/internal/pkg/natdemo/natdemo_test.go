package natdemo

import (
	"context"
	"testing"
	"time"

	messaging "grouter/pkg/messaging/nats"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type mockPublisher struct {
}

func (m *mockPublisher) Publish(ctx context.Context, subject string, msgType string, data interface{}, opts *messaging.PublishOptions) error {
	return nil
}

func (m *mockPublisher) PublishError(ctx context.Context, subject string, errMsg string) error {
	return nil
}

func (m *mockPublisher) Request(ctx context.Context, subject string, msgType string, data interface{}, timeout time.Duration) (*messaging.MessageEnvelope, error) {
	return nil, nil
}

func (m *mockPublisher) PublishJS(ctx context.Context, subject string, msgType string, data interface{}, opts ...nats.PubOpt) (*nats.PubAck, error) {
	return nil, nil
}

func (m *mockPublisher) PublishAsyncJS(ctx context.Context, subject string, msgType string, data interface{}, opts ...nats.PubOpt) (nats.PubAckFuture, error) {
	return nil, nil
}

func (m *mockPublisher) Use(mw ...messaging.PublisherMiddleware)      {}
func (m *mockPublisher) UseRequest(mw ...messaging.RequestMiddleware) {}
func (m *mockPublisher) SetValidator(v messaging.Validator)           {}

func TestNATDemo_New(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	pub := &mockPublisher{}
	cfg := NATDemoConfig{Enabled: true}

	demo := NewNATDemo(pub, logger, cfg)
	assert.NotNil(t, demo)
	assert.Equal(t, "natdemo", demo.Name())
}

func TestNATDemo_Lifecycle(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	pub := &mockPublisher{}
	cfg := NATDemoConfig{Enabled: true}
	demo := NewNATDemo(pub, logger, cfg)
	ctx := context.Background()

	assert.NoError(t, demo.Ready(ctx))
	assert.NoError(t, demo.Start(ctx))
	assert.NoError(t, demo.Stop(ctx))
}

func TestNATDemo_Handle(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	pub := &mockPublisher{}
	cfg := NATDemoConfig{Enabled: true}
	demo := NewNATDemo(pub, logger, cfg)
	ctx := context.Background()

	// Test natdemo.create
	env := &messaging.MessageEnvelope{
		Type: "natdemo.create",
	}
	err := demo.Handle(ctx, "topic", env)
	assert.NoError(t, err)

	// Test unknown
	env2 := &messaging.MessageEnvelope{
		Type: "unknown",
	}
	err = demo.Handle(ctx, "topic", env2)
	assert.NoError(t, err)
}
