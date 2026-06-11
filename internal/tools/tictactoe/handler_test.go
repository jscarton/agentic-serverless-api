package tictactoe_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	"github.com/jscarton/agentic-serverless-api/internal/middleware"
	"github.com/jscarton/agentic-serverless-api/internal/tools/tictactoe"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestRouter() (*gin.Engine, *miniredis.Miniredis) {
	gin.SetMode(gin.TestMode)
	mr := miniredis.RunT(&testing.T{})
	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})

	store := tictactoe.NewStore(client)
	svc := tictactoe.NewService(store)

	r := gin.New()
	r.Use(middleware.RequestID())
	tictactoe.RegisterRoutes(r.Group("/tools/tictactoe"), svc)
	return r, mr
}

func TestCreateGame(t *testing.T) {
	r, _ := newTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/tools/tictactoe/games", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var body map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	data := body["data"].(map[string]any)
	assert.NotEmpty(t, data["id"])
	assert.Equal(t, "X", data["current_player"])
	assert.Equal(t, "in_progress", data["status"])
}

func TestMakeMove(t *testing.T) {
	r, _ := newTestRouter()

	// Create a game first
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/tools/tictactoe/games", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)
	var createBody map[string]any
	json.Unmarshal(w.Body.Bytes(), &createBody)
	gameID := createBody["data"].(map[string]any)["id"].(string)

	// Make a move
	moveBody, _ := json.Marshal(map[string]any{"position": 4, "player": "X"})
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest(http.MethodPost, "/tools/tictactoe/games/"+gameID+"/moves", bytes.NewReader(moveBody))
	req2.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusOK, w2.Code)
	var moveRespBody map[string]any
	json.Unmarshal(w2.Body.Bytes(), &moveRespBody)
	data := moveRespBody["data"].(map[string]any)
	board := data["board"].([]any)
	assert.Equal(t, "X", board[4])
	assert.Equal(t, "O", data["current_player"])
}

func TestMakeMove_InvalidPosition(t *testing.T) {
	r, _ := newTestRouter()

	// Create a game first
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/tools/tictactoe/games", nil)
	r.ServeHTTP(w, req)
	var createBody map[string]any
	json.Unmarshal(w.Body.Bytes(), &createBody)
	gameID := createBody["data"].(map[string]any)["id"].(string)

	// Invalid position
	moveBody, _ := json.Marshal(map[string]any{"position": 9, "player": "X"})
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest(http.MethodPost, "/tools/tictactoe/games/"+gameID+"/moves", bytes.NewReader(moveBody))
	req2.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusBadRequest, w2.Code)
}

func TestGetGame(t *testing.T) {
	r, _ := newTestRouter()

	// Create game
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/tools/tictactoe/games", nil)
	r.ServeHTTP(w, req)
	var createBody map[string]any
	json.Unmarshal(w.Body.Bytes(), &createBody)
	gameID := createBody["data"].(map[string]any)["id"].(string)

	// Get game
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest(http.MethodGet, "/tools/tictactoe/games/"+gameID, nil)
	r.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusOK, w2.Code)
	var body map[string]any
	json.Unmarshal(w2.Body.Bytes(), &body)
	data := body["data"].(map[string]any)
	assert.Equal(t, gameID, data["id"])
}

func TestGetGame_NotFound(t *testing.T) {
	r, _ := newTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/tools/tictactoe/games/nonexistent-id", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// Compile-time check that Store can be used via Service
var _ = tictactoe.NewStore
var _ = func(ctx context.Context) { _ = ctx }
