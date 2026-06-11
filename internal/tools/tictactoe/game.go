package tictactoe

import "errors"

var (
	ErrGameOver         = errors.New("game is already over")
	ErrWrongPlayer      = errors.New("not this player's turn")
	ErrInvalidPosition  = errors.New("position must be between 0 and 8")
	ErrPositionOccupied = errors.New("position is already occupied")
)

const (
	StatusInProgress = "in_progress"
	StatusXWins      = "x_wins"
	StatusOWins      = "o_wins"
	StatusDraw       = "draw"
)

type Board [9]string

type Game struct {
	ID            string `json:"id"`
	Board         Board  `json:"board"`
	CurrentPlayer string `json:"current_player"`
	Status        string `json:"status"`
}

func NewGame(id string) Game {
	return Game{
		ID:            id,
		Board:         Board{},
		CurrentPlayer: "X",
		Status:        StatusInProgress,
	}
}

func ApplyMove(g Game, position int, player string) (Game, error) {
	if g.Status != StatusInProgress {
		return g, ErrGameOver
	}
	if player != g.CurrentPlayer {
		return g, ErrWrongPlayer
	}
	if position < 0 || position > 8 {
		return g, ErrInvalidPosition
	}
	if g.Board[position] != "" {
		return g, ErrPositionOccupied
	}

	g.Board[position] = player
	g.Status = checkStatus(g.Board)
	if g.Status == StatusInProgress {
		if player == "X" {
			g.CurrentPlayer = "O"
		} else {
			g.CurrentPlayer = "X"
		}
	}
	return g, nil
}

var winLines = [8][3]int{
	{0, 1, 2}, {3, 4, 5}, {6, 7, 8},
	{0, 3, 6}, {1, 4, 7}, {2, 5, 8},
	{0, 4, 8}, {2, 4, 6},
}

func checkStatus(b Board) string {
	for _, line := range winLines {
		a, c := b[line[0]], b[line[2]]
		if a != "" && a == b[line[1]] && b[line[1]] == c {
			if a == "X" {
				return StatusXWins
			}
			return StatusOWins
		}
	}
	for _, cell := range b {
		if cell == "" {
			return StatusInProgress
		}
	}
	return StatusDraw
}
