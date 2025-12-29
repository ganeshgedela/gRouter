package app

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestApp_Init(t *testing.T) {
	// Setup temporary config
	pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.yaml")

	configContent := `
app:
  name: "nats-demo-test"
  version: "1.0.0"
  environment: "test"

nats:
  enabled: false
  url: "nats://localhost:4222"

log:
  level: "info"
  format: "console"
  output_path: "stdout"
`
	err := os.WriteFile(configFile, []byte(configContent), 0644)
	assert.NoError(t, err)

	// Reset viper
	viper.Reset()
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	// Point to our temp config
	os.Args = []string{"test_app", "--config", configFile}

	app := New()
	err = app.Init()
	assert.NoError(t, err)

	assert.Equal(t, "nats-demo-test", app.GetAppName())
	assert.NotEmpty(t, app.AppId)
	assert.Contains(t, app.AppId, "nats-demo-test-")
}

func TestApp_New(t *testing.T) {
	app := New()
	assert.NotNil(t, app)
	assert.NotNil(t, app.manager)
	assert.NotNil(t, app.startChan)
	assert.NotNil(t, app.stopChan)
}
