package business

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestItemService(t *testing.T) {
	logger := zap.NewNop()
	service := NewItemService(logger)
	ctx := context.Background()

	t.Run("CreateItem", func(t *testing.T) {
		item, err := service.CreateItem(ctx, "Test Item", "Test Description", 10.0, 5)
		assert.NoError(t, err)
		assert.NotNil(t, item)
		assert.Equal(t, "Test Item", item.Name)
		assert.Equal(t, 10.0, item.Price)
		assert.Equal(t, int32(5), item.Quantity)
		assert.NotEmpty(t, item.ID)
	})

	t.Run("GetItem", func(t *testing.T) {
		// First create an item
		created, _ := service.CreateItem(ctx, "Get Item", "Desc", 20.0, 1)

		// Then get it
		item, err := service.GetItem(ctx, created.ID)
		assert.NoError(t, err)
		assert.Equal(t, created.ID, item.ID)
		assert.Equal(t, "Get Item", item.Name)
	})

	t.Run("GetItemNotFound", func(t *testing.T) {
		_, err := service.GetItem(ctx, "non-existent")
		assert.Error(t, err)
	})

	t.Run("ListItems", func(t *testing.T) {
		// We already created 2 items above
		items, err := service.ListItems(ctx)
		assert.NoError(t, err)
		assert.Len(t, items, 2)
	})

	t.Run("DeleteItem", func(t *testing.T) {
		created, _ := service.CreateItem(ctx, "Delete Item", "Desc", 30.0, 1)

		err := service.DeleteItem(ctx, created.ID)
		assert.NoError(t, err)

		// Verify it's gone
		_, err = service.GetItem(ctx, created.ID)
		assert.Error(t, err)
	})
}
