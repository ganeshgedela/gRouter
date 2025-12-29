package nats

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestMapValidator(t *testing.T) {
	v := NewMapValidator()

	// Register a validator for "test.type"
	v.Register("test.type", func(data []byte) error {
		var m map[string]string
		if err := json.Unmarshal(data, &m); err != nil {
			return err
		}
		if _, ok := m["required"]; !ok {
			return fmt.Errorf("missing required field")
		}
		return nil
	})

	t.Run("Valid data", func(t *testing.T) {
		data := []byte(`{"required": "value"}`)
		err := v.Validate("test.type", data)
		assert.NoError(t, err)
	})

	t.Run("Invalid data", func(t *testing.T) {
		data := []byte(`{"wrong": "value"}`)
		err := v.Validate("test.type", data)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "missing required field")
	})

	t.Run("Unknown type", func(t *testing.T) {
		data := []byte(`{}`)
		err := v.Validate("unknown.type", data)
		assert.NoError(t, err) // Default behavior is to allow unknown types
	})
}

func TestPublisher_Validation(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	client, _ := NewNATSClient(Config{URL: "nats://localhost:4222"}, logger)
	// Note: We don't need to connect for this unit test as we check validation before connection check in Publish

	pub := NewPublisher(client, "test-source")
	v := NewMapValidator()
	v.Register("test.type", func(data []byte) error {
		return fmt.Errorf("validation failed")
	})
	pub.SetValidator(v)

	err := pub.Publish(context.Background(), "subject", "test.type", map[string]string{"foo": "bar"}, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation failed")
}
