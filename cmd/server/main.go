package main

import (
	"context"
	"log"
	"net"

	pb "github.com/NirajDonga/rl/api/ratelimit/v1"
	"github.com/NirajDonga/rl/internal/config"
	"github.com/NirajDonga/rl/internal/limiter"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type rateLimiterServer struct {
	// Embedding this is required by the generated gRPC code for forward compatibility
	pb.UnimplementedRateLimiterServiceServer

	// We hold instances of our algorithm implementations here
	tokenBucket *limiter.TokenBucket
	// We will add fixedWindow and slidingWindow here in the future
}

func (s *rateLimiterServer) IsAllowed(ctx context.Context, req *pb.IsAllowedRequest) (*pb.IsAllowedResponse, error) {
	if req.Key == "" {
		return nil, status.Error(codes.InvalidArgument, "rate limit key is required")
	}

	// Route to the correct algorithm based on the gRPC request
	switch req.Algorithm {
	case pb.Algorithm_ALGORITHM_TOKEN_BUCKET:
		return s.tokenBucket.Allow(ctx, req)

	case pb.Algorithm_ALGORITHM_FIXED_WINDOW, pb.Algorithm_ALGORITHM_SLIDING_WINDOW:
		return nil, status.Errorf(codes.Unimplemented, "algorithm %v not yet implemented", req.Algorithm)

	default:
		return nil, status.Errorf(codes.InvalidArgument, "unknown or unspecified algorithm: %v", req.Algorithm)
	}
}

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

	// 1. Initialize our limiters
	tbLimiter := limiter.NewTokenBucket(rdb)

	// 2. Create our gRPC server implementation
	srv := &rateLimiterServer{
		tokenBucket: tbLimiter,
	}

	// 3. Open a TCP listener on port 50051
	port := ":50051"
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// 4. Create the gRPC framework server and register our implementation
	grpcServer := grpc.NewServer()
	pb.RegisterRateLimiterServiceServer(grpcServer, srv)

	log.Printf("Starting Rate Limiter gRPC Service on port %s...", port)

	// 5. Start serving! This will block and listen forever.
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
