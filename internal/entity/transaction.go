package entity

import (
	"time"

	"gopkg.in/go-playground/validator.v9"

	validation "github.com/yoanesber/go-idempotency-api/pkg/util/validation-util"
)

const (
	TransactionStatusPending = "pending"
)

// Transaction represents the transaction entity in the database.
type Transaction struct {
	ID                  string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	IdempotencyCacheKey string     `gorm:"type:uuid;not null;unique" json:"idempotencyCacheKey" validate:"required"`
	Type                string     `gorm:"type:varchar(20);not null;check:type IN ('payment','withdrawal','disbursement')" json:"type" validate:"required,max=20,oneof=payment withdrawal disbursement"`
	Amount              float64    `gorm:"type:decimal(10,2);not null" json:"amount" validate:"required,numeric"`
	Status              string     `gorm:"type:varchar(20);not null;check:status IN ('pending','processing','completed','failed')" json:"status"`
	ConsumerID          string     `gorm:"type:uuid;not null" json:"consumerId" validate:"required,uuid4"`
	Consumer            *Consumer  `gorm:"foreignKey:ConsumerID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"consumer,omitempty"`
	CreatedAt           *time.Time `gorm:"type:timestamptz;autoCreateTime;default:now()" json:"createdAt,omitempty"`
	UpdatedAt           *time.Time `gorm:"type:timestamptz;autoUpdateTime;default:now()" json:"updatedAt,omitempty"`
}

// Override the TableName method to specify the table name
// in the database. This is optional if you want to use the default naming convention.
func (Transaction) TableName() string {
	return "transactions"
}

// Equals compares two Transaction objects for equality.
func (t *Transaction) Equals(other *Transaction) bool {
	if t == nil && other == nil {
		return true
	}

	if t == nil || other == nil {
		return false
	}

	if (t.ID != other.ID) ||
		(t.IdempotencyCacheKey != other.IdempotencyCacheKey) ||
		(t.Type != other.Type) ||
		(t.Amount != other.Amount) ||
		(t.Status != other.Status) ||
		(t.ConsumerID != other.ConsumerID) {
		return false
	}

	return true
}

// Validate validates the Transaction struct using the validator package.
func (t *Transaction) Validate() error {
	var v *validator.Validate = validation.GetValidator()

	if err := v.Struct(t); err != nil {
		return err
	}
	return nil
}
