package grpc

import (
	"context"
	"fmt"

	pb "grouter/templates/rest-grpc-service/api/proto"
	"grouter/templates/rest-grpc-service/internal/pkg/business"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server implements the gRPC ItemService
type Server struct {
	pb.UnimplementedItemServiceServer
	logger      *zap.Logger
	itemService *business.ItemService
}

// NewServer creates a new gRPC server
func NewServer(logger *zap.Logger, itemService *business.ItemService) *Server {
	return &Server{
		logger:      logger,
		itemService: itemService,
	}
}

// CreateItem creates a new item
func (s *Server) CreateItem(ctx context.Context, req *pb.CreateItemRequest) (*pb.Item, error) {
	s.logger.Info("gRPC CreateItem called", zap.String("name", req.Name))

	item, err := s.itemService.CreateItem(ctx, req.Name, req.Description, req.Price, req.Quantity)
	if err != nil {
		s.logger.Error("failed to create item", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to create item: %v", err)
	}

	return &pb.Item{
		Id:          item.ID,
		Name:        item.Name,
		Description: item.Description,
		Price:       item.Price,
		Quantity:    item.Quantity,
	}, nil
}

// GetItem retrieves an item by ID
func (s *Server) GetItem(ctx context.Context, req *pb.GetItemRequest) (*pb.Item, error) {
	s.logger.Info("gRPC GetItem called", zap.String("id", req.Id))

	item, err := s.itemService.GetItem(ctx, req.Id)
	if err != nil {
		s.logger.Error("item not found", zap.String("id", req.Id), zap.Error(err))
		return nil, status.Errorf(codes.NotFound, "item not found: %s", req.Id)
	}

	return &pb.Item{
		Id:          item.ID,
		Name:        item.Name,
		Description: item.Description,
		Price:       item.Price,
		Quantity:    item.Quantity,
	}, nil
}

// ListItems lists all items
func (s *Server) ListItems(ctx context.Context, req *pb.ListItemsRequest) (*pb.ListItemsResponse, error) {
	s.logger.Info("gRPC ListItems called")

	items, err := s.itemService.ListItems(ctx)
	if err != nil {
		s.logger.Error("failed to list items", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to list items: %v", err)
	}

	pbItems := make([]*pb.Item, 0, len(items))
	for _, item := range items {
		pbItems = append(pbItems, &pb.Item{
			Id:          item.ID,
			Name:        item.Name,
			Description: item.Description,
			Price:       item.Price,
			Quantity:    item.Quantity,
		})
	}

	return &pb.ListItemsResponse{
		Items:      pbItems,
		TotalCount: int32(len(pbItems)),
	}, nil
}

// DeleteItem deletes an item by ID
func (s *Server) DeleteItem(ctx context.Context, req *pb.DeleteItemRequest) (*pb.DeleteItemResponse, error) {
	s.logger.Info("gRPC DeleteItem called", zap.String("id", req.Id))

	err := s.itemService.DeleteItem(ctx, req.Id)
	if err != nil {
		s.logger.Error("failed to delete item", zap.String("id", req.Id), zap.Error(err))
		return nil, status.Errorf(codes.NotFound, "item not found: %s", req.Id)
	}

	return &pb.DeleteItemResponse{
		Success: true,
		Message: fmt.Sprintf("item %s deleted successfully", req.Id),
	}, nil
}
