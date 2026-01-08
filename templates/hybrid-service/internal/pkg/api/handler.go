package api

import (
	"context"
	"fmt"
	"net/http"

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

// Handler handles HTTP requests and publishes events
type Handler struct {
	id        string
	logger    *zap.Logger
	messenger *nats.Messenger
	status    manager.HealthStatus
}

// NewHandler creates a new API handler
func NewHandler(logger *zap.Logger, messenger *nats.Messenger) *Handler {
	return &Handler{
		id:        "api",
		logger:    logger,
		messenger: messenger,
		status:    manager.StatusCreated,
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
	api := router.Group("/api")
	{
		api.POST("/orders", h.CreateOrder)
		api.GET("/health", h.Health)
	}
}

// CreateOrder creates an order and publishes an event
// @Summary Create order
// @Description Create new order and publish event to NATS
// @Tags orders
// @Accept json
// @Produce json
// @Success 201 {object} map[string]interface{}
// @Router /api/orders [post]
func (h *Handler) CreateOrder(c *gin.Context) {
	h.logger.Info("creating order via HTTP")

	orderData := map[string]interface{}{
		"order_id": "order-123",
		"user_id":  "user-456",
		"amount":   99.99,
	}

	// Publish event to NATS
	if h.messenger != nil {
		opts := &nats.PublishOptions{}
		err := h.messenger.Publish(context.Background(), "orders.created", "order.created", orderData, opts)
		if err != nil {
			h.logger.Error("failed to publish event", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish event"})
			return
		}
		h.logger.Info("published order.created event")
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Order created and event published",
		"order":   orderData,
	})
}

// Health handles health check
func (h *Handler) Health(c *gin.Context) {
	status := "healthy"
	if h.messenger != nil && !h.messenger.IsConnected() {
		status = "degraded - NATS disconnected"
	}

	c.JSON(http.StatusOK, gin.H{
		"status": status,
	})
}
