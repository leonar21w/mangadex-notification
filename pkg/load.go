package pkg

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	RedisURL     string
	RedisToken   string
	Username     string
	Password     string
	ClientID     string
	ClientSecret string
}

func Load() (*Config, error) {

	_ = godotenv.Load()

	cfg := &Config{
		RedisURL:     "",
		RedisToken:   "",
		Username:     os.Getenv("MGDEX_USERNAME"),
		Password:     os.Getenv("MGDEX_PASSWORD"),
		ClientID:     os.Getenv("MGDEX_CLIENT"),
		ClientSecret: os.Getenv("MGDEX_SECRET"),
	}

	if os.Getenv("CURRENT_ENV") != "dev" {
		cfg.RedisURL = os.Getenv("REDIS_URL")
		cfg.RedisToken = os.Getenv("REDIS_TOKEN")
	} else {
		cfg.RedisURL = os.Getenv("REDIS_URL_TEST")
		cfg.RedisToken = os.Getenv("REDIS_TOKEN_TEST")
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
