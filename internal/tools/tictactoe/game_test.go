package tictactoe_test

import (
	"testing"

	"github.com/jscarton/agentic-serverless-api/internal/tools/tictactoe"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewGame_InitialState(t *testing.T) {
	g := tictactoe.NewGame("game-1")
	assert.Equal(t, "game-1", g.ID)
	assert.Equal(t, "X", g.CurrentPlayer)
	assert.Equal(t, tictactoe.StatusInProgress, g.Status)
	for _, cell := range g.Board {
		assert.Equal(t, "", cell)
	}
}

func TestApplyMove_BasicMove(t *testing.T) {
	g := tictactoe.NewGame("game-1")
	updated, err := tictactoe.ApplyMove(g, 4, "X")
	require.NoError(t, err)
	assert.Equal(t, "X", updated.Board[4])
	assert.Equal(t, "O", updated.CurrentPlayer)
	assert.Equal(t, tictactoe.StatusInProgress, updated.Status)
}

func TestApplyMove_WrongPlayer(t *testing.T) {
	g := tictactoe.NewGame("game-1")
	_, err := tictactoe.ApplyMove(g, 0, "O")
	assert.ErrorIs(t, err, tictactoe.ErrWrongPlayer)
}

func TestApplyMove_OccupiedPosition(t *testing.T) {
	g := tictactoe.NewGame("game-1")
	g, _ = tictactoe.ApplyMove(g, 0, "X")
	_, err := tictactoe.ApplyMove(g, 0, "O")
	assert.ErrorIs(t, err, tictactoe.ErrPositionOccupied)
}

func TestApplyMove_InvalidPosition(t *testing.T) {
	g := tictactoe.NewGame("game-1")
	_, err := tictactoe.ApplyMove(g, 9, "X")
	assert.ErrorIs(t, err, tictactoe.ErrInvalidPosition)
	_, err = tictactoe.ApplyMove(g, -1, "X")
	assert.ErrorIs(t, err, tictactoe.ErrInvalidPosition)
}

func TestApplyMove_XWins_Row(t *testing.T) {
	g := tictactoe.NewGame("game-1")
	// X: 0,1,2 — O: 3,4
	moves := []struct{ pos int; player string }{
		{0, "X"}, {3, "O"}, {1, "X"}, {4, "O"}, {2, "X"},
	}
	var err error
	for _, m := range moves {
		g, err = tictactoe.ApplyMove(g, m.pos, m.player)
		require.NoError(t, err)
	}
	assert.Equal(t, tictactoe.StatusXWins, g.Status)
}

func TestApplyMove_OWins_Column(t *testing.T) {
	g := tictactoe.NewGame("game-1")
	// X: 0,1,3 — O: 2,5,8 (column 2)
	moves := []struct{ pos int; player string }{
		{0, "X"}, {2, "O"}, {1, "X"}, {5, "O"}, {3, "X"}, {8, "O"},
	}
	var err error
	for _, m := range moves {
		g, err = tictactoe.ApplyMove(g, m.pos, m.player)
		require.NoError(t, err)
	}
	assert.Equal(t, tictactoe.StatusOWins, g.Status)
}

func TestApplyMove_Draw(t *testing.T) {
	g := tictactoe.NewGame("game-1")
	// Draw board: X O X / O X X / O X O
	moves := []struct{ pos int; player string }{
		{0, "X"}, {1, "O"}, {2, "X"},
		{3, "O"}, {4, "X"}, {6, "O"},
		{5, "X"}, {8, "O"}, {7, "X"},
	}
	var err error
	for _, m := range moves {
		g, err = tictactoe.ApplyMove(g, m.pos, m.player)
		require.NoError(t, err)
	}
	assert.Equal(t, tictactoe.StatusDraw, g.Status)
}

func TestApplyMove_GameOver(t *testing.T) {
	g := tictactoe.NewGame("game-1")
	moves := []struct{ pos int; player string }{
		{0, "X"}, {3, "O"}, {1, "X"}, {4, "O"}, {2, "X"},
	}
	var err error
	for _, m := range moves {
		g, err = tictactoe.ApplyMove(g, m.pos, m.player)
		require.NoError(t, err)
	}
	_, err = tictactoe.ApplyMove(g, 5, "O")
	assert.ErrorIs(t, err, tictactoe.ErrGameOver)
}
