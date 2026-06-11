package tictactoe

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jscarton/agentic-serverless-api/internal/response"
)

// Service holds business logic for the tictactoe tool.
type Service struct {
	store *Store
}

// NewService creates a new Service backed by the given Store.
func NewService(store *Store) *Service {
	return &Service{store: store}
}

// RegisterRoutes mounts the tictactoe endpoints on the given RouterGroup.
func RegisterRoutes(rg *gin.RouterGroup, svc *Service) {
	rg.POST("/games", svc.CreateGame)
	rg.POST("/games/:id/moves", svc.MakeMove)
	rg.GET("/games/:id", svc.GetGame)
}

type moveRequest struct {
	Position *int   `json:"position" binding:"required,gte=0,lte=8"`
	Player   string `json:"player"   binding:"required,oneof=X O"`
}

// CreateGame godoc
// @Summary Create a new tic-tac-toe game
// @Description Creates a new game session. X always goes first.
// @Tags tictactoe
// @Produce json
// @Success 201 {object} map[string]any
// @Router /tools/tictactoe/games [post]
func (s *Service) CreateGame(c *gin.Context) {
	game := NewGame(uuid.New().String())
	if err := s.store.Save(c.Request.Context(), game); err != nil {
		response.Fail(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	response.Created(c, gin.H{
		"id":             game.ID,
		"board":          game.Board,
		"current_player": game.CurrentPlayer,
		"status":         game.Status,
	})
}

// MakeMove godoc
// @Summary Make a move in a tic-tac-toe game
// @Description Submits a player's move. Returns updated board state.
// @Tags tictactoe
// @Accept json
// @Produce json
// @Param id path string true "Game ID"
// @Param request body moveRequest true "Move details"
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Failure 404 {object} map[string]any
// @Router /tools/tictactoe/games/{id}/moves [post]
func (s *Service) MakeMove(c *gin.Context) {
	id := c.Param("id")
	var req moveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	game, err := s.store.Get(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrGameNotFound) {
			response.Fail(c, http.StatusNotFound, "GAME_NOT_FOUND", "game not found")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	game, err = ApplyMove(game, *req.Position, req.Player)
	if err != nil {
		code, msg := moveErrToResponse(err)
		response.Fail(c, http.StatusBadRequest, code, msg)
		return
	}

	if err := s.store.Save(c.Request.Context(), game); err != nil {
		response.Fail(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	response.OK(c, gin.H{
		"id":             game.ID,
		"board":          game.Board,
		"current_player": game.CurrentPlayer,
		"status":         game.Status,
	})
}

// GetGame godoc
// @Summary Get current game state
// @Description Returns the current state of a tic-tac-toe game.
// @Tags tictactoe
// @Produce json
// @Param id path string true "Game ID"
// @Success 200 {object} map[string]any
// @Failure 404 {object} map[string]any
// @Router /tools/tictactoe/games/{id} [get]
func (s *Service) GetGame(c *gin.Context) {
	id := c.Param("id")
	game, err := s.store.Get(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrGameNotFound) {
			response.Fail(c, http.StatusNotFound, "GAME_NOT_FOUND", "game not found")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	response.OK(c, gin.H{
		"id":             game.ID,
		"board":          game.Board,
		"current_player": game.CurrentPlayer,
		"status":         game.Status,
	})
}

func moveErrToResponse(err error) (code, message string) {
	switch {
	case errors.Is(err, ErrGameOver):
		return "GAME_OVER", "game is already over"
	case errors.Is(err, ErrWrongPlayer):
		return "WRONG_PLAYER", "not this player's turn"
	case errors.Is(err, ErrPositionOccupied):
		return "POSITION_OCCUPIED", "position is already occupied"
	default:
		return "INVALID_MOVE", err.Error()
	}
}
