package limiter

import (
	"context"
	"time"

	"github.com/NirajDonga/rl/internal/store/lua"
	"github.com/redis/go-redis/v9"
)

type TokenBucket struct {
	client *redis.Client
	script *redis.Script
}

func NewTokenBucket(client *redis.Client) *TokenBucket {
	return &TokenBucket{
		client: client,
		script: redis.NewScript(lua.TokenBucket),
	}
}

func (tb *TokenBucket) Allow(ctx context.Context, req *RateRequest) (*RateResponse, error) {
	tokensKey := req.Key + ":tokens"
	tsKey := req.Key + ":ts"

	nowMs := time.Now().UnixMilli()
	keys := []string{tokensKey, tsKey}

	// Read Limit and WindowMs from the standard Go struct
	args := []interface{}{req.Limit, req.WindowMs, nowMs}

	result, err := tb.script.Run(ctx, tb.client, keys, args...).Result()
	if err != nil {
		return nil, err
	}

	allowed := result.(int64) == 1

	return &RateResponse{
		Allowed: allowed,
	}, nil
}
