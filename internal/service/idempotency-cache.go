package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"gorm.io/gorm"

	"github.com/yoanesber/go-idempotency-api/config/database"
	"github.com/yoanesber/go-idempotency-api/internal/entity"
	"github.com/yoanesber/go-idempotency-api/internal/repository"
	metacontext "github.com/yoanesber/go-idempotency-api/pkg/context-data/meta-context"
	redisutil "github.com/yoanesber/go-idempotency-api/pkg/util/redis-util"
)

const (
	ttl_hour = 24 // Time to live for idempotency keys in Redis, set to 24 hours
)

// Interface for idempotency key service
// This interface defines the methods that the idempotency key service should implement
type IdempotencyCacheService interface {
	GetAllIdempotencyCaches() ([]entity.IdempotencyCache, error)
	GetIdempotencyCacheByKey(key string) (entity.IdempotencyCache, error)
	CreateIdempotencyCache(ctx context.Context, responsePayload interface{}) (entity.IdempotencyCache, error)
	UpdateIdempotencyCache(key string, responsePayload interface{}) (entity.IdempotencyCache, error)
}

// This struct defines the IdempotencyCacheService that contains a repository field of type IdempotencyCacheRepository
// It implements the IdempotencyCacheService interface and provides methods for idempotency key-related operations
type idempotencyCacheService struct {
	repo repository.IdempotencyCacheRepository
}

// NewIdempotencyCacheService creates a new instance of IdempotencyCacheService with the given repository.
// It initializes the idempotencyCacheService struct and returns it.
func NewIdempotencyCacheService(repo repository.IdempotencyCacheRepository) IdempotencyCacheService {
	return &idempotencyCacheService{repo: repo}
}

// GetAllIdempotencyCaches retrieves all idempotency keys from the database.
func (s *idempotencyCacheService) GetAllIdempotencyCaches() ([]entity.IdempotencyCache, error) {
	db := database.GetPostgres()
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	// Retrieve all idempotency keys from the repository
	idempotencyCaches, err := s.repo.GetAllIdempotencyCaches(db)
	if err != nil {
		return nil, err
	}

	return idempotencyCaches, nil
}

// GetIdempotencyCacheByKey retrieves an idempotency key by its key from the database.
func (s *idempotencyCacheService) GetIdempotencyCacheByKey(key string) (entity.IdempotencyCache, error) {
	db := database.GetPostgres()
	if db == nil {
		return entity.IdempotencyCache{}, fmt.Errorf("database connection is nil")
	}

	// Retrieve the idempotency key by key from the repository
	idempotencyCache, err := s.repo.GetIdempotencyCacheByKey(db, key)
	if err != nil {
		return entity.IdempotencyCache{}, err
	}

	return idempotencyCache, nil
}

// CreateIdempotencyCache creates a new idempotency key in the database.
func (s *idempotencyCacheService) CreateIdempotencyCache(ctx context.Context, responsePayload interface{}) (entity.IdempotencyCache, error) {
	db := database.GetPostgres()
	if db == nil {
		return entity.IdempotencyCache{}, fmt.Errorf("database connection is nil")
	}

	// Extract the idempotency key and body hash from the context
	meta, ok := metacontext.ExtractIdemCompetencyMeta(ctx)
	if !ok {
		return entity.IdempotencyCache{}, fmt.Errorf("idempotency metadata not found in context")
	}

	// Create a new idempotency key object
	idemKey := meta.Key
	bodyHash := meta.BodyHash
	resp, err := json.Marshal(responsePayload)
	if err != nil {
		return entity.IdempotencyCache{}, fmt.Errorf("failed to marshal response payload: %w", err)
	}

	// Convert the response payload to JSON string
	respStr := string(resp)

	// Create a new IdempotencyCache object
	now := time.Now()
	idemData := entity.IdempotencyCache{
		Key:             idemKey,
		BodyHash:        bodyHash,
		ResponsePayload: respStr,
		CreatedAt:       now,
	}

	// Set the expiration time for the idempotency key
	ttlStr := os.Getenv("IDEMPOTENCY_TTL_HOURS")
	ttl, err := strconv.Atoi(ttlStr)
	if err != nil {
		return entity.IdempotencyCache{}, fmt.Errorf("invalid Idempotency TTL hours: %w", err)
	}
	idemData.ExpiredAt = now.Add(time.Duration(ttl) * time.Hour)

	createdIdemData := entity.IdempotencyCache{}
	err = db.Transaction(func(tx *gorm.DB) error {
		// Check if the idempotency key already exists
		existingIdem, err := s.repo.GetIdempotencyCacheByKey(tx, idemKey)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		// If the key already exists, return an error
		if existingIdem.Key != "" {
			return fmt.Errorf("idempotency key %s already exists", idemKey)
		}

		// Create the new idempotency key
		createdIdemData, err = s.repo.CreateIdempotencyCache(tx, idemData)
		if err != nil {
			return err
		}

		// Store the idempotency key and body hash in Redis with a TTL
		ttl := time.Duration(ttl_hour * time.Hour)
		idemPrefix := os.Getenv("IDEMPOTENCY_PREFIX")
		redisKey := idemPrefix + idemKey
		if err := redisutil.SetJSON(redisKey, createdIdemData, ttl); err != nil {
			return fmt.Errorf("failed to set idempotency key in Redis: %w", err)
		}

		return nil
	})

	if err != nil {
		return entity.IdempotencyCache{}, err
	}

	return createdIdemData, nil
}

// UpdateIdempotencyCache updates an existing idempotency key in the database.
func (s *idempotencyCacheService) UpdateIdempotencyCache(key string, responsePayload interface{}) (entity.IdempotencyCache, error) {
	db := database.GetPostgres()
	if db == nil {
		return entity.IdempotencyCache{}, fmt.Errorf("database connection is nil")
	}

	// Convert the response payload to JSON string
	resp, err := json.Marshal(responsePayload)
	if err != nil {
		return entity.IdempotencyCache{}, fmt.Errorf("failed to marshal response payload: %w", err)
	}
	respStr := string(resp)

	updatedIdemData := entity.IdempotencyCache{}
	err = db.Transaction(func(tx *gorm.DB) error {
		// Retrieve the existing idempotency key
		existingIdem, err := s.repo.GetIdempotencyCacheByKey(db, key)
		if err != nil {
			return err
		}

		if existingIdem.Key == "" {
			return fmt.Errorf("idempotency key %s does not exist", key)
		}

		// Update the idempotency key with the new response payload
		existingIdem.ResponsePayload = respStr
		existingIdem.UpdatedAt = time.Now()

		// Update the idempotency key in the database
		updatedIdemData, err = s.repo.UpdateIdempotencyCache(tx, existingIdem)
		if err != nil {
			return err
		}

		// Update the idempotency key in Redis with a TTL
		ttl := time.Duration(ttl_hour * time.Hour)
		idemPrefix := os.Getenv("IDEMPOTENCY_PREFIX")
		redisKey := idemPrefix + key
		if err := redisutil.SetJSON(redisKey, updatedIdemData, ttl); err != nil {
			return fmt.Errorf("failed to update idempotency key in Redis: %w", err)
		}

		return nil
	})

	if err != nil {
		return entity.IdempotencyCache{}, err
	}

	return updatedIdemData, nil
}
