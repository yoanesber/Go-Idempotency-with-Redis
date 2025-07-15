package redis_util

import (
	"context"
	"fmt"
	"time"

	"github.com/yoanesber/go-idempotency-api/config/cache"
)

// Set sets a string value in Redis with a specified key and TTL.
func Set(key string, value string, ttl time.Duration) error {
	// Get the Redis client from the context
	client := cache.GetRedisClient()
	if client == nil {
		return fmt.Errorf("redis client is nil")
	}

	return client.Set(context.Background(), key, value, ttl).Err()
}

// Get retrieves a string value from Redis with a specified key.
func Get(key string) (string, error) {
	// Get the Redis client from the context
	client := cache.GetRedisClient()
	if client == nil {
		return "", fmt.Errorf("redis client is nil")
	}

	value, err := client.Get(context.Background(), key).Result()
	if err != nil {
		return "", err
	}
	return value, nil
}

// DeleteKey deletes a key from Redis.
func DeleteKey(key string) error {
	// Get the Redis client from the context
	client := cache.GetRedisClient()
	if client == nil {
		return fmt.Errorf("redis client is nil")
	}

	return client.Del(context.Background(), key).Err()
}
