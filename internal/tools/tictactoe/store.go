package tictactoe

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var ErrGameNotFound = errors.New("game not found")

const gameTTL = time.Hour

type Store struct {
	client *redis.Client
}

func NewStore(client *redis.Client) *Store {
	return &Store{client: client}
}

func gameKey(id string) string {
	return fmt.Sprintf("game:%s", id)
}

func (s *Store) Save(ctx context.Context, g Game) error {
	data, err := json.Marshal(g)
	if err != nil {
		return err
	}
	return s.client.Set(ctx, gameKey(g.ID), data, gameTTL).Err()
}

func (s *Store) Get(ctx context.Context, id string) (Game, error) {
	data, err := s.client.Get(ctx, gameKey(id)).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return Game{}, ErrGameNotFound
		}
		return Game{}, err
	}
	var g Game
	if err := json.Unmarshal(data, &g); err != nil {
		return Game{}, err
	}
	return g, nil
}
