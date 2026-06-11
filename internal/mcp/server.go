// internal/mcp/server.go
package mcp

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type jsonRPCRequest struct {
	JSONRPC string         `json:"jsonrpc"`
	ID      any            `json:"id"`
	Method  string         `json:"method"`
	Params  map[string]any `json:"params"`
}

type jsonRPCResponse struct {
	JSONRPC string        `json:"jsonrpc"`
	ID      any           `json:"id"`
	Result  any           `json:"result,omitempty"`
	Error   *jsonRPCError `json:"error,omitempty"`
}

type jsonRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Server struct {
	registry *Registry
}

func NewServer(registry *Registry) *Server {
	return &Server{registry: registry}
}

// Handle godoc
// @Summary MCP endpoint (JSON-RPC 2.0)
// @Description Handles MCP Streamable HTTP transport. Supports tools/list and tools/call.
// @Tags mcp
// @Accept json
// @Produce json
// @Param request body object true "JSON-RPC 2.0 request"
// @Success 200 {object} map[string]any
// @Router /mcp [post]
func (s *Server) Handle(c *gin.Context) {
	var req jsonRPCRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, jsonRPCResponse{
			JSONRPC: "2.0",
			ID:      nil,
			Error:   &jsonRPCError{Code: -32700, Message: "parse error"},
		})
		return
	}

	switch req.Method {
	case "initialize":
		c.JSON(http.StatusOK, jsonRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: map[string]any{
				"protocolVersion": "2024-11-05",
				"capabilities":    map[string]any{"tools": map[string]any{}},
				"serverInfo":      map[string]any{"name": "tic-tac-toe", "version": "1.0.0"},
			},
		})

	case "tools/list":
		tools := s.registry.List()
		toolDefs := make([]map[string]any, 0, len(tools))
		for _, t := range tools {
			toolDefs = append(toolDefs, map[string]any{
				"name":        t.Name,
				"description": t.Description,
				"inputSchema": t.InputSchema,
			})
		}
		c.JSON(http.StatusOK, jsonRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  map[string]any{"tools": toolDefs},
		})

	case "tools/call":
		name, _ := req.Params["name"].(string)
		args, _ := req.Params["arguments"].(map[string]any)
		result, err := s.registry.Call(c.Request.Context(), name, args)
		if err != nil {
			if errors.Is(err, ErrToolNotFound) {
				c.JSON(http.StatusOK, jsonRPCResponse{
					JSONRPC: "2.0",
					ID:      req.ID,
					Error:   &jsonRPCError{Code: -32601, Message: "tool not found: " + name},
				})
				return
			}
			// Tool executed but returned an error — return as isError=true content
			c.JSON(http.StatusOK, jsonRPCResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Result: map[string]any{
					"content": []map[string]any{
						{"type": "text", "text": err.Error()},
					},
					"isError": true,
				},
			})
			return
		}
		text, _ := json.Marshal(result)
		c.JSON(http.StatusOK, jsonRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: map[string]any{
				"content": []map[string]any{
					{"type": "text", "text": string(text)},
				},
				"isError": false,
			},
		})

	default:
		c.JSON(http.StatusOK, jsonRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error:   &jsonRPCError{Code: -32601, Message: "method not found: " + req.Method},
		})
	}
}
