package app

import (
	"context"
	"testing"
	"time"

	messaging "grouter/pkg/messaging/nats"

	"github.com/stretchr/testify/assert"
)

func TestBootstrapService_Handle(t *testing.T) {
	trigger := make(chan struct{}, 1)
	svc := NewBootstrapService(trigger)

	ctx := context.Background()
	env := &messaging.MessageEnvelope{}

	err := svc.Handle(ctx, "start", env)
	assert.NoError(t, err)

	select {
	case <-trigger:
		// Success
	case <-time.After(1 * time.Second):
		t.Fatal("Trigger channel was not signaled")
	}
}

func TestBootstrapService_Name(t *testing.T) {
	svc := NewBootstrapService(nil)
	assert.Equal(t, "start", svc.Name())
}
