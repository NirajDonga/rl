package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	RedisURL string
}

func Load() (Config, error) {
	_ = godotenv.Load()

	redisURL, err := requiredEnv("REDIS_URL")
	if err != nil {
		return Config{}, err
	}

	cfg := Config{
		RedisURL: redisURL,
	}

	return cfg, nil
}

func requiredEnv(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("%s is required", key)
	}

	return value, nil
}
