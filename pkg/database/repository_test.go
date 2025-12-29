package database

import (
	"context"
	"testing"

	"grouter/pkg/config"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type User struct {
	ID    uint `gorm:"primarykey"`
	Name  string
	Email string
}

func TestGORMRepository(t *testing.T) {
	// Setup
	logger := zap.NewNop()
	cfg := config.DatabaseConfig{
		Driver:   "sqlite",
		DBName:   ":memory:",
		LogLevel: "silent",
	}

	db, err := New(cfg, logger)
	assert.NoError(t, err)
	assert.NotNil(t, db)

	err = db.AutoMigrate(&User{})
	assert.NoError(t, err)

	// Init Repository
	repo := NewRepository[User](db.DB)

	ctx := context.Background()

	// 1. Test Create
	user := &User{Name: "Alice", Email: "alice@example.com"}
	err = repo.Create(ctx, user)
	assert.NoError(t, err)
	assert.NotZero(t, user.ID)

	// 2. Test FindByID
	foundUser, err := repo.FindByID(ctx, user.ID)
	assert.NoError(t, err)
	assert.Equal(t, user.Name, foundUser.Name)

	// 3. Test Update
	foundUser.Name = "Alice Updated"
	err = repo.Update(ctx, foundUser)
	assert.NoError(t, err)

	updatedUser, err := repo.FindByID(ctx, user.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Alice Updated", updatedUser.Name)

	// 4. Test List
	repo.Create(ctx, &User{Name: "Bob", Email: "bob@example.com"})
	repo.Create(ctx, &User{Name: "Charlie", Email: "charlie@example.com"})

	// Test Pagination and Filtering
	p := Pagination{Page: 1, PageSize: 2, Sort: "name asc"}
	users, total, err := repo.List(ctx, p)
	assert.NoError(t, err)
	assert.Equal(t, int64(3), total) // Total count ignoring limit
	assert.Len(t, users, 2)
	assert.Equal(t, "Alice Updated", users[0].Name) // Alphabetical: Alice, Bob

	// Test Filtering
	pFilter := Pagination{
		Filters: map[string]interface{}{"name": "Charlie"},
	}
	usersF, totalF, err := repo.List(ctx, pFilter)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), totalF)
	assert.Len(t, usersF, 1)
	assert.Equal(t, "Charlie", usersF[0].Name)

	// Test Offset
	p2 := Pagination{Page: 2, PageSize: 2, Sort: "name asc"}
	users2, _, err := repo.List(ctx, p2)
	assert.NoError(t, err)
	assert.Len(t, users2, 1)
	assert.Equal(t, "Charlie", users2[0].Name)

	// 5. Test Delete
	err = repo.Delete(ctx, user.ID)
	assert.NoError(t, err)

	_, err = repo.FindByID(ctx, user.ID)
	assert.Error(t, err) // Should not find
}
