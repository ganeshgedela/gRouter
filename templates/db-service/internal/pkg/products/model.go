package products

import (
	"time"

	productv1 "grouter/templates/db-service/api/v1"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// Product represents a product entity in the database
type Product struct {
	ID          uint      `gorm:"primarykey"`
	Name        string    `gorm:"not null;index"`
	Description string    `gorm:"type:text"`
	Price       float64   `gorm:"not null;check:price >= 0"`
	SKU         string    `gorm:"uniqueIndex;not null"`
	Stock       int       `gorm:"default:0;check:stock >= 0"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

// TableName overrides the table name for GORM
func (Product) TableName() string {
	return "products"
}

// ToProto converts a Product model to a protobuf Product message
func (p *Product) ToProto() *productv1.Product {
	return &productv1.Product{
		Id:          uint64(p.ID),
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Sku:         p.SKU,
		Stock:       int32(p.Stock),
		CreatedAt:   timestamppb.New(p.CreatedAt),
		UpdatedAt:   timestamppb.New(p.UpdatedAt),
	}
}

// FromCreateRequest converts a CreateProductRequest to a Product model
func FromCreateRequest(req *productv1.CreateProductRequest) *Product {
	return &Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		SKU:         req.Sku,
		Stock:       int(req.Stock),
	}
}

// FromUpdateRequest updates a Product model from an UpdateProductRequest
func (p *Product) FromUpdateRequest(req *productv1.UpdateProductRequest) {
	if req.Name != "" {
		p.Name = req.Name
	}
	if req.Description != "" {
		p.Description = req.Description
	}
	if req.Price > 0 {
		p.Price = req.Price
	}
	if req.Stock >= 0 {
		p.Stock = int(req.Stock)
	}
}
