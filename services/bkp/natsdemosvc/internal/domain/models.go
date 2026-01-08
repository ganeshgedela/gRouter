package domain

// EchoRequest represents the request body for EchoHandler
type EchoRequest struct {
	Message string `json:"message"`
}
