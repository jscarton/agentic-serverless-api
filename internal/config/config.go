package config

import "os"

type Config struct {
	ClerkSecretKey string
	DatabaseURL    string
	RedisURL       string
	SwaggerEnabled bool
	Port           string
}

func Load() *Config {
	return &Config{
		ClerkSecretKey: os.Getenv("CLERK_SECRET_KEY"),
		DatabaseURL:    os.Getenv("DATABASE_URL"),
		RedisURL:       os.Getenv("REDIS_URL"),
		SwaggerEnabled: os.Getenv("SWAGGER_ENABLED") == "true",
		Port:           getEnv("PORT", "8080"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
