// internal/server/routes.go
package server

import (
	"github.com/gin-gonic/gin"
	"github.com/jscarton/agentic-serverless-api/internal/config"
	"github.com/jscarton/agentic-serverless-api/internal/health"
	"github.com/jscarton/agentic-serverless-api/internal/mcp"
	"github.com/jscarton/agentic-serverless-api/internal/tools/tictactoe"
	"github.com/redis/go-redis/v9"
)

func registerRoutes(r *gin.Engine, cfg *config.Config, db pinger, cache pinger, redisClient *redis.Client, registry *mcp.Registry) {
	r.GET("/health", health.NewHandler(db, cache).Handle)
	r.POST("/mcp", mcp.NewServer(registry).Handle)

	// TODO(Task 17): Add Swagger routes here after `make swag` generates docs/swagger.
	// if cfg.SwaggerEnabled {
	// 	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// }

	if redisClient != nil {
		store := tictactoe.NewStore(redisClient)
		svc := tictactoe.NewService(store)
		tictactoe.RegisterMCPTools(registry, store)
		tictactoe.RegisterRoutes(r.Group("/tools/tictactoe"), svc)
	}
}
