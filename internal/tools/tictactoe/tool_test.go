package tictactoe_test

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/jscarton/agentic-serverless-api/internal/mcp"
	"github.com/jscarton/agentic-serverless-api/internal/tools/tictactoe"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestMCPSetup(t *testing.T) (*mcp.Registry, *miniredis.Miniredis) {
	mr := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	store := tictactoe.NewStore(client)

	registry := mcp.NewRegistry()
	tictactoe.RegisterMCPTools(registry, store)
	return registry, mr
}

func TestMCPTool_NewGame(t *testing.T) {
	registry, _ := newTestMCPSetup(t)
	ctx := context.Background()

	result, err := registry.Call(ctx, "tictactoe_new_game", nil)
	require.NoError(t, err)

	data := result.(map[string]any)
	assert.NotEmpty(t, data["game_id"])
	assert.Equal(t, "X", data["current_player"])
	assert.Equal(t, "in_progress", data["status"])
}

func TestMCPTool_MakeMove(t *testing.T) {
	registry, _ := newTestMCPSetup(t)
	ctx := context.Background()

	// Create game
	result, err := registry.Call(ctx, "tictactoe_new_game", nil)
	require.NoError(t, err)
	gameID := result.(map[string]any)["game_id"].(string)

	// Make move — position comes as float64 from JSON decode
	result2, err := registry.Call(ctx, "tictactoe_make_move", map[string]any{
		"game_id":  gameID,
		"position": float64(4),
		"player":   "X",
	})
	require.NoError(t, err)

	data := result2.(map[string]any)
	boardAny := data["board"].(tictactoe.Board)
	assert.Equal(t, "X", boardAny[4])
	assert.Equal(t, "O", data["current_player"])
}

func TestMCPTool_GetState(t *testing.T) {
	registry, _ := newTestMCPSetup(t)
	ctx := context.Background()

	// Create game
	result, err := registry.Call(ctx, "tictactoe_new_game", nil)
	require.NoError(t, err)
	gameID := result.(map[string]any)["game_id"].(string)

	// Get state
	result2, err := registry.Call(ctx, "tictactoe_get_state", map[string]any{
		"game_id": gameID,
	})
	require.NoError(t, err)

	data := result2.(map[string]any)
	assert.Equal(t, gameID, data["game_id"])
	assert.Equal(t, "in_progress", data["status"])
}
