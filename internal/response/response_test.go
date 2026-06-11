package response_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jscarton/agentic-serverless-api/internal/response"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setup() (*gin.Engine, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	r := gin.New()
	return r, w
}

func TestOK_EnvelopeShape(t *testing.T) {
	r, w := setup()
	r.GET("/test", func(c *gin.Context) {
		c.Set("request_id", "test-id-123")
		response.OK(c, gin.H{"value": 42})
	})
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var body map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.NotNil(t, body["data"])
	assert.Nil(t, body["error"])
	meta := body["meta"].(map[string]any)
	assert.Equal(t, "test-id-123", meta["request_id"])
	assert.NotEmpty(t, meta["timestamp"])
}

func TestFail_EnvelopeShape(t *testing.T) {
	r, w := setup()
	r.GET("/test", func(c *gin.Context) {
		c.Set("request_id", "test-id-456")
		response.Fail(c, http.StatusBadRequest, "INVALID_INPUT", "bad value")
	})
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var body map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Nil(t, body["data"])
	errObj := body["error"].(map[string]any)
	assert.Equal(t, "INVALID_INPUT", errObj["code"])
	assert.Equal(t, "bad value", errObj["message"])
}
