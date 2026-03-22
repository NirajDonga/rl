package limiter

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var tokenBucketLua string

// TokenBucket implements the RateLimiter interface using the token bucket algorithm.
type TokenBucket struct {
	client *redis.Client
	script *redis.Script
}

// NewTokenBucket initializes the token bucket limiter.
func NewTokenBucket(client *redis.Client) *TokenBucket {
	return &TokenBucket{
		client: client,
		script: redis.NewScript(tokenBucketLua),
	}
}

// Allow executes the specific token bucket logic in Redis.
func (tb *TokenBucket) Allow(ctx context.Context, req RateRequest) (bool, error) {
	// 1. Construct the keys specific to the token bucket algorithm
	tokensKey := req.Key + ":tokens"
	tsKey := req.Key + ":ts"

	nowMs := time.Now().UnixMilli()

	keys := []string{tokensKey, tsKey}

	// 2. These arguments map exactly to ARGV[1], ARGV[2], ARGV[3] in token_bucket.lua
	args := []interface{}{req.Limit, req.Windowms, nowMs}

	result, err := tb.script.Run(ctx, tb.client, keys, args...).Result()
	if err != nil {
		return false, err
	}

	allowed := result.(int64) == 1
	return allowed, nil
}
