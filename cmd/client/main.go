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
	// Note: We use insecure credentials here because we don't have TLS/SSL set up locally.
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
		Key:       "{clientA:user1}", // A unique identifier for the user
		Limit:     5,                 // Allow 5 requests...
		WindowMs:  10000,             // ...every 10 seconds
		Algorithm: pb.Algorithm_ALGORITHM_TOKEN_BUCKET,
	}

	log.Println("Starting to send requests to the Rate Limiter Service over gRPC...")

	// 4. Simulate 8 rapid incoming HTTP requests
	for i := 1; i <= 8; i++ {
		// This makes a network call to your running server!
		res, err := client.IsAllowed(ctx, req)
		if err != nil {
			log.Fatalf("could not call rate limiter: %v", err)
		}

		if res.Allowed {
			log.Printf("Request %d: ALLOWED ✅", i)
		} else {
			log.Printf("Request %d: BLOCKED ⛔", i)
		}

		// Sleep for half a second between requests
		time.Sleep(500 * time.Millisecond)
	}
}
