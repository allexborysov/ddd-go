package bookingsync

import (
	"context"
	"time"

	redisv9 "github.com/redis/go-redis/v9"
)

type BookingMutex struct {
	redisClient *redisv9.Client
	ttl         time.Duration
	prefix      string
}

func New(client *redisv9.Client) *BookingMutex {
	return &BookingMutex{
		redisClient: client,
		ttl:         time.Minute,
		prefix:      "lock:seat:",
	}
}

func (m *BookingMutex) Lock(ctx context.Context, key string) bool {
	rKey := m.prefix + key
	locked, err := m.redisClient.SetNX(ctx, rKey, "", m.ttl).Result()
	if err != nil {
		return false
	}
	return locked
}

func (m *BookingMutex) Unlock(ctx context.Context, key string) {
	rKey := m.prefix + key
	m.redisClient.Del(ctx, rKey)
}
