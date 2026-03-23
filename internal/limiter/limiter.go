package limiter

import (
	"context"

	pb "github.com/NirajDonga/rl/api/ratelimit/v1" // Import the generated code
)

type RateLimiter interface {
	// Notice we now return *pb.IsAllowedResponse
	Allow(ctx context.Context, req *pb.IsAllowedRequest) (*pb.IsAllowedResponse, error)
}