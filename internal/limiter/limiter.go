package limiter

import "context"

// RateRequest represents the incoming JSON payload
type RateRequest struct {
	Key       string `json:"key"`
	Limit     int64  `json:"limit"`
	WindowMs  int64  `json:"window_ms"`
	Algorithm string `json:"algorithm"`
}

// RateResponse represents the outgoing JSON payload
type RateResponse struct {
	Allowed bool `json:"allowed"`
}

type RateLimiter interface {
	Allow(ctx context.Context, req *RateRequest) (*RateResponse, error)
}
