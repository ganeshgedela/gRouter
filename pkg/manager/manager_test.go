package manager

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"grouter/pkg/config"
	messaging "grouter/pkg/messaging/nats"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type mockPublisher struct {
	publishedSubject string
	publishedType    string
	publishedData    interface{}
}

func (m *mockPublisher) Publish(ctx context.Context, subject string, msgType string, data interface{}, opts *messaging.PublishOptions) error {
	m.publishedSubject = subject
	m.publishedType = msgType
	m.publishedData = data
	return nil
}

func (m *mockPublisher) PublishError(ctx context.Context, subject string, errMsg string) error {
	m.publishedSubject = subject
	m.publishedType = "error"
	m.publishedData = map[string]string{"error": errMsg}
	return nil
}

func (m *mockPublisher) Request(ctx context.Context, subject string, msgType string, data interface{}, timeout time.Duration) (*messaging.MessageEnvelope, error) {
	return nil, nil
}

func (m *mockPublisher) PublishJS(ctx context.Context, subject string, msgType string, data interface{}, opts ...nats.PubOpt) (*nats.PubAck, error) {
	m.publishedSubject = subject
	m.publishedType = msgType
	m.publishedData = data
	return &nats.PubAck{}, nil
}

func (m *mockPublisher) PublishAsyncJS(ctx context.Context, subject string, msgType string, data interface{}, opts ...nats.PubOpt) (nats.PubAckFuture, error) {
	m.publishedSubject = subject
	m.publishedType = msgType
	m.publishedData = data
	return nil, nil
}

func (m *mockPublisher) Use(mw ...messaging.PublisherMiddleware) {
	// no-op for mock
}

func (m *mockPublisher) UseRequest(mw ...messaging.RequestMiddleware) {
	// no-op for mock
}

func (m *mockPublisher) SetValidator(v messaging.Validator) {
	// no-op for mock
}

func TestServiceManager_OnMessage(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	router := NewServiceRouter()

	svc := &mockService{name: "test"}
	router.Register("test", svc)

	pub := &mockPublisher{}
	// Create mock messenger
	messenger := &messaging.Messenger{
		Publisher: pub,
	}

	mgr := &ServiceManager{
		log:       logger,
		router:    router,
		messenger: messenger,
		timeout:   1 * time.Second,
		cfg: &config.Config{
			App: config.AppConfig{Name: "grouter"},
		},
	}

	ctx := context.Background()

	t.Run("Successful routing", func(t *testing.T) {
		env := &messaging.MessageEnvelope{
			ID:     "123",
			Type:   "test.op",
			Source: "client",
			Data:   json.RawMessage(`{"foo":"bar"}`),
		}

		err := mgr.onNATSMessage(ctx, "grouter.test.op", env)
		assert.NoError(t, err)
	})

	t.Run("Routing failure", func(t *testing.T) {
		env := &messaging.MessageEnvelope{
			ID:     "456",
			Type:   "unknown.op",
			Source: "client",
			Data:   json.RawMessage(`{}`),
		}

		err := mgr.onNATSMessage(ctx, "grouter.unknown.op", env)
		assert.NoError(t, err)
		// Should have published an error reply if Reply was set, but it wasn't
	})

	t.Run("Message with reply", func(t *testing.T) {
		env := &messaging.MessageEnvelope{
			ID:    "789",
			Type:  "test.op",
			Reply: "inbox.123",
			Data:  json.RawMessage(`{"foo":"bar"}`),
		}

		err := mgr.onNATSMessage(ctx, "grouter.test.op", env)
		assert.NoError(t, err)
		// Auto-reply is removed from manager, so we don't expect a publish here from the manager itself.
		// The Service is responsible for replying.
	})

	t.Run("Routing error replies", func(t *testing.T) {
		// Mock a service that returns an error
		errSvc := &errorService{mockService{name: "error"}}
		router.Register("error", errSvc)

		env := &messaging.MessageEnvelope{
			ID:     "999",
			Type:   "error.op",
			Source: "client",
			Reply:  "inbox.error",
			Data:   json.RawMessage(`{}`),
		}

		err := mgr.onNATSMessage(ctx, "grouter.error.op", env)
		assert.NoError(t, err)

		// Verify that PublishError was called on the mock publisher
		assert.Equal(t, "inbox.error", pub.publishedSubject)
		assert.Equal(t, "error", pub.publishedType)
		// Check data if needed
		dataMap, ok := pub.publishedData.(map[string]string)
		assert.True(t, ok)
		assert.Equal(t, "intentional error", dataMap["error"])
	})
}

func TestServiceManager_Stop(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	mgr := &ServiceManager{
		log: logger,
	}

	err := mgr.Stop(context.Background())
	assert.NoError(t, err)
}

type errorService struct {
	mockService
}

func (s *errorService) Handle(ctx context.Context, topic string, msg *messaging.MessageEnvelope) error {
	return fmt.Errorf("intentional error")
}
