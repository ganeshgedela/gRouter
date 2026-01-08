package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	// Common API errors
	ErrBadRequest = APIError{
		StatusCode: http.StatusBadRequest,
		Error:      "BAD_REQUEST",
		Message:    "invalid request parameters",
		Code:       "BAD_REQUEST",
	}

	ErrUnauthorized = APIError{
		StatusCode: http.StatusUnauthorized,
		Error:      "UNAUTHORIZED",
		Message:    "authentication required",
		Code:       "UNAUTHORIZED",
	}

	ErrForbidden = APIError{
		StatusCode: http.StatusForbidden,
		Error:      "FORBIDDEN",
		Message:    "insufficient permissions",
		Code:       "FORBIDDEN",
	}

	ErrNotFound = APIError{
		StatusCode: http.StatusNotFound,
		Error:      "NOT_FOUND",
		Message:    "resource not found",
		Code:       "NOT_FOUND",
	}

	ErrInternalServer = APIError{
		StatusCode: http.StatusInternalServerError,
		Error:      "INTERNAL_ERROR",
		Message:    "internal server error",
		Code:       "INTERNAL_ERROR",
	}
)

// ErrorResponse represents a standardized API error response
type ErrorResponse struct {
	Error     string `json:"error"`
	Message   string `json:"message"`
	Code      string `json:"code"`
	RequestID string `json:"request_id,omitempty"`
}

// APIError represents an API error
type APIError struct {
	StatusCode int
	Error      string
	Message    string
	Code       string
}

// RespondError sends a standardized error response
func RespondError(c *gin.Context, apiErr APIError) {
	requestID := c.GetString("request_id")
	response := ErrorResponse{
		Error:     apiErr.Error,
		Message:   apiErr.Message,
		Code:      apiErr.Code,
		RequestID: requestID,
	}
	c.JSON(apiErr.StatusCode, response)
}
