// internal/mcp/server_test.go
package mcp_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jscarton/agentic-serverless-api/internal/mcp"
	"github.com/jscarton/agentic-serverless-api/internal/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestRouter(registry *mcp.Registry) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.RequestID())
	r.POST("/mcp", mcp.NewServer(registry).Handle)
	return r
}

func TestMCP_ToolsList(t *testing.T) {
	registry := mcp.NewRegistry()
	registry.Register(mcp.Tool{
		Name:        "echo",
		Description: "echoes input",
		InputSchema: map[string]any{"type": "object"},
	})

	body, _ := json.Marshal(map[string]any{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "tools/list",
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/mcp", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	newTestRouter(registry).ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "2.0", resp["jsonrpc"])
	result := resp["result"].(map[string]any)
	tools := result["tools"].([]any)
	require.Len(t, tools, 1)
	tool := tools[0].(map[string]any)
	assert.Equal(t, "echo", tool["name"])
}

func TestMCP_ToolsCall_Success(t *testing.T) {
	registry := mcp.NewRegistry()
	registry.Register(mcp.Tool{
		Name: "echo",
		Handler: func(ctx context.Context, args map[string]any) (any, error) {
			return args["message"], nil
		},
	})

	body, _ := json.Marshal(map[string]any{
		"jsonrpc": "2.0",
		"id":      2,
		"method":  "tools/call",
		"params": map[string]any{
			"name":      "echo",
			"arguments": map[string]any{"message": "hello"},
		},
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/mcp", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	newTestRouter(registry).ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	result := resp["result"].(map[string]any)
	content := result["content"].([]any)
	require.Len(t, content, 1)
	item := content[0].(map[string]any)
	assert.Equal(t, "text", item["type"])
	assert.Contains(t, item["text"].(string), "hello")
	assert.Equal(t, false, result["isError"])
}

func TestMCP_ToolsCall_ToolNotFound(t *testing.T) {
	registry := mcp.NewRegistry()
	body, _ := json.Marshal(map[string]any{
		"jsonrpc": "2.0",
		"id":      3,
		"method":  "tools/call",
		"params": map[string]any{
			"name":      "nonexistent",
			"arguments": map[string]any{},
		},
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/mcp", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	newTestRouter(registry).ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	errObj := resp["error"].(map[string]any)
	assert.EqualValues(t, -32601, errObj["code"])
}

func TestMCP_UnknownMethod(t *testing.T) {
	registry := mcp.NewRegistry()
	body, _ := json.Marshal(map[string]any{
		"jsonrpc": "2.0",
		"id":      4,
		"method":  "unknown/method",
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/mcp", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	newTestRouter(registry).ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	errObj := resp["error"].(map[string]any)
	assert.EqualValues(t, -32601, errObj["code"])
}
