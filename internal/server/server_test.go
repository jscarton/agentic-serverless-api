package server_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jscarton/agentic-serverless-api/internal/config"
	"github.com/jscarton/agentic-serverless-api/internal/mcp"
	"github.com/jscarton/agentic-serverless-api/internal/server"
	"github.com/stretchr/testify/assert"
)

func TestServer_HealthRoute(t *testing.T) {
	cfg := &config.Config{}
	registry := mcp.NewRegistry()
	r := server.New(cfg, nil, nil, nil, registry)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/health", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestServer_MCPRoute(t *testing.T) {
	cfg := &config.Config{}
	registry := mcp.NewRegistry()
	r := server.New(cfg, nil, nil, nil, registry)

	w := httptest.NewRecorder()
	body := []byte(`{"jsonrpc":"2.0","id":1,"method":"tools/list"}`)
	req, _ := http.NewRequest(http.MethodPost, "/mcp", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
