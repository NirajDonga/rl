package main

import (
	"context"
	"log"
	"net/http"

	"github.com/NirajDonga/rl/internal/config"
	"github.com/NirajDonga/rl/internal/handler"
	"github.com/NirajDonga/rl/internal/limiter"
	"github.com/redis/go-redis/v9"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("invalid configuration: %v", err)
	}

	redisOpts, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		log.Fatalf("invalid REDIS_URL: %v", err)
	}

	rdb := redis.NewClient(redisOpts)
	defer rdb.Close()

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	tbLimiter := limiter.NewTokenBucket(rdb)
	algorithms := map[string]limiter.RateLimiter{
		"token_bucket": tbLimiter,
	}

	apiHandler := handler.NewRateLimiterHandler(algorithms)

	http.HandleFunc("/api/v1/allow", apiHandler.HandleAllow)

	port := ":8080"
	log.Printf("Starting HTTP Rate Limiter Service on port %s...", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
