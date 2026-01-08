package interceptors

import (
	"context"
	"sync"
	"time"

	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RateLimiter creates a server interceptor that rate limits requests using token bucket.
func RateLimiter(rps int, burst int) grpc.UnaryServerInterceptor {
	limiter := rate.NewLimiter(rate.Limit(rps), burst)

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if !limiter.Allow() {
			return nil, status.Errorf(codes.ResourceExhausted, "rate limit exceeded")
		}

		return handler(ctx, req)
	}
}

// PerMethodRateLimiter creates a server interceptor with per-method rate limits.
type MethodRateLimits map[string]*rate.Limiter

// NewMethodRateLimits creates a new per-method rate limiter.
func NewMethodRateLimits() MethodRateLimits {
	return make(MethodRateLimits)
}

// Add adds a rate limit for a specific method.
func (m MethodRateLimits) Add(method string, rps int, burst int) {
	m[method] = rate.NewLimiter(rate.Limit(rps), burst)
}

// PerMethodRateLimiter creates an interceptor with per-method rate limiting.
func PerMethodRateLimiterInterceptor(limits MethodRateLimits, defaultRPS int, defaultBurst int) grpc.UnaryServerInterceptor {
	defaultLimiter := rate.NewLimiter(rate.Limit(defaultRPS), defaultBurst)

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		limiter, ok := limits[info.FullMethod]
		if !ok {
			limiter = defaultLimiter
		}

		if !limiter.Allow() {
			return nil, status.Errorf(codes.ResourceExhausted, "rate limit exceeded for method %s", info.FullMethod)
		}

		return handler(ctx, req)
	}
}

// CircuitBreakerState represents the state of a circuit breaker.
type CircuitBreakerState int

const (
	StateClosed CircuitBreakerState = iota
	StateOpen
	StateHalfOpen
)

// CircuitBreaker implements the circuit breaker pattern for client calls.
type CircuitBreaker struct {
	mu sync.RWMutex

	threshold   int           // Failures before opening
	timeout     time.Duration // Time before half-open
	maxRequests int           // Max requests in half-open

	state         CircuitBreakerState
	failures      int
	successes     int
	requests      int
	lastStateTime time.Time
}

// NewCircuitBreaker creates a new circuit breaker.
func NewCircuitBreaker(threshold int, timeout time.Duration, maxRequests int) *CircuitBreaker {
	return &CircuitBreaker{
		threshold:     threshold,
		timeout:       timeout,
		maxRequests:   maxRequests,
		state:         StateClosed,
		lastStateTime: time.Now(),
	}
}

// Call executes the function if the circuit breaker allows it.
func (cb *CircuitBreaker) Call(fn func() error) error {
	if !cb.allowRequest() {
		return status.Errorf(codes.Unavailable, "circuit breaker is open")
	}

	err := fn()
	cb.recordResult(err)
	return err
}

func (cb *CircuitBreaker) allowRequest() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case StateClosed:
		return true
	case StateOpen:
		if time.Since(cb.lastStateTime) > cb.timeout {
			cb.state = StateHalfOpen
			cb.requests = 0
			cb.lastStateTime = time.Now()
			return true
		}
		return false
	case StateHalfOpen:
		return cb.requests < cb.maxRequests
	}
	return false
}

func (cb *CircuitBreaker) recordResult(err error) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if err != nil {
		cb.failures++
		cb.successes = 0

		if cb.state == StateHalfOpen {
			cb.state = StateOpen
			cb.lastStateTime = time.Now()
		} else if cb.failures >= cb.threshold {
			cb.state = StateOpen
			cb.lastStateTime = time.Now()
		}
	} else {
		cb.successes++
		cb.failures = 0

		if cb.state == StateHalfOpen && cb.successes >= cb.maxRequests {
			cb.state = StateClosed
			cb.lastStateTime = time.Now()
		}
	}

	if cb.state == StateHalfOpen {
		cb.requests++
	}
}

// CircuitBreakerInterceptor creates a client interceptor with circuit breaker.
func CircuitBreakerInterceptor(cb *CircuitBreaker) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		return cb.Call(func() error {
			return invoker(ctx, method, req, reply, cc, opts...)
		})
	}
}
