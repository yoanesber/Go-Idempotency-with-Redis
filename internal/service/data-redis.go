package service

import (
	redisutil "github.com/yoanesber/go-idempotency-api/pkg/util/redis-util"
)

// Interface for the DataRedisService
// This interface defines the methods that the DataRedisService should implement
type DataRedisService interface {
	GetStringValue(key string) (string, error)
	GetJSONValue(key string) (interface{}, error)
}

// This struct defines the DataRedisService
type dataRedisService struct{}

// NewDataRedisService creates a new instance of DataRedisService
// It initializes the dataRedisService struct and returns it.
func NewDataRedisService() DataRedisService {
	return &dataRedisService{}
}

// GetStringValue retrieves a string value from Redis by its key
func (s *dataRedisService) GetStringValue(key string) (string, error) {
	value, err := redisutil.Get(key)
	if err != nil {
		return "", err
	}

	return value, nil
}

// GetJSONValue retrieves a JSON value from Redis by its key
func (s *dataRedisService) GetJSONValue(key string) (interface{}, error) {
	value, err := redisutil.GetJSON[any](key)
	if err != nil {
		return nil, err
	}

	return value, nil
}
