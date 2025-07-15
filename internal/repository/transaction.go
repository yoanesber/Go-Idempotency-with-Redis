package repository

import (
	"fmt"

	"github.com/yoanesber/go-idempotency-api/internal/entity"
	"gorm.io/gorm" // Import GORM for ORM functionalities
)

// Interface for transaction repository
// This interface defines the methods that the transaction repository should implement
type TransactionRepository interface {
	GetAllTransactions(tx *gorm.DB, page int, limit int) ([]entity.Transaction, error)
	GetTransactionByID(tx *gorm.DB, id string) (entity.Transaction, error)
	CreateTransaction(tx *gorm.DB, d entity.Transaction) (entity.Transaction, error)
}

// This struct defines the transactionRepository that implements the TransactionRepository interface.
// It contains methods for interacting with the transaction data in the database.
type transactionRepository struct{}

// NewTransactionRepository creates a new instance of TransactionRepository.
// It initializes the transactionRepository struct and returns it.
func NewTransactionRepository() TransactionRepository {
	return &transactionRepository{}
}

// GetAllTransactions retrieves all transactions from the database.
func (r *transactionRepository) GetAllTransactions(tx *gorm.DB, page int, limit int) ([]entity.Transaction, error) {
	var transactions []entity.Transaction
	err := tx.Order("created_at ASC").
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&transactions).Error

	if err != nil {
		return nil, err
	}

	return transactions, nil
}

// It returns a single transaction by its ID from the database.
func (r *transactionRepository) GetTransactionByID(tx *gorm.DB, id string) (entity.Transaction, error) {
	var transaction entity.Transaction
	err := tx.First(&transaction, "id = ?", id).Error

	if err != nil {
		return entity.Transaction{}, err
	}

	return transaction, nil
}

func (r *transactionRepository) GetAllTransactionsByStatus(tx *gorm.DB, status string, page int, limit int) ([]entity.Transaction, error) {
	var transactions []entity.Transaction
	err := tx.Where("status = ?", status).
		Order("created_at ASC").
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&transactions).
		Error

	if err != nil {
		return nil, err
	}

	return transactions, nil
}

// GetTransactionByConsumerByStatus retrieves a transaction by consumer ID and status from the database.
// It returns a single transaction that matches the consumer ID and status, paginated by page and limit.
func (r *transactionRepository) GetAllTransactionsByConsumerByStatus(tx *gorm.DB, consumerId string, status string, page int, limit int) ([]entity.Transaction, error) {
	var transactions []entity.Transaction
	err := tx.Where("consumer_id = ? AND status = ?", consumerId, status).
		Order("created_at ASC").
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&transactions).
		Error

	if err != nil {
		return nil, err
	}

	return transactions, nil
}

// CreateTransaction creates a new transaction in the database and returns the created transaction.
func (r *transactionRepository) CreateTransaction(tx *gorm.DB, t entity.Transaction) (entity.Transaction, error) {
	// Insert new transaction
	if err := tx.Create(&t).Error; err != nil {
		return entity.Transaction{}, fmt.Errorf("failed to create transaction: %w", err)
	}

	return t, nil
}
