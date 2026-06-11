// internal/server/server.go
package server

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/jscarton/agentic-serverless-api/internal/config"
	"github.com/jscarton/agentic-serverless-api/internal/mcp"
	"github.com/jscarton/agentic-serverless-api/internal/middleware"
	"github.com/redis/go-redis/v9"
)

type pinger interface {
	Ping(ctx context.Context) error
}

func New(cfg *config.Config, db pinger, cache pinger, redisClient *redis.Client, registry *mcp.Registry) *gin.Engine {
	r := gin.New()
	r.Use(middleware.RequestID())
	r.Use(middleware.Logger())
	registerRoutes(r, cfg, db, cache, redisClient, registry)
	return r
}
