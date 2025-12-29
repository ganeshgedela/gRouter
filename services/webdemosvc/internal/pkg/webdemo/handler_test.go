package webdemo

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestService_InterfaceCheck(t *testing.T) {
	s := NewService()

	// Test Name
	assert.Equal(t, "webdemo", s.Name())

	// Test Lifecycle methods (simple no-ops currently)
	ctx := context.Background()
	assert.NoError(t, s.Ready(ctx))
	assert.NoError(t, s.Start(ctx))
	assert.NoError(t, s.Stop(ctx))
}

func TestService_RegisterRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	s := NewService()
	r := gin.New()

	// Create a group and register
	g := r.Group("/")
	s.RegisterRoutes(g)

	// Verify routes are registered by making requests

	// Hello
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/hello", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Echo
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/echo?msg=test", nil)
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestService_HelloHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	s := NewService()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	s.HelloHandler(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "Hello from WebDemoSvc!", resp["message"])
}

func TestService_EchoHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	s := NewService()

	// Case 1: With query param
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/echo?msg=test", nil)

	s.EchoHandler(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "test", resp["echo"])

	// Case 2: Without query param (default)
	w2 := httptest.NewRecorder()
	c2, _ := gin.CreateTestContext(w2)
	c2.Request, _ = http.NewRequest("GET", "/echo", nil)

	s.EchoHandler(c2)

	assert.Equal(t, http.StatusOK, w2.Code)

	var resp2 map[string]string
	err = json.Unmarshal(w2.Body.Bytes(), &resp2)
	assert.NoError(t, err)
	assert.Equal(t, "nothing", resp2["echo"])
}
