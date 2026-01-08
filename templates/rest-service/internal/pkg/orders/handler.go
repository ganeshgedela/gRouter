package orders

import (
	"context"
	"net/http"

	"grouter/pkg/manager"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Register factory on package import
func init() {
	manager.RegisterFactory("orders", func(deps manager.Deps) (manager.Service, error) {
		return NewHandler(deps.Logger), nil
	})
}

// Handler handles order-related HTTP requests
type Handler struct {
	id     string
	logger *zap.Logger
	status manager.HealthStatus
}

// NewHandler creates a new order handler
func NewHandler(logger *zap.Logger) *Handler {
	return &Handler{
		id:     "orders",
		logger: logger,
		status: manager.StatusCreated,
	}
}

// ID returns the service ID
func (h *Handler) ID() string { return h.id }

// Name returns the service name
func (h *Handler) Name() string { return "Orders Handler" }

// Status returns current health status
func (h *Handler) Status() manager.HealthStatus { return h.status }

// LastError returns the last error
func (h *Handler) LastError() error { return nil }

// Dependencies returns list of service IDs this depends on
func (h *Handler) Dependencies() []string { return nil }

// Init initializes the service
func (h *Handler) Init(ctx context.Context) error {
	h.logger.Info("initializing orders handler")
	h.status = manager.StatusInitialized
	return nil
}

// Start starts the service
func (h *Handler) Start(ctx context.Context) error {
	h.logger.Info("starting orders handler")
	h.status = manager.StatusRunning
	return nil
}

// Stop stops the service
func (h *Handler) Stop(ctx context.Context) error {
	h.logger.Info("stopping orders handler")
	h.status = manager.StatusStopped
	return nil
}

// RegisterRoutes registers order routes
func (h *Handler) RegisterRoutes(router gin.IRouter) {
	orders := router.Group("/orders")
	{
		orders.GET("", h.ListOrders)
		orders.GET("/:id", h.GetOrder)
		orders.POST("", h.CreateOrder)
		orders.PUT("/:id/status", h.UpdateOrderStatus)
	}
}

// ListOrders handles GET /orders
// @Summary List orders
// @Description Get all orders
// @Tags orders
// @Produce json
// @Success 200 {array} map[string]interface{}
// @Router /orders [get]
func (h *Handler) ListOrders(c *gin.Context) {
	h.logger.Info("listing orders")
	c.JSON(http.StatusOK, gin.H{
		"orders": []map[string]interface{}{
			{"id": "1", "user_id": "1", "amount": 99.99, "status": "completed"},
			{"id": "2", "user_id": "2", "amount": 149.99, "status": "pending"},
		},
	})
}

// GetOrder handles GET /orders/:id
// @Summary Get order
// @Description Get order by ID
// @Tags orders
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} map[string]interface{}
// @Router /orders/{id} [get]
func (h *Handler) GetOrder(c *gin.Context) {
	id := c.Param("id")
	h.logger.Info("getting order", zap.String("id", id))
	c.JSON(http.StatusOK, gin.H{
		"id":      id,
		"user_id": "1",
		"amount":  99.99,
		"status":  "completed",
	})
}

// CreateOrder handles POST /orders
// @Summary Create order
// @Description Create new order
// @Tags orders
// @Accept json
// @Produce json
// @Success 201 {object} map[string]interface{}
// @Router /orders [post]
func (h *Handler) CreateOrder(c *gin.Context) {
	h.logger.Info("creating order")
	c.JSON(http.StatusCreated, gin.H{
		"id":      "3",
		"message": "Order created successfully",
	})
}

// UpdateOrderStatus handles PUT /orders/:id/status
// @Summary Update order status
// @Description Update order status by ID
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} map[string]interface{}
// @Router /orders/{id}/status [put]
func (h *Handler) UpdateOrderStatus(c *gin.Context) {
	id := c.Param("id")
	h.logger.Info("updating order status", zap.String("id", id))
	c.JSON(http.StatusOK, gin.H{
		"message": "Order status updated successfully",
	})
}
