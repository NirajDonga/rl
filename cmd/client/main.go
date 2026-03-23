package main

import (
	"context"
	"log"
	"time"

	pb "github.com/NirajDonga/rl/api/ratelimit/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// 1. Connect to the Rate Limiter gRPC Server we started in Phase 3
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// 2. Create the gRPC client using the generated code
	client := pb.NewRateLimiterServiceClient(conn)

	// 3. Define the rate limiting rules for this specific endpoint/user
	ctx := context.Background()
	req := &pb.IsAllowedRequest{
		Key:       "{clientA:user1}",
		Limit:     5,
		WindowMs:  10000,
		Algorithm: pb.Algorithm_ALGORITHM_TOKEN_BUCKET,
	}

	log.Println("Starting to send requests to the Rate Limiter Service over gRPC...")

	// 4. Simulate 8 rapid incoming HTTP requests
	for i := 1; i <= 8; i++ {
		res, err := client.IsAllowed(ctx, req)
		if err != nil {
			log.Fatalf("could not call rate limiter: %v", err)
		}

		if res.Allowed {
			log.Printf("Request %d: ALLOWED", i)
		} else {
			log.Printf("Request %d: BLOCKED", i)
		}

		time.Sleep(500 * time.Millisecond)
	}
}
