package products

import (
	"testing"
	"time"

	productv1 "grouter/templates/db-service/api/v1"

	"github.com/stretchr/testify/assert"
)

func TestProductToProto(t *testing.T) {
	now := time.Now()
	product := &Product{
		ID:          1,
		Name:        "Test Product",
		Description: "Description",
		Price:       10.50,
		SKU:         "TEST-SKU",
		Stock:       100,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	proto := product.ToProto()

	assert.Equal(t, uint64(1), proto.Id)
	assert.Equal(t, "Test Product", proto.Name)
	assert.Equal(t, "Description", proto.Description)
	assert.Equal(t, 10.50, proto.Price)
	assert.Equal(t, "TEST-SKU", proto.Sku)
	assert.Equal(t, int32(100), proto.Stock)
	assert.NotNil(t, proto.CreatedAt)
	assert.NotNil(t, proto.UpdatedAt)
}

func TestFromCreateRequest(t *testing.T) {
	req := &productv1.CreateProductRequest{
		Name:        "New Product",
		Description: "New Desc",
		Price:       20.00,
		Sku:         "NEW-SKU",
		Stock:       50,
	}

	product := FromCreateRequest(req)

	assert.Equal(t, "New Product", product.Name)
	assert.Equal(t, "New Desc", product.Description)
	assert.Equal(t, 20.00, product.Price)
	assert.Equal(t, "NEW-SKU", product.SKU)
	assert.Equal(t, 50, product.Stock)
}

func TestFromUpdateRequest(t *testing.T) {
	product := &Product{
		Name:  "Old Name",
		Price: 10.00,
		Stock: 50,
	}

	// Update only name and price
	req := &productv1.UpdateProductRequest{
		Name:  "New Name",
		Price: 15.00,
		Stock: -1,
	}

	product.FromUpdateRequest(req)

	assert.Equal(t, "New Name", product.Name)
	assert.Equal(t, 15.00, product.Price)
	assert.Equal(t, 50, product.Stock) // Should remain unchanged
}
