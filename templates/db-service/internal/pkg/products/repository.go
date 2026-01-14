package products

import (
	"context"
	"fmt"

	"grouter/pkg/database"

	"gorm.io/gorm"
)

// Repository extends the generic database repository with product-specific methods
type Repository struct {
	database.Repository[Product]
	db *gorm.DB
}

// NewRepository creates a new product repository
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		Repository: database.NewRepository[Product](db),
		db:         db,
	}
}

// FindBySKU finds a product by its SKU
func (r *Repository) FindBySKU(ctx context.Context, sku string) (*Product, error) {
	var product Product
	if err := r.db.WithContext(ctx).Where("sku = ?", sku).First(&product).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

// UpdateStock atomically updates the stock of a product
func (r *Repository) UpdateStock(ctx context.Context, id uint, delta int) (*Product, error) {
	var product Product
	
	// Use transaction for atomic update
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Lock the row for update
		if err := tx.Clauses(gorm.Locking{Strength: "UPDATE"}).First(&product, id).Error; err != nil {
			return err
		}
		
		newStock := product.Stock + delta
		if newStock < 0 {
			return fmt.Errorf("insufficient stock: current=%d, delta=%d", product.Stock, delta)
		}
		
		product.Stock = newStock
		return tx.Save(&product).Error
	})
	
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// BulkCreate creates multiple products in a single transaction
func (r *Repository) BulkCreate(ctx context.Context, products []*Product) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.CreateInBatches(products, 100).Error
	})
}

// SearchByName searches products by name using LIKE query
func (r *Repository) SearchByName(ctx context.Context, query string, pagination database.Pagination) ([]Product, int64, error) {
	var products []Product
	var total int64
	
	db := r.db.WithContext(ctx).Model(&Product{})
	
	// Add search filter
	if query != "" {
		likeQuery := "%" + query + "%"
		db = db.Where("name LIKE ? OR description LIKE ?", likeQuery, likeQuery)
	}
	
	// Apply pagination filters
	if len(pagination.Filters) > 0 {
		db = db.Where(pagination.Filters)
	}
	
	// Count total
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// Apply sorting and pagination
	if pagination.Sort != "" {
		db = db.Order(pagination.Sort)
	}
	
	err := db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit()).Find(&products).Error
	if err != nil {
		return nil, 0, err
	}
	
	return products, total, nil
}
