package health_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jscarton/agentic-serverless-api/internal/health"
	"github.com/jscarton/agentic-serverless-api/internal/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHealth_AllUnconfigured(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	r := gin.New()
	r.Use(middleware.RequestID())
	r.GET("/health", health.NewHandler(nil, nil).Handle)

	req, _ := http.NewRequest(http.MethodGet, "/health", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var body map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	data := body["data"].(map[string]any)
	assert.Equal(t, "ok", data["status"])
	assert.Equal(t, "unconfigured", data["db"])
	assert.Equal(t, "unconfigured", data["cache"])
}
