package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	client *redis.Client
}

func New(url string) (*Cache, error) {
	opts, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(opts)
	return &Cache{client: client}, nil
}

func (c *Cache) Client() *redis.Client {
	return c.client
}

func (c *Cache) Ping(ctx context.Context) error {
	return c.client.Ping(ctx).Err()
}

func (c *Cache) Close() error {
	return c.client.Close()
}
