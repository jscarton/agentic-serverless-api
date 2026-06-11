// cmd/api/main.go
package main

// @title           Agentic Serverless API
// @version         1.0
// @description     A production-ready Go API starter kit designed to be consumed by AI agents.
// @host            localhost:8080
// @BasePath        /

import (
	"context"
	"log/slog"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	clerkgo "github.com/clerk/clerk-sdk-go/v2"
	"github.com/redis/go-redis/v9"

	"github.com/jscarton/agentic-serverless-api/internal/cache"
	"github.com/jscarton/agentic-serverless-api/internal/config"
	"github.com/jscarton/agentic-serverless-api/internal/db"
	"github.com/jscarton/agentic-serverless-api/internal/mcp"
	"github.com/jscarton/agentic-serverless-api/internal/server"
)

func main() {
	cfg := config.Load()

	if cfg.ClerkSecretKey != "" {
		clerkgo.SetKey(cfg.ClerkSecretKey)
	}

	ctx := context.Background()

	var dbPinger interface{ Ping(ctx context.Context) error }
	if cfg.DatabaseURL != "" {
		dbConn, err := db.New(ctx, cfg.DatabaseURL)
		if err != nil {
			slog.Error("failed to connect to database", "error", err)
			os.Exit(1)
		}
		defer dbConn.Close()
		dbPinger = dbConn
	}

	var cachePinger interface{ Ping(ctx context.Context) error }
	var redisClient *redis.Client
	if cfg.RedisURL != "" {
		cacheConn, err := cache.New(cfg.RedisURL)
		if err != nil {
			slog.Error("failed to connect to Redis", "error", err)
			os.Exit(1)
		}
		defer cacheConn.Close()
		cachePinger = cacheConn
		redisClient = cacheConn.Client()
	}

	registry := mcp.NewRegistry()
	router := server.New(cfg, dbPinger, cachePinger, redisClient, registry)

	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
		adapter := ginadapter.NewV2(router)
		lambda.Start(adapter.ProxyWithContext)
	} else {
		slog.Info("starting local server", "port", cfg.Port)
		if err := router.Run(":" + cfg.Port); err != nil {
			slog.Error("server failed", "error", err)
			os.Exit(1)
		}
	}
}
