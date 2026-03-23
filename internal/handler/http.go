package handler

import (
	"encoding/json"
	"net/http"

	"github.com/NirajDonga/rl/internal/limiter"
)

type RateLimiterHandler struct {
	algorithms map[string]limiter.RateLimiter
}

func NewRateLimiterHandler(algorithms map[string]limiter.RateLimiter) *RateLimiterHandler {
	return &RateLimiterHandler{
		algorithms: algorithms,
	}
}

func (h *RateLimiterHandler) HandleAllow(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req limiter.RateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if req.Key == "" {
		http.Error(w, "Rate limit key is required", http.StatusBadRequest)
		return
	}

	var res *limiter.RateResponse
	var err error

	algorithm := req.Algorithm
	if algorithm == "" {
		algorithm = "token_bucket"
	}

	selectedLimiter, ok := h.algorithms[algorithm]
	if !ok {
		http.Error(w, "Unknown algorithm", http.StatusBadRequest)
		return
	}

	res, err = selectedLimiter.Allow(r.Context(), &req)

	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
