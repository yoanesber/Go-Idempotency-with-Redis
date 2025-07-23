package cache

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/yoanesber/go-idempotency-with-redis/pkg/logger"

	"github.com/go-redis/redis/v8" // Redis client for Go
)

var (
	once        sync.Once
	RedisClient *redis.Client
	RedisDB     string
	RedisHost   string
	RedisPort   string
	RedisUser   string
	RedisPass   string
	IsFlushDB   string
)

// LoadRedisEnv loads Redis configuration from environment variables
func LoadRedisEnv() bool {
	RedisDB = os.Getenv("REDIS_DB")
	RedisHost = os.Getenv("REDIS_HOST")
	RedisPort = os.Getenv("REDIS_PORT")
	RedisUser = os.Getenv("REDIS_USER")
	RedisPass = os.Getenv("REDIS_PASS")
	IsFlushDB = os.Getenv("REDIS_FLUSH_DB")

	if RedisDB == "" || RedisHost == "" || RedisPort == "" {
		logger.Panic("One or more required environment variables for Redis are not set", nil)
		return false
	}

	return true
}

// InitRedis initializes the Redis client using environment variables
// It constructs the connection string and calls ConnectRedis to establish the connection
func InitRedis() bool {
	isSuccess := true
	once.Do(func() {
		if !LoadRedisEnv() {
			isSuccess = false
			return
		}

		logger.Info("Connecting to Redis...", nil)

		// Initialize the Redis client
		redisDb, _ := strconv.Atoi(RedisDB)
		RedisClient = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", RedisHost, RedisPort),
			Username: RedisUser,
			Password: RedisPass,
			DB:       redisDb,
			// DialTimeout:        10 * time.Second,
			// ReadTimeout:        30 * time.Second,
			// WriteTimeout:       30 * time.Second,
			// PoolSize:           10,
			// PoolTimeout:        30 * time.Second,
			// IdleTimeout:        500 * time.Millisecond,
			// IdleCheckFrequency: 500 * time.Millisecond,
			// TLSConfig: &tls.Config{
			// 	InsecureSkipVerify: true,
			// },
		})

		_, err := RedisClient.Ping(context.Background()).Result()
		if err != nil {
			logger.Fatal(fmt.Sprintf("Failed to connect to Redis: %v", err), nil)
			isSuccess = false
			return
		}

		logger.Info("Connected to Redis", nil)

		// Flush all keys in the Redis database
		// This is typically used for testing or development purposes
		if IsFlushDB == "TRUE" {
			logger.Info("Flushing Redis database...", nil)
			if status, err := RedisClient.FlushDBAsync(context.Background()).Result(); err != nil {
				logger.Error(fmt.Sprintf("Failed to flush Redis database: %v", err), nil)
				isSuccess = false
				return
			} else {
				logger.Info(fmt.Sprintf("Redis database flushed successfully: %s", status), nil)
			}
		}
	})

	return isSuccess
}

// GetRedisClient retrieves the Redis client instance
// If the client is not initialized, it calls InitRedis to set it up
func GetRedisClient() *redis.Client {
	if RedisClient == nil {
		if !InitRedis() {
			logger.Error("Failed to initialize Redis client", nil)
			return nil
		}
	}

	return RedisClient
}

// CloseRedis closes the Redis client connection
func CloseRedis() {
	if RedisClient != nil {
		if err := RedisClient.Close(); err != nil {
			logger.Error(fmt.Sprintf("Failed to close Redis client: %v", err), nil)
		}

		RedisClient = nil
		logger.Info("Redis client closed successfully", nil)
		return
	}

	once = sync.Once{} // Reset the once to allow re-initialization
	RedisClient = nil  // Clear the RedisClient variable to prevent further use
	logger.Warn("Redis client is nil, nothing to close", nil)
}
