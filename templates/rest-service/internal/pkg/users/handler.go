package users

import (
	"context"
	"net/http"

	"grouter/pkg/manager"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Register factory on package import
func init() {
	manager.RegisterFactory("users", func(deps manager.Deps) (manager.Service, error) {
		return NewHandler(deps.Logger), nil
	})
}

// Handler handles user-related HTTP requests
type Handler struct {
	id     string
	logger *zap.Logger
	status manager.HealthStatus
}

// NewHandler creates a new user handler
func NewHandler(logger *zap.Logger) *Handler {
	return &Handler{
		id:     "users",
		logger: logger,
		status: manager.StatusCreated,
	}
}

// ID returns the service ID
func (h *Handler) ID() string { return h.id }

// Name returns the service name
func (h *Handler) Name() string { return "Users Handler" }

// Status returns current health status
func (h *Handler) Status() manager.HealthStatus { return h.status }

// LastError returns the last error
func (h *Handler) LastError() error { return nil }

// Dependencies returns list of service IDs this depends on
func (h *Handler) Dependencies() []string { return nil }

// Init initializes the service
func (h *Handler) Init(ctx context.Context) error {
	h.logger.Info("initializing users handler")
	h.status = manager.StatusInitialized
	return nil
}

// Start starts the service
func (h *Handler) Start(ctx context.Context) error {
	h.logger.Info("starting users handler")
	h.status = manager.StatusRunning
	return nil
}

// Stop stops the service
func (h *Handler) Stop(ctx context.Context) error {
	h.logger.Info("stopping users handler")
	h.status = manager.StatusStopped
	return nil
}

// RegisterRoutes registers user routes
func (h *Handler) RegisterRoutes(router gin.IRouter) {
	users := router.Group("/users")
	{
		users.GET("", h.ListUsers)
		users.GET("/:id", h.GetUser)
		users.POST("", h.CreateUser)
		users.PUT("/:id", h.UpdateUser)
		users.DELETE("/:id", h.DeleteUser)
	}
}

// ListUsers handles GET /users
// @Summary List users
// @Description Get all users
// @Tags users
// @Produce json
// @Success 200 {array} map[string]interface{}
// @Router /users [get]
func (h *Handler) ListUsers(c *gin.Context) {
	h.logger.Info("listing users")
	c.JSON(http.StatusOK, gin.H{
		"users": []map[string]interface{}{
			{"id": "1", "name": "Alice", "email": "alice@example.com"},
			{"id": "2", "name": "Bob", "email": "bob@example.com"},
		},
	})
}

// GetUser handles GET /users/:id
// @Summary Get user
// @Description Get user by ID
// @Tags users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} map[string]interface{}
// @Router /users/{id} [get]
func (h *Handler) GetUser(c *gin.Context) {
	id := c.Param("id")
	h.logger.Info("getting user", zap.String("id", id))
	c.JSON(http.StatusOK, gin.H{
		"id":    id,
		"name":  "Alice",
		"email": "alice@example.com",
	})
}

// CreateUser handles POST /users
// @Summary Create user
// @Description Create new user
// @Tags users
// @Accept json
// @Produce json
// @Success 201 {object} map[string]interface{}
// @Router /users [post]
func (h *Handler) CreateUser(c *gin.Context) {
	h.logger.Info("creating user")
	c.JSON(http.StatusCreated, gin.H{
		"id":      "3",
		"message": "User created successfully",
	})
}

// UpdateUser handles PUT /users/:id
// @Summary Update user
// @Description Update user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} map[string]interface{}
// @Router /users/{id} [put]
func (h *Handler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	h.logger.Info("updating user", zap.String("id", id))
	c.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
	})
}

// DeleteUser handles DELETE /users/:id
// @Summary Delete user
// @Description Delete user by ID
// @Tags users
// @Produce json
// @Param id path string true "User ID"
// @Success 204
// @Router /users/{id} [delete]
func (h *Handler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	h.logger.Info("deleting user", zap.String("id", id))
	c.Status(http.StatusNoContent)
}
