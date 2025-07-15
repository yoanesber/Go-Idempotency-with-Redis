package redis_util

import (
	"context"
	"fmt"

	"github.com/yoanesber/go-idempotency-api/config/cache"
)

// AddToSet adds one or more members to a Redis Set
// If the key does not exist, it will be created.
func AddToSet(key string, members ...string) error {
	// Get the Redis client from the context
	client := cache.GetRedisClient()
	if client == nil {
		return fmt.Errorf("redis client is nil")
	}

	return client.SAdd(context.Background(), key, members).Err()
}

// GetSetMembers retrieves all members of a Redis Set
// It returns a slice of strings representing the members of the set.
func GetSetMembers(key string) ([]string, error) {
	// Get the Redis client from the context
	client := cache.GetRedisClient()
	if client == nil {
		return nil, fmt.Errorf("redis client is nil")
	}

	return client.SMembers(context.Background(), key).Result()
}
