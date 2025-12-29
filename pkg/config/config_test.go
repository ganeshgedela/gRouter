package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// resetConfig resets global state between tests
func resetConfig() {
	globalConfig = nil
	viper.Reset()
	pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)
}

func TestLoad(t *testing.T) {
	resetConfig()

	// Create a temporary config file
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.yaml")

	configContent := `
app:
  name: "test-app"
  version: "1.0.0"
  environment: "test"

nats:
  url: "nats://localhost:4222"
  max_reconnects: 10
  reconnect_wait: 2s
  connection_timeout: 5s

log:
  level: "info"
  format: "console"
  output_path: "stdout"

services:
  natdemo:
    enabled: true
    subject: "natdemo"
`

	err := os.WriteFile(configFile, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	// Set command line args
	os.Args = []string{"test", "--config", configFile}

	// Load config
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	// Verify config values
	if cfg.App.Name != "test-app" {
		t.Errorf("App.Name = %v, want %v", cfg.App.Name, "test-app")
	}

	if cfg.App.Version != "1.0.0" {
		t.Errorf("App.Version = %v, want %v", cfg.App.Version, "1.0.0")
	}

	if cfg.NATS.URL != "nats://localhost:4222" {
		t.Errorf("NATS.URL = %v, want %v", cfg.NATS.URL, "nats://localhost:4222")
	}

	if cfg.Log.Level != "info" {
		t.Errorf("Log.Level = %v, want %v", cfg.Log.Level, "info")
	}

	// Check that natdemo config is present in map
	if cfg.Services["natdemo"] == nil {
		t.Error("Services['natdemo'] should be present")
	}
}

func TestLoad_InvalidFile(t *testing.T) {
	resetConfig()

	// Set non-existent config file
	os.Args = []string{"test", "--config", "/nonexistent/config.yaml"}

	_, err := Load()
	if err == nil {
		t.Error("Load() should return error for non-existent file")
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	resetConfig()

	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "invalid.yaml")

	// Write invalid YAML
	invalidYAML := `
app:
  name: "test
  invalid yaml here
`
	err := os.WriteFile(configFile, []byte(invalidYAML), 0644)
	if err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	os.Args = []string{"test", "--config", configFile}

	_, err = Load()
	if err == nil {
		t.Error("Load() should return error for invalid YAML")
	}
}

func TestLoad_WithFlags(t *testing.T) {
	resetConfig()

	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.yaml")

	configContent := `
app:
  name: "test-app"
  version: "1.0.0"
  environment: "test"

nats:
  url: "nats://localhost:4222"
  max_reconnects: 10
  reconnect_wait: 2s
  connection_timeout: 5s

log:
  level: "info"
  format: "console"
  output_path: "stdout"

services:
  natdemo:
    enabled: true
    subject: "natdemo"
`

	err := os.WriteFile(configFile, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	// Set command line args with overrides
	os.Args = []string{
		"test",
		"--config", configFile,
		"--log-level", "debug",
		"--nats-url", "nats://custom:4222",
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	// Verify flag overrides
	if cfg.Log.Level != "debug" {
		t.Errorf("Log.Level = %v, want %v (should be overridden by flag)", cfg.Log.Level, "debug")
	}

	if cfg.NATS.URL != "nats://custom:4222" {
		t.Errorf("NATS.URL = %v, want %v (should be overridden by flag)", cfg.NATS.URL, "nats://custom:4222")
	}
}

func TestLoad_WithEnvVars(t *testing.T) {
	resetConfig()

	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.yaml")

	configContent := `
app:
  name: "test-app"
  version: "1.0.0"
  environment: "test"

nats:
  url: "nats://localhost:4222"
  max_reconnects: 10
  reconnect_wait: 2s
  connection_timeout: 5s

log:
  level: "info"
  format: "console"
  output_path: "stdout"

services:
  natdemo:
    enabled: true
    subject: "natdemo"
`

	err := os.WriteFile(configFile, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	// Set environment variables
	os.Setenv("GROUTER_LOG_LEVEL", "warn")
	os.Setenv("GROUTER_NATS_URL", "nats://env:4222")
	defer func() {
		os.Unsetenv("GROUTER_LOG_LEVEL")
		os.Unsetenv("GROUTER_NATS_URL")
	}()

	os.Args = []string{"test", "--config", configFile}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	// Verify env var overrides
	if cfg.Log.Level != "warn" {
		t.Errorf("Log.Level = %v, want %v (should be overridden by env var)", cfg.Log.Level, "warn")
	}

	if cfg.NATS.URL != "nats://env:4222" {
		t.Errorf("NATS.URL = %v, want %v (should be overridden by env var)", cfg.NATS.URL, "nats://env:4222")
	}
}

func TestGet(t *testing.T) {
	resetConfig()

	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.yaml")

	configContent := `
app:
  name: "test-app"
  version: "1.0.0"
  environment: "test"

nats:
  url: "nats://localhost:4222"
  max_reconnects: 10
  reconnect_wait: 2s
  connection_timeout: 5s

log:
  level: "info"
  format: "console"
  output_path: "stdout"

services:
  natdemo:
    enabled: true
    subject: "natdemo"
`

	err := os.WriteFile(configFile, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	os.Args = []string{"test", "--config", configFile}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	// Get should return the same config
	retrievedCfg := Get()
	if retrievedCfg != cfg {
		t.Error("Get() did not return the loaded config")
	}

	// Verify config values
	if retrievedCfg.App.Name != "test-app" {
		t.Errorf("Retrieved config App.Name = %v, want %v", retrievedCfg.App.Name, "test-app")
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: Config{
				App: AppConfig{
					Name:        "test-app",
					Version:     "1.0.0",
					Environment: "test",
				},
				NATS: NATSConfig{
					URL:               "nats://localhost:4222",
					MaxReconnects:     10,
					ReconnectWait:     2 * time.Second,
					ConnectionTimeout: 5 * time.Second,
				},
				Log: LogConfig{
					Level:      "info",
					Format:     "console",
					OutputPath: "stdout",
				},
			},
			wantErr: false,
		},
		{
			name: "missing app name",
			config: Config{
				App: AppConfig{
					Name: "",
				},
				NATS: NATSConfig{
					URL: "nats://localhost:4222",
				},
				Log: LogConfig{
					Level: "info",
				},
			},
			wantErr: true,
		},
		{
			name: "missing nats url",
			config: Config{
				App: AppConfig{
					Name: "test-app",
				},
				NATS: NATSConfig{
					Enabled: true,
					URL:     "",
				},
				Log: LogConfig{
					Level: "info",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid log level",
			config: Config{
				App: AppConfig{
					Name: "test-app",
				},
				NATS: NATSConfig{
					URL: "nats://localhost:4222",
				},
				Log: LogConfig{
					Level: "invalid",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate(&tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
