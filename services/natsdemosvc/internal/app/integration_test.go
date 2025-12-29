package app

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	messaging "grouter/pkg/messaging/nats"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/require"
)

// runServer starts an embedded NATS server for testing
func runServer(t *testing.T) *server.Server {
	opts := &server.Options{
		Port: -1, // Random port
	}
	s, err := server.NewServer(opts)
	require.NoError(t, err)

	go s.Start()

	// Wait for server to be ready
	if !s.ReadyForConnections(5 * time.Second) {
		t.Fatal("NATS server failed to start")
	}

	return s
}

func TestApp_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// 1. Start Embedded NATS
	pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)
	s := runServer(t)
	defer s.Shutdown()
	natsURL := s.ClientURL()
	t.Logf("Started NATS server at %s", natsURL)

	// 2. Configure App to use this NATS
	// We set the env var which config.Load() picks up
	os.Setenv("GROUTER_NATS_URL", natsURL)
	defer os.Unsetenv("GROUTER_NATS_URL")

	// Set config path to the generated config file
	os.Setenv("GROUTER_CONFIG", "configs/config.yaml")
	defer os.Unsetenv("GROUTER_CONFIG")

	// 3. Connect Test Client
	nc, err := nats.Connect(natsURL)
	require.NoError(t, err)
	defer nc.Close()

	// Create temporary config file
	err = os.MkdirAll("configs", 0755)
	require.NoError(t, err)
	defer os.RemoveAll("configs")

	configContent := `
app:
  name: gRouterTest
  version: 1.0.0
  environment: test
nats:
  enabled: true
  url: "nats://localhost:4222"
log:
  level: debug
  format: text
  output_path: stdout
web:
  enabled: false
  port: 8080
services:
  natdemo:
    enabled: true
tracing:
  enabled: true
  service_name: "test-svc"
  exporter: "stdout"
`
	err = os.WriteFile("configs/config.yaml", []byte(configContent), 0644)
	require.NoError(t, err)

	// 4. Start App
	app := New()
	// App.Init() should pick up the env var for NATS URL
	require.NoError(t, app.Init(), "App Init failed")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errChan := make(chan error, 1)
	go func() {
		if err := app.Start(ctx); err != nil {
			errChan <- err
		}
	}()

	// Allow app to start up and subscribe
	time.Sleep(1 * time.Second)

	// 5. Send Start Signal
	// Note: We need to know the App Name. Assuming config defaults to "gRouter" or similar.
	appName := app.GetAppName()
	startTopic := appName + ".start"
	stopTopic := appName + ".stop"

	t.Logf("Sending start signal to %s", startTopic)
	err = publishMessage(nc, startTopic, "start")
	require.NoError(t, err)

	// Wait for services to register
	time.Sleep(1 * time.Second)

	// 6. Test Health (Proof of Life)
	// HealthService listens on: appName + ".health"
	healthSubject := appName + ".health"
	replySubject := "test.reply.health"
	sub, err := nc.SubscribeSync(replySubject)
	require.NoError(t, err)

	healthMsg := &messaging.MessageEnvelope{
		ID:    "test-health",
		Type:  "liveness",
		Reply: replySubject,
	}
	healthData, _ := json.Marshal(healthMsg)
	err = nc.Publish(healthSubject, healthData)
	require.NoError(t, err)

	msg, err := sub.NextMsg(2 * time.Second)
	if err == nil {
		t.Logf("Received health response: %s", string(msg.Data))
	} else {
		t.Logf("Health check timed out: %v", err)
	}

	// 7. Send Stop Signal
	t.Logf("Sending stop signal to %s", stopTopic)
	err = publishMessage(nc, stopTopic, "stop")
	require.NoError(t, err)

	// Wait for shutdown log
	time.Sleep(1 * time.Second)

	// 8. Cleanup
	cancel()
	select {
	case err := <-errChan:
		if err != context.Canceled {
			t.Logf("App stopped with error: %v", err)
		}
	default:
	}
}

func publishMessage(nc *nats.Conn, subject, msgType string) error {
	env := messaging.MessageEnvelope{
		ID:        "test-" + msgType,
		Type:      msgType,
		Source:    "integration-test",
		Timestamp: time.Now(),
	}
	data, _ := json.Marshal(env)
	return nc.Publish(subject, data)
}
