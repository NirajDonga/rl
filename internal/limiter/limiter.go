package limiter

import (
	"context"

	pb "github.com/NirajDonga/rl/api/ratelimit/v1"
)

type RateLimiter interface {
	Allow(ctx context.Context, req *pb.IsAllowedRequest) (*pb.IsAllowedResponse, error)
}
