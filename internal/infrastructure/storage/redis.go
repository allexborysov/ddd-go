package redis

import (
	"context"
	"log"

	redisv9 "github.com/redis/go-redis/v9"
)

func MustConnect(options *redisv9.Options) *redisv9.Client {
	client := redisv9.NewClient(options)
	ctx := context.Background()

	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	return client
}
