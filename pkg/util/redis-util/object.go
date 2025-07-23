package redis_util

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/yoanesber/go-idempotency-with-redis/config/cache"
)

// SetJSON sets a JSON value in Redis with a specified key and TTL.
// It marshals the value into JSON format and stores it in Redis.
func SetJSON(key string, value interface{}, ttl time.Duration) error {
	// Get the Redis client from the context
	client := cache.GetRedisClient()
	if client == nil {
		return fmt.Errorf("redis client is nil")
	}

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return client.Set(context.Background(), key, data, ttl).Err()
}

// GetJSON retrieves a JSON value from Redis with a specified key.
// It unmarshals the JSON data into the provided value.
func GetJSON[T any](key string) (*T, error) {
	// Get the Redis client from the context
	client := cache.GetRedisClient()
	if client == nil {
		return nil, fmt.Errorf("redis client is nil")
	}

	data, err := client.Get(context.Background(), key).Bytes()
	if err != nil {
		return nil, err
	}

	var result T
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
