package grpc

import (
	"context"
	"fmt"

	commonv1 "grouter/templates/db-service/api/v1"
	productv1 "grouter/templates/db-service/api/v1"
	"grouter/templates/db-service/internal/pkg/products"

	"grouter/pkg/database"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

// ProductServer implements the ProductService gRPC server
type ProductServer struct {
	productv1.UnimplementedProductServiceServer
	repo   *products.Repository
	logger *zap.Logger
}

// NewProductServer creates a new product gRPC server
func NewProductServer(repo *products.Repository, logger *zap.Logger) *ProductServer {
	return &ProductServer{
		repo:   repo,
		logger: logger,
	}
}

// CreateProduct creates a new product
func (s *ProductServer) CreateProduct(ctx context.Context, req *productv1.CreateProductRequest) (*productv1.Product, error) {
	s.logger.Info("creating product", zap.String("sku", req.Sku), zap.String("name", req.Name))

	// Validate request
	if err := s.validateCreateRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Check if SKU already exists
	existing, err := s.repo.FindBySKU(ctx, req.Sku)
	if err == nil && existing != nil {
		return nil, status.Errorf(codes.AlreadyExists, "product with SKU %s already exists", req.Sku)
	}

	// Convert request to model
	product := products.FromCreateRequest(req)

	// Create product
	if err := s.repo.Create(ctx, product); err != nil {
		s.logger.Error("failed to create product", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to create product")
	}

	s.logger.Info("product created successfully", zap.Uint("id", product.ID))
	return product.ToProto(), nil
}

// GetProduct retrieves a product by ID
func (s *ProductServer) GetProduct(ctx context.Context, req *productv1.GetProductRequest) (*productv1.Product, error) {
	s.logger.Info("getting product", zap.Uint64("id", req.Id))

	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	product, err := s.repo.FindByID(ctx, uint(req.Id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Errorf(codes.NotFound, "product with id %d not found", req.Id)
		}
		s.logger.Error("failed to get product", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to get product")
	}

	return product.ToProto(), nil
}

// ListProducts lists products with pagination and filtering
func (s *ProductServer) ListProducts(ctx context.Context, req *productv1.ListProductsRequest) (*productv1.ListProductsResponse, error) {
	s.logger.Info("listing products")

	// Convert proto pagination to database pagination
	dbPagination := protoToDBPagination(req.Pagination)

	// List products
	productList, total, err := s.repo.List(ctx, dbPagination)
	if err != nil {
		s.logger.Error("failed to list products", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to list products")
	}

	// Convert to proto
	protoProducts := make([]*productv1.Product, len(productList))
	for i, p := range productList {
		protoProducts[i] = p.ToProto()
	}

	// Calculate total pages
	totalPages := int32(total) / int32(dbPagination.PageSize)
	if int32(total)%int32(dbPagination.PageSize) > 0 {
		totalPages++
	}

	return &productv1.ListProductsResponse{
		Products: protoProducts,
		Pagination: &commonv1.PaginationResponse{
			Total:      total,
			Page:       int32(dbPagination.Page),
			PageSize:   int32(dbPagination.PageSize),
			TotalPages: totalPages,
		},
	}, nil
}

// UpdateProduct updates an existing product
func (s *ProductServer) UpdateProduct(ctx context.Context, req *productv1.UpdateProductRequest) (*productv1.Product, error) {
	s.logger.Info("updating product", zap.Uint64("id", req.Id))

	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	// Get existing product
	product, err := s.repo.FindByID(ctx, uint(req.Id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Errorf(codes.NotFound, "product with id %d not found", req.Id)
		}
		s.logger.Error("failed to get product", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to get product")
	}

	// Update fields
	product.FromUpdateRequest(req)

	// Save changes
	if err := s.repo.Update(ctx, product); err != nil {
		s.logger.Error("failed to update product", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to update product")
	}

	s.logger.Info("product updated successfully", zap.Uint("id", product.ID))
	return product.ToProto(), nil
}

// DeleteProduct deletes a product by ID
func (s *ProductServer) DeleteProduct(ctx context.Context, req *productv1.DeleteProductRequest) (*emptypb.Empty, error) {
	s.logger.Info("deleting product", zap.Uint64("id", req.Id))

	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	// Check if product exists
	_, err := s.repo.FindByID(ctx, uint(req.Id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Errorf(codes.NotFound, "product with id %d not found", req.Id)
		}
		s.logger.Error("failed to check product", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to check product")
	}

	// Delete product
	if err := s.repo.Delete(ctx, uint(req.Id)); err != nil {
		s.logger.Error("failed to delete product", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to delete product")
	}

	s.logger.Info("product deleted successfully", zap.Uint64("id", req.Id))
	return &emptypb.Empty{}, nil
}

// UpdateStock atomically updates product stock
func (s *ProductServer) UpdateStock(ctx context.Context, req *productv1.UpdateStockRequest) (*productv1.Product, error) {
	s.logger.Info("updating stock", zap.Uint64("id", req.Id), zap.Int32("delta", req.StockDelta))

	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	// Update stock atomically
	product, err := s.repo.UpdateStock(ctx, uint(req.Id), int(req.StockDelta))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Errorf(codes.NotFound, "product with id %d not found", req.Id)
		}
		s.logger.Error("failed to update stock", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	s.logger.Info("stock updated successfully", zap.Uint("id", product.ID), zap.Int("new_stock", product.Stock))
	return product.ToProto(), nil
}

// BulkCreateProducts creates multiple products using server streaming
func (s *ProductServer) BulkCreateProducts(stream productv1.ProductService_BulkCreateProductsServer) error {
	s.logger.Info("bulk creating products")

	var productsToCreate []*products.Product

	// Receive all products from stream
	for {
		req, err := stream.Recv()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			s.logger.Error("error receiving from stream", zap.Error(err))
			return status.Error(codes.Internal, "error receiving products")
		}

		// Validate request
		if err := s.validateCreateRequest(req); err != nil {
			return status.Error(codes.InvalidArgument, err.Error())
		}

		product := products.FromCreateRequest(req)
		productsToCreate = append(productsToCreate, product)
	}

	if len(productsToCreate) == 0 {
		return status.Error(codes.InvalidArgument, "no products to create")
	}

	// Bulk create in transaction
	if err := s.repo.BulkCreate(stream.Context(), productsToCreate); err != nil {
		s.logger.Error("failed to bulk create products", zap.Error(err))
		return status.Error(codes.Internal, "failed to create products")
	}

	// Convert to proto
	protoProducts := make([]*productv1.Product, len(productsToCreate))
	for i, p := range productsToCreate {
		protoProducts[i] = p.ToProto()
	}

	s.logger.Info("bulk created products", zap.Int("count", len(productsToCreate)))
	return stream.SendAndClose(&productv1.BulkCreateProductsResponse{
		CreatedCount: int32(len(productsToCreate)),
		Products:     protoProducts,
	})
}

// Helper functions

func (s *ProductServer) validateCreateRequest(req *productv1.CreateProductRequest) error {
	if req.Name == "" {
		return fmt.Errorf("name is required")
	}
	if req.Sku == "" {
		return fmt.Errorf("sku is required")
	}
	if req.Price < 0 {
		return fmt.Errorf("price must be non-negative")
	}
	if req.Stock < 0 {
		return fmt.Errorf("stock must be non-negative")
	}
	return nil
}

func protoToDBPagination(p *commonv1.PaginationRequest) database.Pagination {
	if p == nil {
		return database.Pagination{
			Page:     1,
			PageSize: 10,
		}
	}

	page := int(p.Page)
	if page < 1 {
		page = 1
	}

	pageSize := int(p.PageSize)
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100 // Max limit
	}

	filters := make(map[string]interface{})
	for k, v := range p.Filters {
		filters[k] = v
	}

	return database.Pagination{
		Page:     page,
		PageSize: pageSize,
		Sort:     p.Sort,
		Filters:  filters,
	}
}
