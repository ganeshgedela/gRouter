package manager

import (
	"context"
	"testing"

	messaging "grouter/pkg/messaging/nats"

	"github.com/stretchr/testify/assert"
)

type mockService struct {
	name string
}

func (m *mockService) Name() string                    { return m.name }
func (m *mockService) Ready(ctx context.Context) error { return nil }
func (m *mockService) Start(ctx context.Context) error { return nil }
func (m *mockService) Stop(ctx context.Context) error  { return nil }
func (m *mockService) Handle(ctx context.Context, topic string, msg *messaging.MessageEnvelope) error {
	return nil
}

func TestServiceStore(t *testing.T) {
	store := NewServiceStore()
	svc := &mockService{name: "test-svc"}

	// Test Add and Get
	store.Add("test-svc", svc)
	retrieved, ok := store.Get("test-svc")
	assert.True(t, ok)
	assert.Equal(t, svc, retrieved)

	// Test normalization
	retrieved, ok = store.Get("  TEST-SVC  ")
	assert.True(t, ok)
	assert.Equal(t, svc, retrieved)

	// Test Exists
	assert.True(t, store.Exists("test-svc"))

	// Test List
	assert.ElementsMatch(t, []string{"test-svc"}, store.List())

	// Test Delete
	assert.True(t, store.Delete("test-svc"))
	assert.False(t, store.Exists("test-svc"))
}

func TestServiceRouter(t *testing.T) {
	router := NewServiceRouter()
	svc := &mockService{name: "test-svc"}
	router.Register("test-svc", svc)

	tests := []struct {
		topic    string
		expected string
		wantErr  bool
	}{
		{"test-svc.create", "test-svc", false},
		{"test-svc.delete", "test-svc", false},
		{"other.op", "", true},
		{"invalid", "", true},
		{"", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.topic, func(t *testing.T) {
			s, err := router.RouteByTopic(tt.topic)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, s.Name())
			}
		})
	}
}
