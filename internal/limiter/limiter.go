package limiter

import "context"

type RateRequest struct {
	Key      string
	Limit    int64
	Windowms int64
}

type RateLimiter interface {
	Allow(ctx context.Context, req RateRequest) (bool, error)
}
