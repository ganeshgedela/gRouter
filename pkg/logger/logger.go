package logger

import (
	"fmt"
	"os"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	globalLogger *zap.Logger
	sugar        *zap.SugaredLogger
	mu           sync.RWMutex
)

// Config holds logger configuration
type Config struct {
	Level      string
	Format     string // json or console
	OutputPath string

	// Log rotation settings (only applies when OutputPath is a file)
	MaxSize    int  // Maximum size in megabytes before rotation (default: 100)
	MaxBackups int  // Maximum number of old log files to retain (default: 3)
	MaxAge     int  // Maximum number of days to retain old log files (default: 28)
	Compress   bool // Whether to compress rotated files (default: true)

	// Sampling configuration for high-volume environments
	EnableSampling     bool // Enable log sampling to prevent flooding
	SamplingInitial    int  // Initial number of messages logged per second (default: 100)
	SamplingThereafter int  // Messages logged after initial threshold (default: 100)
}

// New creates a new logger instance with production-ready features
func New(cfg Config) (*zap.Logger, error) {
	// Set defaults
	if cfg.MaxSize == 0 {
		cfg.MaxSize = 100
	}
	if cfg.MaxBackups == 0 {
		cfg.MaxBackups = 3
	}
	if cfg.MaxAge == 0 {
		cfg.MaxAge = 28
	}
	if cfg.SamplingInitial == 0 {
		cfg.SamplingInitial = 100
	}
	if cfg.SamplingThereafter == 0 {
		cfg.SamplingThereafter = 100
	}

	// Parse log level
	level, err := zapcore.ParseLevel(cfg.Level)
	if err != nil {
		return nil, fmt.Errorf("invalid log level: %w", err)
	}

	// Configure encoder
	var encoderConfig zapcore.EncoderConfig
	if cfg.Format == "json" {
		encoderConfig = zap.NewProductionEncoderConfig()
	} else {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.StacktraceKey = "" // Disable stacktrace

	// Create encoder
	var encoder zapcore.Encoder
	if cfg.Format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// Configure output with proper file handle management
	var writer zapcore.WriteSyncer
	if cfg.OutputPath == "" || cfg.OutputPath == "stdout" {
		// Use locked stdout to prevent concurrent write issues
		writer = zapcore.Lock(os.Stdout)
	} else {
		// Use lumberjack for automatic log rotation
		writer = zapcore.AddSync(&lumberjack.Logger{
			Filename:   cfg.OutputPath,
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
			Compress:   cfg.Compress,
		})
	}

	// Create core
	core := zapcore.NewCore(encoder, writer, level)

	// Apply sampling if enabled
	if cfg.EnableSampling {
		core = zapcore.NewSamplerWithOptions(
			core,
			time.Second,
			cfg.SamplingInitial,
			cfg.SamplingThereafter,
		)
	}

	// Create logger with caller info
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(0))

	// Thread-safe update of global logger
	mu.Lock()
	globalLogger = logger
	sugar = logger.Sugar()
	mu.Unlock()

	return logger, nil
}

// Get returns the global logger (thread-safe)
func Get() *zap.Logger {
	mu.RLock()
	logger := globalLogger
	mu.RUnlock()

	if logger == nil {
		// Create a default logger if none exists
		mu.Lock()
		if globalLogger == nil {
			logger, _ := zap.NewProduction()
			globalLogger = logger
			sugar = logger.Sugar()
		}
		mu.Unlock()
		return globalLogger
	}
	return logger
}

// Sugar returns the global sugared logger (thread-safe)
func Sugar() *zap.SugaredLogger {
	mu.RLock()
	s := sugar
	mu.RUnlock()

	if s == nil {
		mu.Lock()
		if sugar == nil {
			sugar = Get().Sugar()
		}
		s = sugar
		mu.Unlock()
	}
	return s
}

// WithFields creates a logger with additional fields
func WithFields(fields ...zap.Field) *zap.Logger {
	return Get().With(fields...)
}

// Debug logs a debug message
func Debug(msg string, fields ...zap.Field) {
	Get().Debug(msg, fields...)
}

// Info logs an info message
func Info(msg string, fields ...zap.Field) {
	Get().Info(msg, fields...)
}

// Warn logs a warning message
func Warn(msg string, fields ...zap.Field) {
	Get().Warn(msg, fields...)
}

// Error logs an error message
func Error(msg string, fields ...zap.Field) {
	Get().Error(msg, fields...)
}

// Fatal logs a fatal message and exits
func Fatal(msg string, fields ...zap.Field) {
	Get().Fatal(msg, fields...)
}

// Sync flushes any buffered log entries
func Sync() error {
	if globalLogger != nil {
		return globalLogger.Sync()
	}
	return nil
}
