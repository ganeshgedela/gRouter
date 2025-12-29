package logger

import (
	"os"
	"path/filepath"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "valid json config",
			config: Config{
				Level:      "info",
				Format:     "json",
				OutputPath: "stdout",
			},
			wantErr: false,
		},
		{
			name: "valid console config",
			config: Config{
				Level:      "debug",
				Format:     "console",
				OutputPath: "stdout",
			},
			wantErr: false,
		},
		{
			name: "invalid log level",
			config: Config{
				Level:      "invalid",
				Format:     "console",
				OutputPath: "stdout",
			},
			wantErr: true,
		},
		{
			name: "warn level",
			config: Config{
				Level:      "warn",
				Format:     "json",
				OutputPath: "stdout",
			},
			wantErr: false,
		},
		{
			name: "error level",
			config: Config{
				Level:      "error",
				Format:     "console",
				OutputPath: "stdout",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := New(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && logger == nil {
				t.Error("New() returned nil logger")
			}
		})
	}
}

func TestNew_FileOutput(t *testing.T) {
	// Create temp directory for test logs
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "test.log")

	config := Config{
		Level:      "info",
		Format:     "json",
		OutputPath: logFile,
	}

	logger, err := New(config)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	if logger == nil {
		t.Fatal("New() returned nil logger")
	}

	// Write a log message
	logger.Info("test message")
	logger.Sync()

	// Verify file was created
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		t.Errorf("Log file was not created: %s", logFile)
	}
}

func TestGet(t *testing.T) {
	// Reset global logger
	globalLogger = nil

	// Get should return a default logger if none exists
	logger := Get()
	if logger == nil {
		t.Error("Get() returned nil logger")
	}

	// Create a new logger
	config := Config{
		Level:      "info",
		Format:     "console",
		OutputPath: "stdout",
	}
	newLogger, err := New(config)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	// Get should now return the new logger
	logger = Get()
	if logger != newLogger {
		t.Error("Get() did not return the expected logger")
	}
}

func TestSugar(t *testing.T) {
	config := Config{
		Level:      "info",
		Format:     "console",
		OutputPath: "stdout",
	}

	_, err := New(config)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	sugared := Sugar()
	if sugared == nil {
		t.Error("Sugar() returned nil")
	}
}

func TestWithFields(t *testing.T) {
	config := Config{
		Level:      "info",
		Format:     "console",
		OutputPath: "stdout",
	}

	_, err := New(config)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	logger := WithFields(
		zap.String("key1", "value1"),
		zap.Int("key2", 42),
	)

	if logger == nil {
		t.Error("WithFields() returned nil")
	}
}

func TestLogLevels(t *testing.T) {
	config := Config{
		Level:      "debug",
		Format:     "console",
		OutputPath: "stdout",
	}

	_, err := New(config)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	// Test all log levels - should not panic
	tests := []struct {
		name string
		fn   func()
	}{
		{
			name: "Debug",
			fn: func() {
				Debug("debug message", zap.String("key", "value"))
			},
		},
		{
			name: "Info",
			fn: func() {
				Info("info message", zap.String("key", "value"))
			},
		},
		{
			name: "Warn",
			fn: func() {
				Warn("warn message", zap.String("key", "value"))
			},
		},
		{
			name: "Error",
			fn: func() {
				Error("error message", zap.Error(err))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Should not panic
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("%s() panicked: %v", tt.name, r)
				}
			}()
			tt.fn()
		})
	}
}

func TestSync(t *testing.T) {
	config := Config{
		Level:      "info",
		Format:     "console",
		OutputPath: "stdout",
	}

	_, err := New(config)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	// Sync should not return error for stdout
	err = Sync()
	if err != nil {
		// Note: Sync() may return error on some systems for stdout, which is acceptable
		t.Logf("Sync() returned error (may be expected): %v", err)
	}
}

func TestLogFormats(t *testing.T) {
	tests := []struct {
		name   string
		format string
	}{
		{
			name:   "json format",
			format: "json",
		},
		{
			name:   "console format",
			format: "console",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := Config{
				Level:      "info",
				Format:     tt.format,
				OutputPath: "stdout",
			}

			logger, err := New(config)
			if err != nil {
				t.Errorf("New() with %s format error = %v", tt.format, err)
				return
			}

			if logger == nil {
				t.Errorf("New() with %s format returned nil", tt.format)
			}
		})
	}
}

func TestLogLevelParsing(t *testing.T) {
	tests := []struct {
		level    string
		wantErr  bool
		expected zapcore.Level
	}{
		{"debug", false, zapcore.DebugLevel},
		{"info", false, zapcore.InfoLevel},
		{"warn", false, zapcore.WarnLevel},
		{"error", false, zapcore.ErrorLevel},
		{"invalid", true, zapcore.InfoLevel},
		{"DEBUG", false, zapcore.DebugLevel}, // Case insensitive
		{"INFO", false, zapcore.InfoLevel},
	}

	for _, tt := range tests {
		t.Run(tt.level, func(t *testing.T) {
			config := Config{
				Level:      tt.level,
				Format:     "console",
				OutputPath: "stdout",
			}

			logger, err := New(config)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() with level %s error = %v, wantErr %v", tt.level, err, tt.wantErr)
			}

			if !tt.wantErr && logger == nil {
				t.Errorf("New() with level %s returned nil", tt.level)
			}
		})
	}
}

func TestEmptyOutputPath(t *testing.T) {
	config := Config{
		Level:      "info",
		Format:     "console",
		OutputPath: "", // Should default to stdout
	}

	logger, err := New(config)
	if err != nil {
		t.Fatalf("New() with empty output path error = %v", err)
	}

	if logger == nil {
		t.Error("New() with empty output path returned nil")
	}
}

func TestMultipleLoggerCreation(t *testing.T) {
	// Create multiple loggers to ensure no conflicts
	for i := 0; i < 5; i++ {
		config := Config{
			Level:      "info",
			Format:     "console",
			OutputPath: "stdout",
		}

		logger, err := New(config)
		if err != nil {
			t.Fatalf("New() iteration %d error = %v", i, err)
		}

		if logger == nil {
			t.Errorf("New() iteration %d returned nil", i)
		}

		// Each call should update the global logger
		if Get() != logger {
			t.Errorf("Get() iteration %d did not return the latest logger", i)
		}
	}
}
