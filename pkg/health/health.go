package health

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

// HealthChecker is a function that returns an error if the check fails
type HealthChecker func() error

// HealthService manages health checks
type HealthService struct {
	mu        sync.RWMutex
	readiness map[string]HealthChecker
	liveness  map[string]HealthChecker
}

// NewHealthService creates a new HealthService
func NewHealthService() *HealthService {
	return &HealthService{
		readiness: make(map[string]HealthChecker),
		liveness:  make(map[string]HealthChecker),
	}
}

// AddReadinessCheck adds a readiness check
func (s *HealthService) AddReadinessCheck(name string, check HealthChecker) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.readiness[name] = check
}

// AddLivenessCheck adds a liveness check
func (s *HealthService) AddLivenessCheck(name string, check HealthChecker) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.liveness[name] = check
}

// RemoveReadinessCheck removes a readiness check
func (s *HealthService) RemoveReadinessCheck(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.readiness, name)
}

// RemoveLivenessCheck removes a liveness check
func (s *HealthService) RemoveLivenessCheck(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.liveness, name)
}

// CheckLiveness performs all liveness checks
func (s *HealthService) CheckLiveness() (map[string]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	errors := make(map[string]string)
	hasError := false

	for name, check := range s.liveness {
		if err := check(); err != nil {
			errors[name] = err.Error()
			hasError = true
		} else {
			errors[name] = "OK"
		}
	}

	if hasError {
		return errors, fmt.Errorf("liveness check failed")
	}
	return errors, nil
}

// CheckReadiness performs all readiness checks
func (s *HealthService) CheckReadiness() (map[string]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	errors := make(map[string]string)
	hasError := false

	for name, check := range s.readiness {
		if err := check(); err != nil {
			errors[name] = err.Error()
			hasError = true
		} else {
			errors[name] = "OK"
		}
	}

	if hasError {
		return errors, fmt.Errorf("readiness check failed")
	}
	return errors, nil
}

// LivenessHandler handles liveness probes
func (s *HealthService) LivenessHandler(c *gin.Context) {
	checks, err := s.CheckLiveness()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "down",
			"checks": checks,
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "up",
		"checks": checks,
	})
}

// ReadinessHandler handles readiness probes
func (s *HealthService) ReadinessHandler(c *gin.Context) {
	checks, err := s.CheckReadiness()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "not ready",
			"checks": checks,
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
		"checks": checks,
	})
}
