package main

import (
	"context"
	"log"
	"time"

	"github.com/NirajDonga/rl/internal/config"
	"github.com/NirajDonga/rl/internal/limiter"
	"github.com/redis/go-redis/v9"
)

const (
	rateKey      = "{clientA:user1}"
	rateLimit    = int64(5)
	rateWindowMs = int64(10000)
	simRequests  = 8
	simSleepMs   = 500
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

	ctx := context.Background()

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v. Is Redis running?", err)
	}

	tbLimiter := limiter.NewTokenBucket(rdb)

	reqInfo := limiter.RateRequest{
		Key:      rateKey,
		Limit:    rateLimit,
		Windowms: rateWindowMs,
	}

	for i := 1; i <= simRequests; i++ {
		allowed, err := tbLimiter.Allow(ctx, reqInfo)
		if err != nil {
			log.Fatalf("Error checking rate limit: %v", err)
		}

		if allowed {
			log.Printf("request=%d status=allowed", i)
		} else {
			log.Printf("request=%d status=blocked", i)
		}

		time.Sleep(time.Duration(simSleepMs) * time.Millisecond)
	}
}
