package health

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNewHealthService(t *testing.T) {
	s := NewHealthService()
	assert.NotNil(t, s)
	assert.NotNil(t, s.readiness)
	assert.NotNil(t, s.liveness)
}

func TestHealthService_AddLivenessCheck(t *testing.T) {
	s := NewHealthService()
	s.AddLivenessCheck("test", func() error { return nil })

	checks, err := s.CheckLiveness()
	assert.NoError(t, err)
	assert.Equal(t, "OK", checks["test"])
}

func TestHealthService_AddReadinessCheck(t *testing.T) {
	s := NewHealthService()
	s.AddReadinessCheck("test", func() error { return nil })

	checks, err := s.CheckReadiness()
	assert.NoError(t, err)
	assert.Equal(t, "OK", checks["test"])
}

func TestHealthService_CheckLiveness_Failure(t *testing.T) {
	s := NewHealthService()
	s.AddLivenessCheck("fail", func() error { return errors.New("failed") })
	s.AddLivenessCheck("pass", func() error { return nil })

	checks, err := s.CheckLiveness()
	assert.Error(t, err)
	assert.Equal(t, "failed", checks["fail"])
	assert.Equal(t, "OK", checks["pass"])
}

func TestHealthService_CheckReadiness_Failure(t *testing.T) {
	s := NewHealthService()
	s.AddReadinessCheck("fail", func() error { return errors.New("failed") })
	s.AddReadinessCheck("pass", func() error { return nil })

	checks, err := s.CheckReadiness()
	assert.Error(t, err)
	assert.Equal(t, "failed", checks["fail"])
	assert.Equal(t, "OK", checks["pass"])
}

func TestHealthService_ConcurrentAccess(t *testing.T) {
	s := NewHealthService()
	var wg sync.WaitGroup
	count := 100

	// Concurrent writes
	for i := 0; i < count; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			s.AddLivenessCheck("check", func() error { return nil })
		}(i)
	}

	// Concurrent reads
	for i := 0; i < count; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = s.CheckLiveness()
		}()
	}

	wg.Wait()
}

func TestLivenessHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	s := NewHealthService()
	s.AddLivenessCheck("test", func() error { return nil })

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	s.LivenessHandler(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "up", resp["status"])
}

func TestLivenessHandler_Fail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	s := NewHealthService()
	s.AddLivenessCheck("test", func() error { return errors.New("oops") })

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	s.LivenessHandler(c)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "down", resp["status"])
}

func TestReadinessHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	s := NewHealthService()
	s.AddReadinessCheck("test", func() error { return nil })

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	s.ReadinessHandler(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "ready", resp["status"])
}

func TestHealthService_RemoveLivenessCheck(t *testing.T) {
	s := NewHealthService()
	s.AddLivenessCheck("test", func() error { return nil })

	// Verify it exists
	checks, _ := s.CheckLiveness()
	assert.Contains(t, checks, "test")

	// Remove it
	s.RemoveLivenessCheck("test")

	// Verify it's gone
	checks, _ = s.CheckLiveness()
	assert.NotContains(t, checks, "test")
}

func TestHealthService_RemoveReadinessCheck(t *testing.T) {
	s := NewHealthService()
	s.AddReadinessCheck("test", func() error { return nil })

	// Verify it exists
	checks, _ := s.CheckReadiness()
	assert.Contains(t, checks, "test")

	// Remove it
	s.RemoveReadinessCheck("test")

	// Verify it's gone
	checks, _ = s.CheckReadiness()
	assert.NotContains(t, checks, "test")
}
