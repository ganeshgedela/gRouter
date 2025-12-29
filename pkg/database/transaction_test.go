package database

import (
	"context"
	"errors"
	"testing"

	"grouter/pkg/config"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestTransactions(t *testing.T) {
	// Setup
	logger := zap.NewNop()
	cfg := config.DatabaseConfig{
		Driver:   "sqlite",
		DBName:   ":memory:",
		LogLevel: "silent",
	}

	db, err := New(cfg, logger)
	assert.NoError(t, err)
	db.AutoMigrate(&User{})

	ctx := context.Background()

	t.Run("Commit", func(t *testing.T) {
		err := db.WithTransaction(ctx, func(txDB *Database) error {
			repo := NewRepository[User](txDB.DB)
			return repo.Create(ctx, &User{Name: "TxUser1"})
		})
		assert.NoError(t, err)

		// Verify persisted
		repo := NewRepository[User](db.DB)
		p := Pagination{Filters: map[string]interface{}{"name": "TxUser1"}}
		users, _, _ := repo.List(ctx, p)
		assert.Len(t, users, 1)
	})

	t.Run("Rollback", func(t *testing.T) {
		err := db.WithTransaction(ctx, func(txDB *Database) error {
			repo := NewRepository[User](txDB.DB)
			repo.Create(ctx, &User{Name: "TxUser2"})
			return errors.New("force rollback")
		})
		assert.Error(t, err)

		// Verify NOT persisted
		repo := NewRepository[User](db.DB)
		p := Pagination{Filters: map[string]interface{}{"name": "TxUser2"}}
		users, _, _ := repo.List(ctx, p)
		assert.Len(t, users, 0)
	})
}
