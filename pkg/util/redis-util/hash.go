package redis_util

import (
	"context"
	"fmt"

	"github.com/yoanesber/go-idempotency-api/config/cache"
)

// SetHashField sets a field in a Redis hash with a specified key and value.
// It adds the field to the hash if it doesn't exist, or updates it if it does.
func SetHashField(key, field, value string) error {
	// Get the Redis client from the context
	client := cache.GetRedisClient()
	if client == nil {
		return fmt.Errorf("redis client is nil")
	}

	return client.HSet(context.Background(), key, field, value).Err()
}

// GetHashField retrieves a field from a Redis hash with a specified key.
// It returns the value of the field if it exists, or an error if it doesn't.
func GetHashField(key, field string) (string, error) {
	// Get the Redis client from the context
	client := cache.GetRedisClient()
	if client == nil {
		return "", fmt.Errorf("redis client is nil")
	}

	return client.HGet(context.Background(), key, field).Result()
}

// GetAllHash retrieves all fields and values from a Redis hash with a specified key.
// It returns a map of field-value pairs.
func GetAllHash(key string) (map[string]string, error) {
	// Get the Redis client from the context
	client := cache.GetRedisClient()
	if client == nil {
		return nil, fmt.Errorf("redis client is nil")
	}

	return client.HGetAll(context.Background(), key).Result()
}
