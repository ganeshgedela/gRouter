package api

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"grouter/pkg/manager"
	"grouter/pkg/messaging/nats"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Register factory on package import
func init() {
	manager.RegisterFactory("api", func(deps manager.Deps) (manager.Service, error) {
		if deps.Messenger == nil {
			return nil, fmt.Errorf("messenger not provided in dependencies")
		}
		return NewHandler(deps.Logger, deps.Messenger), nil
	})
}

// Item represents an API item
type Item struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description,omitempty"`
	Price       float64 `json:"price"`
}

// Handler handles HTTP requests and publishes events
type Handler struct {
	id        string
	logger    *zap.Logger
	messenger *nats.Messenger
	status    manager.HealthStatus
	items     map[string]Item
	mu        sync.RWMutex
}

// NewHandler creates a new API handler
func NewHandler(logger *zap.Logger, messenger *nats.Messenger) *Handler {
	return &Handler{
		id:        "api",
		logger:    logger,
		messenger: messenger,
		status:    manager.StatusCreated,
		items:     make(map[string]Item),
	}
}

// ID returns the service ID
func (h *Handler) ID() string { return h.id }

// Name returns the service name
func (h *Handler) Name() string { return "API Handler" }

// Status returns current health status
func (h *Handler) Status() manager.HealthStatus { return h.status }

// LastError returns the last error
func (h *Handler) LastError() error { return nil }

// Dependencies returns list of service IDs this depends on
func (h *Handler) Dependencies() []string { return nil }

// Init initializes the service
func (h *Handler) Init(ctx context.Context) error {
	h.logger.Info("initializing API handler")
	h.status = manager.StatusInitialized
	return nil
}

// Start starts the service
func (h *Handler) Start(ctx context.Context) error {
	h.logger.Info("starting API handler")
	h.status = manager.StatusRunning
	return nil
}

// Stop stops the service
func (h *Handler) Stop(ctx context.Context) error {
	h.logger.Info("stopping API handler")
	h.status = manager.StatusStopped
	return nil
}

// RegisterRoutes registers API routes
func (h *Handler) RegisterRoutes(router gin.IRouter) {
	v1 := router.Group("/api/v1")
	{
		// Items endpoints
		v1.GET("/items", h.ListItems)
		v1.POST("/items", h.CreateItem)
		v1.GET("/items/:id", h.GetItem)
		v1.DELETE("/items/:id", h.DeleteItem)

		// Health endpoint
		v1.GET("/health", h.Health)
	}
}

// ListItems returns all items
func (h *Handler) ListItems(c *gin.Context) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	items := make([]Item, 0, len(h.items))
	for _, item := range h.items {
		items = append(items, item)
	}

	c.JSON(http.StatusOK, gin.H{
		"items": items,
		"count": len(items),
	})
}

// CreateItem creates a new item
func (h *Handler) CreateItem(c *gin.Context) {
	var item Item
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate ID if not provided
	if item.ID == "" {
		item.ID = fmt.Sprintf("item-%d", len(h.items)+1)
	}

	h.mu.Lock()
	h.items[item.ID] = item
	h.mu.Unlock()

	h.logger.Info("created item", zap.String("id", item.ID))

	// Publish event to NATS
	if h.messenger != nil {
		ctx := c.Request.Context()
		eventData := map[string]interface{}{
			"item_id":     item.ID,
			"name":        item.Name,
			"description": item.Description,
			"price":       item.Price,
		}
		if err := h.messenger.Publish(ctx, "item.created", "item.created", eventData, &nats.PublishOptions{}); err != nil {
			h.logger.Error("failed to publish item.created event", zap.Error(err))
		} else {
			h.logger.Info("published item.created event to NATS", zap.String("item_id", item.ID))
		}
	}

	c.JSON(http.StatusCreated, item)
}

// GetItem returns a specific item
func (h *Handler) GetItem(c *gin.Context) {
	id := c.Param("id")

	h.mu.RLock()
	item, exists := h.items[id]
	h.mu.RUnlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
		return
	}

	c.JSON(http.StatusOK, item)
}

// DeleteItem deletes an item
func (h *Handler) DeleteItem(c *gin.Context) {
	id := c.Param("id")

	h.mu.Lock()
	_, exists := h.items[id]
	if exists {
		delete(h.items, id)
	}
	h.mu.Unlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
		return
	}

	h.logger.Info("deleted item", zap.String("id", id))

	// Publish event to NATS
	if h.messenger != nil {
		ctx := c.Request.Context()
		eventData := map[string]interface{}{"item_id": id}
		if err := h.messenger.Publish(ctx, "item.deleted", "item.deleted", eventData, &nats.PublishOptions{}); err != nil {
			h.logger.Error("failed to publish item.deleted event", zap.Error(err))
		} else {
			h.logger.Info("published item.deleted event to NATS", zap.String("item_id", id))
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "item deleted"})
}

// Health returns API health status
func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "hybrid-api",
		"items":   len(h.items),
	})
}

// WebRoutable interface implementation
func (h *Handler) IsWebRoutable() bool {
	return true
}

func (h *Handler) GetRouteRegistrar() func(router interface{}) {
	return func(router interface{}) {
		if r, ok := router.(gin.IRouter); ok {
			h.RegisterRoutes(r)
		}
	}
}
