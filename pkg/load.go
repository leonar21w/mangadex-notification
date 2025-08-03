package pkg

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	RedisURL   string
	RedisToken string
}

func Load() (*Config, error) {

	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("loading .env: %w", err)
	}

	cfg := &Config{
		RedisURL:   os.Getenv("REDIS_URL"),
		RedisToken: os.Getenv("REDIS_TOKEN"),
	}

	// validate required fields
	if cfg.RedisURL == "" {
		return nil, fmt.Errorf("environment variable REDIS_URL is required")
	}
	if cfg.RedisToken == "" {
		return nil, fmt.Errorf("environment variable REDIS_TOKEN is required")
	}
	return cfg, nil
}
