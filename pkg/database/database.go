package database

import (
	"context"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/plugin/opentelemetry/tracing"

	"grouter/pkg/config"

	"go.uber.org/zap"
)

// Database wraps the GORM DB connection
type Database struct {
	*gorm.DB
}

// New creates a new database connection based on configuration
func New(cfg config.DatabaseConfig, logger *zap.Logger) (*Database, error) {
	var dialect gorm.Dialector

	switch cfg.Driver {
	case "postgres":
		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)
		dialect = postgres.Open(dsn)
	case "sqlite", "sqlite3":
		dialect = sqlite.Open(cfg.DBName)
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Driver)
	}

	// Configure GORM Logger
	gormLog := NewGormLogger(logger, cfg.LogLevel)

	gormConfig := &gorm.Config{
		Logger: gormLog,
	}

	db, err := gorm.Open(dialect, gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Integrate OpenTelemetry Tracing
	if err := db.Use(tracing.NewPlugin()); err != nil {
		logger.Warn("failed to register opentelemetry plugin for gorm", zap.Error(err))
	}

	// Configure Connection Pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB from gorm: %w", err)
	}

	if cfg.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	}
	if cfg.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	}
	if cfg.ConnMaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	}

	return &Database{DB: db}, nil
}

// WithTransaction executes a function within a database transaction
func (d *Database) WithTransaction(ctx context.Context, fn func(txDB *Database) error) error {
	return d.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(&Database{DB: tx})
	})
}

// HealthCheck executes a simple query to verify database connectivity
func (d *Database) HealthCheck(ctx context.Context) error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}
	return sqlDB.PingContext(ctx)
}

// NewGormLogger creates a GORM logger that outputs to a Zap logger
func NewGormLogger(zapLogger *zap.Logger, logLevel string) gormlogger.Interface {
	var lvl gormlogger.LogLevel
	switch logLevel {
	case "silent":
		lvl = gormlogger.Silent
	case "error":
		lvl = gormlogger.Error
	case "warn":
		lvl = gormlogger.Warn
	case "info":
		lvl = gormlogger.Info
	default:
		lvl = gormlogger.Warn // Default to Warn
	}

	return &zapGormLogger{
		ZapLogger:     zapLogger,
		LogLevel:      lvl,
		SlowThreshold: 200 * time.Millisecond,
	}
}

type zapGormLogger struct {
	ZapLogger     *zap.Logger
	LogLevel      gormlogger.LogLevel
	SlowThreshold time.Duration
}

func (l *zapGormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

func (l *zapGormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Info {
		l.ZapLogger.Sugar().Infof(msg, data...)
	}
}

func (l *zapGormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Warn {
		l.ZapLogger.Sugar().Warnf(msg, data...)
	}
}

func (l *zapGormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Error {
		l.ZapLogger.Sugar().Errorf(msg, data...)
	}
}

func (l *zapGormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= gormlogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	fields := []zap.Field{
		zap.String("sql", sql),
		zap.Int64("rows", rows),
		zap.Duration("elapsed", elapsed),
	}

	if err != nil && l.LogLevel >= gormlogger.Error {
		l.ZapLogger.Error("trace", append(fields, zap.Error(err))...)
		return
	}

	if l.SlowThreshold != 0 && elapsed > l.SlowThreshold && l.LogLevel >= gormlogger.Warn {
		l.ZapLogger.Warn("slow sql", fields...)
		return
	}

	if l.LogLevel >= gormlogger.Info {
		l.ZapLogger.Info("trace", fields...)
	}
}
