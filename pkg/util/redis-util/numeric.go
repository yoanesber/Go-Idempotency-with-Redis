package redis_util

import (
	"context"
	"fmt"

	"github.com/yoanesber/go-idempotency-with-redis/config/cache"
)

// Increment increases a key's value by 1 (or given amount)
// If the key does not exist, it will be created with the specified value.
func Increment(key string, by int64) (int64, error) {
	// Get the Redis client from the context
	client := cache.GetRedisClient()
	if client == nil {
		return 0, fmt.Errorf("redis client is nil")
	}

	return client.IncrBy(context.Background(), key, by).Result()
}

// Decrement decreases a key's value by 1 (or given amount)
// If the key does not exist, it will be created with the specified value.
func Decrement(key string, by int64) (int64, error) {
	// Get the Redis client from the context
	client := cache.GetRedisClient()
	if client == nil {
		return 0, fmt.Errorf("redis client is nil")
	}

	return client.DecrBy(context.Background(), key, by).Result()
}
