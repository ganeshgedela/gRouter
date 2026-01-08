package app

import (
	"context"
	"testing"
	"time"

	messaging "grouter/pkg/messaging/nats"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestBootstrapService_Handle(t *testing.T) {
	logger := zap.NewNop()
	svc := NewBootstrapService(nil, logger, "test-app")

	ctx := context.Background()
	env := &messaging.MessageEnvelope{}

	err := svc.HandleStart(ctx, "start", env)
	assert.NoError(t, err)

	select {
	case <-svc.WaitForStart():
		// Success
	case <-time.After(1 * time.Second):
		t.Fatal("Start channel was not signaled")
	}
}

func TestBootstrapService_Name(t *testing.T) {
	svc := NewBootstrapService(nil, zap.NewNop(), "test-app")
	assert.Equal(t, "bootstrap", svc.Name())
}
