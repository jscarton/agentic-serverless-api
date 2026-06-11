package tictactoe_test

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/jscarton/agentic-serverless-api/internal/tools/tictactoe"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestStore(t *testing.T) (*tictactoe.Store, *miniredis.Miniredis) {
	mr := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	return tictactoe.NewStore(client), mr
}

func TestStore_SaveAndGet(t *testing.T) {
	store, _ := newTestStore(t)
	ctx := context.Background()

	game := tictactoe.NewGame("test-id-1")
	require.NoError(t, store.Save(ctx, game))

	got, err := store.Get(ctx, "test-id-1")
	require.NoError(t, err)
	assert.Equal(t, "test-id-1", got.ID)
	assert.Equal(t, "X", got.CurrentPlayer)
	assert.Equal(t, tictactoe.StatusInProgress, got.Status)
}

func TestStore_Get_NotFound(t *testing.T) {
	store, _ := newTestStore(t)
	ctx := context.Background()

	_, err := store.Get(ctx, "nonexistent")
	assert.ErrorIs(t, err, tictactoe.ErrGameNotFound)
}

func TestStore_Save_UpdatesExisting(t *testing.T) {
	store, _ := newTestStore(t)
	ctx := context.Background()

	game := tictactoe.NewGame("test-id-2")
	require.NoError(t, store.Save(ctx, game))

	game, _ = tictactoe.ApplyMove(game, 4, "X")
	require.NoError(t, store.Save(ctx, game))

	got, err := store.Get(ctx, "test-id-2")
	require.NoError(t, err)
	assert.Equal(t, "X", got.Board[4])
	assert.Equal(t, "O", got.CurrentPlayer)
}
