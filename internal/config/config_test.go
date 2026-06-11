package config_test

import (
	"os"
	"testing"

	"github.com/jscarton/agentic-serverless-api/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestLoad_Defaults(t *testing.T) {
	os.Clearenv()
	cfg := config.Load()
	assert.Equal(t, "", cfg.ClerkSecretKey)
	assert.Equal(t, "", cfg.DatabaseURL)
	assert.Equal(t, "", cfg.RedisURL)
	assert.False(t, cfg.SwaggerEnabled)
	assert.Equal(t, "8080", cfg.Port)
}

func TestLoad_FromEnv(t *testing.T) {
	os.Setenv("CLERK_SECRET_KEY", "sk_test_abc")
	os.Setenv("DATABASE_URL", "postgres://user:pass@host/db")
	os.Setenv("REDIS_URL", "redis://localhost:6379")
	os.Setenv("SWAGGER_ENABLED", "true")
	os.Setenv("PORT", "9090")
	t.Cleanup(os.Clearenv)

	cfg := config.Load()
	assert.Equal(t, "sk_test_abc", cfg.ClerkSecretKey)
	assert.Equal(t, "postgres://user:pass@host/db", cfg.DatabaseURL)
	assert.Equal(t, "redis://localhost:6379", cfg.RedisURL)
	assert.True(t, cfg.SwaggerEnabled)
	assert.Equal(t, "9090", cfg.Port)
}
