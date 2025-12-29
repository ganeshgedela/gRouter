package app

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApp_Integration_WebDemo(t *testing.T) {
	// 1. Setup
	// Ensure config is set for testing
	// We rely on default config or env vars. For test, defaults should work if port is free.
	// But we might want to set a random port?
	// Config loader uses viper singleton. We can try to clean it up?
	// But let's assume default port 8080 is okay or we check if used.
	// Ideally we set a different port.
	// Since we can't easily modify config loaded inside New/Init without hacks,
	// checking if we can mock or just run.

	// Ideally we would set env vars to configure port
	t.Setenv("GROUTER_WEB_PORT", "8888")              // Use a specific test port
	t.Setenv("GROUTER_NATS_ENABLED", "false")         // Disable NATS for pure web test
	t.Setenv("GROUTER_WEB_SECURITY_ENABLED", "false") // Disable Security middleware to avoid SSL redirects

	// Point to the root config file relative to this test file
	os.Args = []string{"test", "--config", "../../../../configs/config.yaml"}

	appInstance := New()
	require.NoError(t, appInstance.Init())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 2. Start App in background
	errChan := make(chan error, 1)
	go func() {
		errChan <- appInstance.Start(ctx)
	}()

	// Wait for server to be ready
	baseURL := "http://localhost:8888"
	require.Eventually(t, func() bool {
		resp, err := http.Get(baseURL + "/health/live")
		if err != nil {
			return false
		}
		defer resp.Body.Close()
		return resp.StatusCode == http.StatusOK
	}, 5*time.Second, 100*time.Millisecond, "Web server failed to start")

	// 3. Test Lifecycle

	// Call /start (Bootstrap service) to register services
	resp, err := http.Get(baseURL + "/start")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()

	// Verify /hello works immediately after start
	// Wait for server to restart and register routes
	require.Eventually(t, func() bool {
		resp, err := http.Get(baseURL + "/hello")
		if err != nil {
			return false
		}
		defer resp.Body.Close()
		return resp.StatusCode == http.StatusOK
	}, 5*time.Second, 100*time.Millisecond, "Web server failed to restart or register /hello")

	// Verify /echo works
	resp, err = http.Get(baseURL + "/echo?msg=integration")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()

	// Call /stop
	resp, err = http.Get(baseURL + "/stop")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()

	// 4. Cleanup
	cancel()
	select {
	case err := <-errChan:
		if err != nil && err != context.Canceled {
			assert.NoError(t, err)
		}
	case <-time.After(2 * time.Second):
		t.Log("App shutdown timed out")
	}
}
