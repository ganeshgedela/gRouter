package business

import (
	"context"
	"fmt"
	"sync"

	"go.uber.org/zap"
)

// Item represents a business item
type Item struct {
	ID          string
	Name        string
	Description string
	Price       float64
	Quantity    int32
}

// ItemService manages item business logic
type ItemService struct {
	logger *zap.Logger
	items  map[string]*Item
	mu     sync.RWMutex
	nextID int
}

// NewItemService creates a new item service
func NewItemService(logger *zap.Logger) *ItemService {
	return &ItemService{
		logger: logger,
		items:  make(map[string]*Item),
		nextID: 1,
	}
}

// CreateItem creates a new item
func (s *ItemService) CreateItem(ctx context.Context, name, description string, price float64, quantity int32) (*Item, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := fmt.Sprintf("item-%d", s.nextID)
	s.nextID++

	item := &Item{
		ID:          id,
		Name:        name,
		Description: description,
		Price:       price,
		Quantity:    quantity,
	}

	s.items[id] = item
	s.logger.Info("created item",
		zap.String("id", id),
		zap.String("name", name),
		zap.Float64("price", price))

	return item, nil
}

// GetItem retrieves an item by ID
func (s *ItemService) GetItem(ctx context.Context, id string) (*Item, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item, exists := s.items[id]
	if !exists {
		return nil, fmt.Errorf("item not found: %s", id)
	}

	return item, nil
}

// ListItems returns all items
func (s *ItemService) ListItems(ctx context.Context) ([]*Item, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	items := make([]*Item, 0, len(s.items))
	for _, item := range s.items {
		items = append(items, item)
	}

	return items, nil
}

// DeleteItem deletes an item by ID
func (s *ItemService) DeleteItem(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.items[id]; !exists {
		return fmt.Errorf("item not found: %s", id)
	}

	delete(s.items, id)
	s.logger.Info("deleted item", zap.String("id", id))

	return nil
}

// Count returns the total number of items
func (s *ItemService) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.items)
}
