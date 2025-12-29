package database

import (
	"context"

	"gorm.io/gorm"
)

// Repository defines the standard CRUD interface for any entity T
type Repository[T any] interface {
	Create(ctx context.Context, entity *T) error
	FindByID(ctx context.Context, id interface{}) (*T, error)
	List(ctx context.Context, pagination Pagination) ([]T, int64, error)
	Update(ctx context.Context, entity *T) error
	Delete(ctx context.Context, id interface{}) error
}

// GORMRepository implements Repository[T] using GORM
type GORMRepository[T any] struct {
	db *gorm.DB
}

// NewRepository creates a new generic repository for type T
func NewRepository[T any](db *gorm.DB) Repository[T] {
	return &GORMRepository[T]{db: db}
}

func (r *GORMRepository[T]) Create(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *GORMRepository[T]) FindByID(ctx context.Context, id interface{}) (*T, error) {
	var entity T
	if err := r.db.WithContext(ctx).First(&entity, id).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *GORMRepository[T]) List(ctx context.Context, p Pagination) ([]T, int64, error) {
	var entities []T
	var total int64

	db := r.db.WithContext(ctx).Model(new(T))

	// Apply filters
	if len(p.Filters) > 0 {
		db = db.Where(p.Filters)
	}

	// Count total records (after filters)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	if p.Sort != "" {
		db = db.Order(p.Sort)
	}

	// Apply pagination
	err := db.Offset(p.GetOffset()).Limit(p.GetLimit()).Find(&entities).Error
	if err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}

func (r *GORMRepository[T]) Update(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

func (r *GORMRepository[T]) Delete(ctx context.Context, id interface{}) error {
	var entity T
	return r.db.WithContext(ctx).Delete(&entity, id).Error
}
