package app

import (
	"context"
	"testing"
	"time"

	messaging "grouter/pkg/messaging/nats"

	"github.com/stretchr/testify/assert"
)

func TestStopService_Handle(t *testing.T) {
	trigger := make(chan struct{}, 1)
	svc := NewStopService(trigger)

	ctx := context.Background()
	env := &messaging.MessageEnvelope{}

	err := svc.Handle(ctx, "stop", env)
	assert.NoError(t, err)

	select {
	case <-trigger:
		// Success
	case <-time.After(1 * time.Second):
		t.Fatal("Trigger channel was not signaled")
	}
}

func TestStopService_Name(t *testing.T) {
	svc := NewStopService(nil)
	assert.Equal(t, "stop", svc.Name())
}
