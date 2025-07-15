package repository

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/yoanesber/go-idempotency-api/internal/entity"
)

// Interface for idempotency key repository
// This interface defines the methods that the idempotency key repository should implement
type IdempotencyCacheRepository interface {
	GetAllIdempotencyCaches(tx *gorm.DB) ([]entity.IdempotencyCache, error)
	GetIdempotencyCacheByKey(tx *gorm.DB, key string) (entity.IdempotencyCache, error)
	CreateIdempotencyCache(tx *gorm.DB, key entity.IdempotencyCache) (entity.IdempotencyCache, error)
	UpdateIdempotencyCache(tx *gorm.DB, key entity.IdempotencyCache) (entity.IdempotencyCache, error)
}

// This struct defines the IdempotencyCacheRepository that contains methods for interacting with the database
// It implements the IdempotencyCacheRepository interface and provides methods for idempotency key-related operations
type idempotencyCacheRepository struct{}

// NewIdempotencyCacheRepository creates a new instance of IdempotencyCacheRepository.
// It initializes the idempotencyCacheRepository struct and returns it.
func NewIdempotencyCacheRepository() IdempotencyCacheRepository {
	return &idempotencyCacheRepository{}
}

// GetAllIdempotencyCaches retrieves all idempotency keys from the database.
func (r *idempotencyCacheRepository) GetAllIdempotencyCaches(tx *gorm.DB) ([]entity.IdempotencyCache, error) {
	// Select all idempotency keys from the database
	var idempotencyCaches []entity.IdempotencyCache
	err := tx.Find(&idempotencyCaches).Error
	if err != nil {
		return nil, err
	}

	return idempotencyCaches, nil
}

// GetIdempotencyCacheByKey retrieves an idempotency key by its key string from the database.
func (r *idempotencyCacheRepository) GetIdempotencyCacheByKey(tx *gorm.DB, key string) (entity.IdempotencyCache, error) {
	// Select the idempotency key with the given key string from the database
	var idempotencyCache entity.IdempotencyCache
	err := tx.First(&idempotencyCache, "key = ?", key).Error
	if err != nil {
		return entity.IdempotencyCache{}, err
	}

	return idempotencyCache, nil
}

// CreateIdempotencyCache creates a new idempotency key in the database.
func (r *idempotencyCacheRepository) CreateIdempotencyCache(tx *gorm.DB, key entity.IdempotencyCache) (entity.IdempotencyCache, error) {
	// Create a new idempotency key in the database
	if err := tx.Create(&key).Error; err != nil {
		return entity.IdempotencyCache{}, fmt.Errorf("failed to create idempotency cache: %w", err)
	}

	return key, nil
}

// UpdateIdempotencyCache updates an existing idempotency key in the database.
func (r *idempotencyCacheRepository) UpdateIdempotencyCache(tx *gorm.DB, key entity.IdempotencyCache) (entity.IdempotencyCache, error) {
	// Update the idempotency key in the database
	if err := tx.Save(&key).Error; err != nil {
		return entity.IdempotencyCache{}, fmt.Errorf("failed to update idempotency cache: %w", err)
	}

	return key, nil
}
