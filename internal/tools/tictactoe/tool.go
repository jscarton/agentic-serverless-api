package tictactoe

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jscarton/agentic-serverless-api/internal/mcp"
)

func RegisterMCPTools(registry *mcp.Registry, store *Store) {
	registry.Register(mcp.Tool{
		Name:        "tictactoe_new_game",
		Description: "Start a new tic-tac-toe game. X always goes first. Returns game_id, board, current_player, and status.",
		InputSchema: map[string]any{
			"type":       "object",
			"properties": map[string]any{},
		},
		Handler: func(ctx context.Context, args map[string]any) (any, error) {
			game := NewGame(uuid.New().String())
			if err := store.Save(ctx, game); err != nil {
				return nil, err
			}
			return map[string]any{
				"game_id":        game.ID,
				"board":          game.Board,
				"current_player": game.CurrentPlayer,
				"status":         game.Status,
			}, nil
		},
	})

	registry.Register(mcp.Tool{
		Name:        "tictactoe_make_move",
		Description: "Make a move in a tic-tac-toe game. Position is 0-8 (row-major: 0=top-left, 4=center, 8=bottom-right). Player must be 'X' or 'O' and must match current_player.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"game_id":  map[string]any{"type": "string", "description": "The game ID returned by tictactoe_new_game"},
				"position": map[string]any{"type": "integer", "minimum": 0, "maximum": 8, "description": "Board position 0-8"},
				"player":   map[string]any{"type": "string", "enum": []string{"X", "O"}, "description": "The player making the move"},
			},
			"required": []string{"game_id", "position", "player"},
		},
		Handler: func(ctx context.Context, args map[string]any) (any, error) {
			gameID, _ := args["game_id"].(string)
			posFloat, _ := args["position"].(float64)
			player, _ := args["player"].(string)

			if gameID == "" {
				return nil, errors.New("game_id is required")
			}

			game, err := store.Get(ctx, gameID)
			if err != nil {
				return nil, err
			}

			game, err = ApplyMove(game, int(posFloat), player)
			if err != nil {
				return nil, err
			}

			if err := store.Save(ctx, game); err != nil {
				return nil, err
			}

			return map[string]any{
				"game_id":        game.ID,
				"board":          game.Board,
				"current_player": game.CurrentPlayer,
				"status":         game.Status,
			}, nil
		},
	})

	registry.Register(mcp.Tool{
		Name:        "tictactoe_get_state",
		Description: "Get the current state of a tic-tac-toe game.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"game_id": map[string]any{"type": "string", "description": "The game ID"},
			},
			"required": []string{"game_id"},
		},
		Handler: func(ctx context.Context, args map[string]any) (any, error) {
			gameID, _ := args["game_id"].(string)
			if gameID == "" {
				return nil, fmt.Errorf("game_id is required")
			}
			game, err := store.Get(ctx, gameID)
			if err != nil {
				return nil, err
			}
			return map[string]any{
				"game_id":        game.ID,
				"board":          game.Board,
				"current_player": game.CurrentPlayer,
				"status":         game.Status,
			}, nil
		},
	})
}
