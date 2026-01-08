package web

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimiter manages rate limiting per client
type RateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	rps      float64
	burst    int
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(rps float64, burst int) *RateLimiter {
	return &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
		rps:      rps,
		burst:    burst,
	}
}

// getLimiter returns a rate limiter for a client
func (rl *RateLimiter) getLimiter(clientID string) *rate.Limiter {
	rl.mu.RLock()
	limiter, exists := rl.limiters[clientID]
	rl.mu.RUnlock()

	if exists {
		return limiter
	}

	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Double-check after acquiring write lock
	if limiter, exists := rl.limiters[clientID]; exists {
		return limiter
	}

	limiter = rate.NewLimiter(rate.Limit(rl.rps), rl.burst)
	rl.limiters[clientID] = limiter

	return limiter
}

// RateLimitMiddleware creates rate limiting middleware
func RateLimitMiddleware(rps float64, burst int) gin.HandlerFunc {
	limiter := NewRateLimiter(rps, burst)

	return func(c *gin.Context) {
		// Use client IP as identifier
		clientID := c.ClientIP()

		clientLimiter := limiter.getLimiter(clientID)

		if !clientLimiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, ErrorResponse{
				Error:   "RATE_LIMIT_EXCEEDED",
				Message: "too many requests",
				Code:    "TOO_MANY_REQUESTS",
			})
			return
		}

		c.Next()
	}
}
