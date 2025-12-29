package database

import (
	"context"
	"testing"

	"grouter/pkg/config"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type TestModel struct {
	ID   uint `gorm:"primarykey"`
	Name string
}

func TestNewDatabase_SQLite(t *testing.T) {
	// Setup
	logger := zap.NewNop()
	cfg := config.DatabaseConfig{
		Driver:   "sqlite",
		DBName:   ":memory:",
		LogLevel: "info",
	}

	// Test New
	db, err := New(cfg, logger)
	assert.NoError(t, err)
	assert.NotNil(t, db)

	// Test HealthCheck
	err = db.HealthCheck(context.Background())
	assert.NoError(t, err)

	// Test GORM functionality
	err = db.AutoMigrate(&TestModel{})
	assert.NoError(t, err)

	// Insert
	item := TestModel{Name: "test"}
	result := db.Create(&item)
	assert.NoError(t, result.Error)
	assert.Equal(t, uint(1), item.ID)

	// Query
	var readItem TestModel
	result = db.First(&readItem, 1)
	assert.NoError(t, result.Error)
	assert.Equal(t, "test", readItem.Name)
}
