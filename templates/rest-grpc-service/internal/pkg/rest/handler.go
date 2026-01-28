package rest

import (
	"net/http"

	"grouter/templates/rest-grpc-service/internal/pkg/business"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Handler handles REST API requests
type Handler struct {
	logger      *zap.Logger
	itemService *business.ItemService
}

// NewHandler creates a new REST handler
func NewHandler(logger *zap.Logger, itemService *business.ItemService) *Handler {
	return &Handler{
		logger:      logger,
		itemService: itemService,
	}
}

// RegisterRoutes registers REST API routes
func (h *Handler) RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api/v1")
	{
		api.GET("/health", h.Health)
		api.POST("/items", h.CreateItem)
		api.GET("/items", h.ListItems)
		api.GET("/items/:id", h.GetItem)
		api.DELETE("/items/:id", h.DeleteItem)
	}
}

// CreateItemRequest represents the create item request
type CreateItemRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Quantity    int32   `json:"quantity" binding:"required,gte=0"`
}

// CreateItem creates a new item
func (h *Handler) CreateItem(c *gin.Context) {
	var req CreateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("REST CreateItem called", zap.String("name", req.Name))

	item, err := h.itemService.CreateItem(c.Request.Context(), req.Name, req.Description, req.Price, req.Quantity)
	if err != nil {
		h.logger.Error("failed to create item", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create item"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":          item.ID,
		"name":        item.Name,
		"description": item.Description,
		"price":       item.Price,
		"quantity":    item.Quantity,
	})
}

// GetItem retrieves an item by ID
func (h *Handler) GetItem(c *gin.Context) {
	id := c.Param("id")
	h.logger.Info("REST GetItem called", zap.String("id", id))

	item, err := h.itemService.GetItem(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("item not found", zap.String("id", id), zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":          item.ID,
		"name":        item.Name,
		"description": item.Description,
		"price":       item.Price,
		"quantity":    item.Quantity,
	})
}

// ListItems lists all items
func (h *Handler) ListItems(c *gin.Context) {
	h.logger.Info("REST ListItems called")

	items, err := h.itemService.ListItems(c.Request.Context())
	if err != nil {
		h.logger.Error("failed to list items", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list items"})
		return
	}

	response := make([]gin.H, 0, len(items))
	for _, item := range items {
		response = append(response, gin.H{
			"id":          item.ID,
			"name":        item.Name,
			"description": item.Description,
			"price":       item.Price,
			"quantity":    item.Quantity,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"items": response,
		"count": len(response),
	})
}

// DeleteItem deletes an item by ID
func (h *Handler) DeleteItem(c *gin.Context) {
	id := c.Param("id")
	h.logger.Info("REST DeleteItem called", zap.String("id", id))

	err := h.itemService.DeleteItem(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("failed to delete item", zap.String("id", id), zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "item deleted successfully",
	})
}

// Health returns the health status
func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":      "healthy",
		"service":     "rest-grpc-service",
		"items_count": h.itemService.Count(),
	})
}
