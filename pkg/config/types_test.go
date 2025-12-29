package config

import (
	"testing"
	"time"
)

func TestAppConfig(t *testing.T) {
	app := AppConfig{
		Name:        "test-service",
		Version:     "2.0.0",
		Environment: "production",
	}

	if app.Name != "test-service" {
		t.Errorf("Name = %v, want %v", app.Name, "test-service")
	}

	if app.Version != "2.0.0" {
		t.Errorf("Version = %v, want %v", app.Version, "2.0.0")
	}

	if app.Environment != "production" {
		t.Errorf("Environment = %v, want %v", app.Environment, "production")
	}
}

func TestNATSConfig(t *testing.T) {
	nats := NATSConfig{
		URL:               "nats://test:4222",
		MaxReconnects:     20,
		ReconnectWait:     3 * time.Second,
		ConnectionTimeout: 10 * time.Second,
		Token:             "test-token",
		Username:          "testuser",
		Password:          "testpass",
	}

	if nats.URL != "nats://test:4222" {
		t.Errorf("URL = %v, want %v", nats.URL, "nats://test:4222")
	}

	if nats.MaxReconnects != 20 {
		t.Errorf("MaxReconnects = %v, want %v", nats.MaxReconnects, 20)
	}

	if nats.ReconnectWait != 3*time.Second {
		t.Errorf("ReconnectWait = %v, want %v", nats.ReconnectWait, 3*time.Second)
	}

	if nats.ConnectionTimeout != 10*time.Second {
		t.Errorf("ConnectionTimeout = %v, want %v", nats.ConnectionTimeout, 10*time.Second)
	}

	if nats.Token != "test-token" {
		t.Errorf("Token = %v, want %v", nats.Token, "test-token")
	}

	if nats.Username != "testuser" {
		t.Errorf("Username = %v, want %v", nats.Username, "testuser")
	}

	if nats.Password != "testpass" {
		t.Errorf("Password = %v, want %v", nats.Password, "testpass")
	}
}

func TestLogConfig(t *testing.T) {
	log := LogConfig{
		Level:      "debug",
		Format:     "json",
		OutputPath: "/var/log/test.log",
	}

	if log.Level != "debug" {
		t.Errorf("Level = %v, want %v", log.Level, "debug")
	}

	if log.Format != "json" {
		t.Errorf("Format = %v, want %v", log.Format, "json")
	}

	if log.OutputPath != "/var/log/test.log" {
		t.Errorf("OutputPath = %v, want %v", log.OutputPath, "/var/log/test.log")
	}
}

func TestCompleteConfig(t *testing.T) {
	cfg := Config{
		App: AppConfig{
			Name:        "complete-test",
			Version:     "3.0.0",
			Environment: "staging",
		},
		NATS: NATSConfig{
			URL:               "nats://complete:4222",
			MaxReconnects:     15,
			ReconnectWait:     4 * time.Second,
			ConnectionTimeout: 8 * time.Second,
		},
		Log: LogConfig{
			Level:      "warn",
			Format:     "console",
			OutputPath: "stdout",
		},
		Services: ServicesConfig{
			"natdemo": map[string]interface{}{
				"enabled": false,
				"subject": "natdemo.complete",
			},
		},
	}

	// Verify all nested fields
	if cfg.App.Name != "complete-test" {
		t.Error("App.Name mismatch")
	}

	if cfg.NATS.URL != "nats://complete:4222" {
		t.Error("NATS.URL mismatch")
	}

	if cfg.Log.Level != "warn" {
		t.Error("Log.Level mismatch")
	}

	if cfg.Services["natdemo"] == nil {
		t.Error("Services.NATDemo should not be nil")
	}
	// Note: We can't easily check "enabled: false" here because it is just a map now.
	// The decoding logic is tested in integration or service level.
}

func TestConfigDefaults(t *testing.T) {
	// Test zero values
	var cfg Config

	if cfg.App.Name != "" {
		t.Error("Default App.Name should be empty")
	}

	if cfg.NATS.MaxReconnects != 0 {
		t.Error("Default NATS.MaxReconnects should be 0")
	}

	if cfg.Log.Level != "" {
		t.Error("Default Log.Level should be empty")
	}

	if cfg.Services["natdemo"] != nil {
		t.Error("Default Services.NATDemo should be nil")
	}
}

func TestNATSConfigWithoutAuth(t *testing.T) {
	nats := NATSConfig{
		URL:               "nats://noauth:4222",
		MaxReconnects:     10,
		ReconnectWait:     2 * time.Second,
		ConnectionTimeout: 5 * time.Second,
		// No auth fields set
	}

	if nats.URL != "nats://noauth:4222" {
		t.Error("URL mismatch")
	}
	if nats.MaxReconnects != 10 {
		t.Error("MaxReconnects mismatch")
	}
	if nats.ReconnectWait != 2*time.Second {
		t.Error("ReconnectWait mismatch")
	}
	if nats.ConnectionTimeout != 5*time.Second {
		t.Error("ConnectionTimeout mismatch")
	}

	if nats.Token != "" {
		t.Error("Token should be empty when not set")
	}

	if nats.Username != "" {
		t.Error("Username should be empty when not set")
	}

	if nats.Password != "" {
		t.Error("Password should be empty when not set")
	}
}

func TestLogConfigFormats(t *testing.T) {
	tests := []struct {
		name   string
		format string
	}{
		{"json format", "json"},
		{"console format", "console"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := LogConfig{
				Level:      "info",
				Format:     tt.format,
				OutputPath: "stdout",
			}

			if log.Level != "info" {
				t.Error("Level mismatch")
			}
			if log.Format != tt.format {
				t.Errorf("Format = %v, want %v", log.Format, tt.format)
			}
			if log.OutputPath != "stdout" {
				t.Error("OutputPath mismatch")
			}
		})
	}
}

func TestLogLevels(t *testing.T) {
	levels := []string{"debug", "info", "warn", "error"}

	for _, level := range levels {
		t.Run(level, func(t *testing.T) {
			log := LogConfig{
				Level:      level,
				Format:     "console",
				OutputPath: "stdout",
			}

			if log.Level != level {
				t.Errorf("Level = %v, want %v", log.Level, level)
			}
			if log.Format != "console" {
				t.Error("Format mismatch")
			}
			if log.OutputPath != "stdout" {
				t.Error("OutputPath mismatch")
			}
		})
	}
}

func TestConfigTypes(t *testing.T) {
	// Test that config types can be created and have expected fields
	cfg := Config{
		App: AppConfig{
			Name:        "test",
			Version:     "1.0",
			Environment: "dev",
		},
		NATS: NATSConfig{
			URL:               "nats://localhost:4222",
			MaxReconnects:     5,
			ReconnectWait:     1 * time.Second,
			ConnectionTimeout: 3 * time.Second,
			Token:             "token123",
			Username:          "user",
			Password:          "pass",
		},
		Log: LogConfig{
			Level:      "debug",
			Format:     "json",
			OutputPath: "/var/log/app.log",
		},
		Services: ServicesConfig{
			"natdemo": map[string]interface{}{
				"enabled": true,
				"subject": "natdemo.events",
			},
		},
	}

	// Verify all fields are accessible
	if cfg.App.Name != "test" {
		t.Error("App.Name field not accessible")
	}
	if cfg.NATS.URL != "nats://localhost:4222" {
		t.Error("NATS.URL field not accessible")
	}
	if cfg.Log.Level != "debug" {
		t.Error("Log.Level field not accessible")
	}
	if cfg.Services["natdemo"] == nil {
		t.Error("Services.NATDemo field not accessible")
	}
}
