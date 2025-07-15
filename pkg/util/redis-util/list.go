package redis_util

import (
	"context"
	"fmt"

	"github.com/yoanesber/go-idempotency-api/config/cache"
)

// PushToList pushes a value to a Redis list with a specified key.
// It adds the value to the head of the list.
func PushToList(key string, value string) error {
	// Get the Redis client from the context
	client := cache.GetRedisClient()
	if client == nil {
		return fmt.Errorf("redis client is nil")
	}

	return client.LPush(context.Background(), key, value).Err()
}

// GetListRange retrieves a range of values from a Redis list with a specified key.
// It returns a slice of strings representing the values in the specified range.
func GetListRange(key string, start int64, stop int64) ([]string, error) {
	// Get the Redis client from the context
	client := cache.GetRedisClient()
	if client == nil {
		return nil, fmt.Errorf("redis client is nil")
	}

	values, err := client.LRange(context.Background(), key, start, stop).Result()
	if err != nil {
		return nil, err
	}
	return values, nil
}

// PopFromList pops a value from a Redis list with a specified key.
// It removes the value from the head of the list and returns the updated list.
// If the list is empty, it returns an empty slice.
func PopFromList(key string) ([]string, error) {
	// Get the Redis client from the context
	client := cache.GetRedisClient()
	if client == nil {
		return nil, fmt.Errorf("redis client is nil")
	}

	_, err := client.LPop(context.Background(), key).Result()
	if err != nil {
		return nil, err
	}

	// Get the updated list after popping the value
	updatedList, err := client.LRange(context.Background(), key, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	return updatedList, nil
}
