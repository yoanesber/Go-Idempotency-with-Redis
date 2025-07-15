package entity

import (
	"time"
)

// IdempotencyCache represents an idempotency key entity.
// It is used to ensure that a request is processed only once, even if it is sent multiple times.
type IdempotencyCache struct {
	Key             string    `gorm:"type:uuid;primaryKey" json:"key" validate:"required,uuid4"`
	BodyHash        string    `gorm:"type:text;not null" json:"bodyHash" validate:"required"`
	ResponsePayload string    `gorm:"type:text;not null" json:"responsePayload" validate:"required"`
	CreatedAt       time.Time `gorm:"type:timestamptz;autoCreateTime;default:now()" json:"createdAt,omitempty"`
	UpdatedAt       time.Time `gorm:"type:timestamptz;autoUpdateTime;default:now()" json:"updatedAt,omitempty"`
	ExpiredAt       time.Time `gorm:"type:timestamptz;not null" json:"expiredAt" validate:"required"`
}

// TableName overrides the table name used by GORM to `idempotency_keys` and `idempotency_logs`.
func (IdempotencyCache) TableName() string {
	return "idempotency_cache"
}

// Equals compares two IdempotencyCache objects for equality.
func (ik *IdempotencyCache) Equals(other *IdempotencyCache) bool {
	if ik == nil && other == nil {
		return true
	}

	if ik == nil || other == nil {
		return false
	}

	if (ik.Key != other.Key) ||
		(ik.BodyHash != other.BodyHash) ||
		(ik.ResponsePayload != other.ResponsePayload) ||
		(ik.CreatedAt != other.CreatedAt) ||
		(ik.UpdatedAt != other.UpdatedAt) ||
		(ik.ExpiredAt != other.ExpiredAt) {
		return false
	}

	return true
}
