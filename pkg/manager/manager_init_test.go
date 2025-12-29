package manager

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func resetFlags() {
	pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)
}

func TestServiceManager_Init(t *testing.T) {
	resetFlags()
	// Setup temporary config
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.yaml")

	configContent := `
app:
  name: "test-grouter"
  version: "1.0.0"
  environment: "test"

nats:
  enabled: false
  url: "nats://localhost:4222"

web:
  enabled: false
  port: 8080

log:
  level: "error"
  format: "console"
  output_path: "stdout"

tracing:
  enabled: false
`
	err := os.WriteFile(configFile, []byte(configContent), 0644)
	assert.NoError(t, err)

	// Reset viper
	viper.Reset()
	// Init usually expects certain flags or env vars or default locations.
	// Since Init calls initConfig which calls config.Load, and config.Load uses pflag/viper
	// modifying global flags in parallel tests is risky.
	// However, we can try to set the config file path via viper directly if config.Load respects it?
	// config.Load calls viper.SetConfigFile(configFile) where configFile comes from pflag.

	// We can set arguments to point to our config file
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"test_binary", "--config", configFile}

	mgr := NewServiceManager()
	err = mgr.Init()
	assert.NoError(t, err)
	assert.NotNil(t, mgr.cfg)
	assert.Equal(t, "test-grouter", mgr.cfg.App.Name)
	assert.NotNil(t, mgr.log)

	// Verify NATS is disabled
	assert.Nil(t, mgr.messenger)

	// Verify Web is disabled
	assert.Nil(t, mgr.webServer)
}

func TestServiceManager_Init_WithNATS(t *testing.T) {
	resetFlags()
	// Setup temporary config
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config_nats.yaml")

	configContent := `
app:
  name: "test-grouter-nats"
  version: "1.0.0"
  environment: "test"

nats:
  enabled: true
  url: "nats://localhost:4222"
  max_reconnects: 1
  reconnect_wait: 100ms
  connection_timeout: 100ms

web:
  enabled: false

log:
  level: "error"
  format: "console"
  output_path: "stdout"
`
	err := os.WriteFile(configFile, []byte(configContent), 0644)
	assert.NoError(t, err)

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"test_binary", "--config", configFile}

	// We need to reset viper because config.Load uses it globally
	viper.Reset()

	mgr := NewServiceManager()
	err = mgr.Init()
	assert.NoError(t, err)

	err = mgr.InitNATS()
	// It might find NATS enabled and try to connect.
	// If NATS is not running, it might fail or just log error depending on Messenger.Init implementation.
	// Messenger.Init does Connect() which returns error if connection fails.
	// So we expect error if no NATS available.

	if err != nil {
		t.Logf("Init failed with NATS enabled (expected if no server): %v", err)
		assert.Contains(t, err.Error(), "failed to initialize messenger")
	} else {
		assert.NotNil(t, mgr.messenger)
	}
}
