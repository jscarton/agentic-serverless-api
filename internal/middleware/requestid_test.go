package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jscarton/agentic-serverless-api/internal/middleware"
	"github.com/stretchr/testify/assert"
)

func TestRequestID_GeneratesID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	r := gin.New()
	r.Use(middleware.RequestID())
	var capturedID string
	r.GET("/test", func(c *gin.Context) {
		capturedID = c.GetString("request_id")
		c.Status(http.StatusOK)
	})

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	assert.NotEmpty(t, capturedID)
	assert.Equal(t, capturedID, w.Header().Get("X-Request-ID"))
}

func TestRequestID_PassthroughExistingID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	r := gin.New()
	r.Use(middleware.RequestID())
	var capturedID string
	r.GET("/test", func(c *gin.Context) {
		capturedID = c.GetString("request_id")
		c.Status(http.StatusOK)
	})

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-Request-ID", "my-custom-id")
	r.ServeHTTP(w, req)

	assert.Equal(t, "my-custom-id", capturedID)
	assert.Equal(t, "my-custom-id", w.Header().Get("X-Request-ID"))
}
