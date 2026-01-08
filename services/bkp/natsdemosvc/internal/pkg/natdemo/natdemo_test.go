package natdemo

import (
	"context"
	"testing"
	"time"

	messaging "grouter/pkg/messaging/nats"

	"github.com/nats-io/nats.go"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"grouter/pkg/manager"
)

type mockPublisher struct {
}

func (m *mockPublisher) Publish(ctx context.Context, subject string, msgType string, data interface{}, opts *messaging.PublishOptions) error {
	return nil
}

func (m *mockPublisher) PublishError(ctx context.Context, subject string, errMsg string) error {
	return nil
}

func (m *mockPublisher) Request(ctx context.Context, subject string, msgType string, data interface{}, timeout time.Duration, opts *messaging.PublishOptions) (*messaging.MessageEnvelope, error) {
	return nil, nil
}

func (m *mockPublisher) PublishJS(ctx context.Context, subject string, msgType string, data interface{}, opts *messaging.PublishOptions) (*nats.PubAck, error) {
	return nil, nil
}

func (m *mockPublisher) PublishAsyncJS(ctx context.Context, subject string, msgType string, data interface{}, opts *messaging.PublishOptions) (nats.PubAckFuture, error) {
	return nil, nil
}

func (m *mockPublisher) Use(mw ...messaging.PublisherMiddleware)               {}
func (m *mockPublisher) UseRequest(mw ...messaging.RequestMiddleware)          {}
func (m *mockPublisher) UseJS(mw ...messaging.JSPublisherMiddleware)           {}
func (m *mockPublisher) UseAsyncJS(mw ...messaging.JSAsyncPublisherMiddleware) {}
func (m *mockPublisher) SetValidator(v messaging.Validator)                    {}

func TestNATDemo_New(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	messenger := &messaging.Messenger{
		Publisher: &mockPublisher{},
	}
	cfg := NATDemoConfig{Enabled: true, Name: "natdemo"}

	deps := manager.Deps{
		Logger:    logger,
		Messenger: messenger,
	}
	demo := NewNATDemo(deps)
	demo.config = cfg
	assert.NotNil(t, demo)
	assert.Equal(t, "natdemo", demo.Name())
}

func TestNATDemo_Lifecycle(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	messenger := &messaging.Messenger{
		Publisher: &mockPublisher{},
	}
	cfg := NATDemoConfig{Enabled: true, Name: "natdemo"}

	deps := manager.Deps{
		Logger:    logger,
		Messenger: messenger,
	}
	demo := NewNATDemo(deps)
	demo.config = cfg

	err := demo.Init(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, demo.mappings)
	assert.Contains(t, demo.mappings, "create")
}

func TestNATDemo_Handle(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	messenger := &messaging.Messenger{
		Publisher: &mockPublisher{},
	}
	cfg := NATDemoConfig{Enabled: true, Name: "natdemo"}

	deps := manager.Deps{
		Logger:    logger,
		Messenger: messenger,
	}
	demo := NewNATDemo(deps)
	demo.config = cfg
	_ = demo.Init(context.Background())
	ctx := context.Background()

	// Test natdemo.create
	env := &messaging.MessageEnvelope{
		Type: "natdemo.create",
	}
	err := demo.HandleCreate(ctx, "natdemo.create", env)
	assert.NoError(t, err)

	// Test through mappings (more realistic)
	handler, ok := demo.mappings["create"]
	assert.True(t, ok)
	if ok {
		err = handler(ctx, "natdemo.create", env)
		assert.NoError(t, err)
	}
	assert.NoError(t, err)
}

func TestNATDemo_Metrics(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	messenger := &messaging.Messenger{
		Publisher: &mockPublisher{},
	}
	cfg := NATDemoConfig{Enabled: true, Name: "natdemo"}

	deps := manager.Deps{
		Logger:    logger,
		Messenger: messenger,
	}
	demo := NewNATDemo(deps)
	demo.config = cfg

	err := demo.Init(context.Background())
	assert.NoError(t, err)

	assert.NotNil(t, demo.metrics, "metrics should be initialized after Init")
	assert.NotNil(t, demo.metrics.RequestsTotal, "requests total counter should be initialized")

	// Ensure no panic when incrementing
	assert.NotPanics(t, func() {
		demo.metrics.RequestsTotal.WithLabelValues("test").Inc()
	})

	// Verify value using testutil
	// We need to pass the collector that corresponds to the specific label
	metric := demo.metrics.RequestsTotal.WithLabelValues("test")
	val := testutil.ToFloat64(metric)
	assert.Equal(t, 1.0, val, "metric value using testutil should be 1.0")
}
