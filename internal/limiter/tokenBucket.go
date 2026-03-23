package limiter

import (
	"context"
	"time"

	pb "github.com/NirajDonga/rl/api/ratelimit/v1"
	"github.com/NirajDonga/rl/internal/store/lua" // Import the new lua package
	"github.com/redis/go-redis/v9"
)

type TokenBucket struct {
	client *redis.Client
	script *redis.Script
}

func NewTokenBucket(client *redis.Client) *TokenBucket {
	return &TokenBucket{
		client: client,
		// Use the exported string from the lua package
		script: redis.NewScript(lua.TokenBucket),
	}
}

func (tb *TokenBucket) Allow(ctx context.Context, req *pb.IsAllowedRequest) (*pb.IsAllowedResponse, error) {
	tokensKey := req.Key + ":tokens"
	tsKey := req.Key + ":ts"

	nowMs := time.Now().UnixMilli()
	keys := []string{tokensKey, tsKey}

	args := []interface{}{req.Limit, req.WindowMs, nowMs}

	result, err := tb.script.Run(ctx, tb.client, keys, args...).Result()
	if err != nil {
		return nil, err
	}

	allowed := result.(int64) == 1

	return &pb.IsAllowedResponse{
		Allowed: allowed,
	}, nil
}
