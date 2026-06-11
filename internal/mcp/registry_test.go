// internal/mcp/registry_test.go
package mcp_test

import (
	"context"
	"testing"

	"github.com/jscarton/agentic-serverless-api/internal/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegistry_RegisterAndList(t *testing.T) {
	r := mcp.NewRegistry()
	r.Register(mcp.Tool{
		Name:        "echo",
		Description: "echoes input",
		InputSchema: map[string]any{"type": "object"},
		Handler: func(ctx context.Context, args map[string]any) (any, error) {
			return args["message"], nil
		},
	})

	tools := r.List()
	require.Len(t, tools, 1)
	assert.Equal(t, "echo", tools[0].Name)
}

func TestRegistry_Call_Success(t *testing.T) {
	r := mcp.NewRegistry()
	r.Register(mcp.Tool{
		Name: "echo",
		Handler: func(ctx context.Context, args map[string]any) (any, error) {
			return args["message"], nil
		},
	})

	result, err := r.Call(context.Background(), "echo", map[string]any{"message": "hello"})
	require.NoError(t, err)
	assert.Equal(t, "hello", result)
}

func TestRegistry_Call_ToolNotFound(t *testing.T) {
	r := mcp.NewRegistry()
	_, err := r.Call(context.Background(), "nonexistent", nil)
	assert.ErrorIs(t, err, mcp.ErrToolNotFound)
}
