package service

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/yoanesber/go-idempotency-api/config/database"
	"github.com/yoanesber/go-idempotency-api/internal/entity"
	"github.com/yoanesber/go-idempotency-api/internal/repository"
	metacontext "github.com/yoanesber/go-idempotency-api/pkg/context-data/meta-context"
)

const (
	paymentEventTopic      = "payment-event"
	withdrawalEventTopic   = "withdrawal-event"
	disbursementEventTopic = "disbursement-event"

	externalServiceURL = "https://example.com/api/transactions" // Example external service URL
)

// Interface for transaction service
// This interface defines the methods that the transaction service should implement
type TransactionService interface {
	GetAllTransactions(page int, limit int) ([]entity.Transaction, error)
	GetTransactionByID(id string) (entity.Transaction, error)
	CreateTransaction(ctx context.Context, t entity.Transaction) (entity.Transaction, error)
}

// This struct defines the TransactionService that contains a repository field of type TransactionRepository
// It implements the TransactionService interface and provides methods for transaction-related operations
type transactionService struct {
	repo repository.TransactionRepository
}

// NewTransactionService creates a new instance of TransactionService with the given repository.
// This function initializes the transactionService struct and returns it.
func NewTransactionService(repo repository.TransactionRepository) TransactionService {
	return &transactionService{repo: repo}
}

// GetAllTransactions retrieves all transactions from the database.
func (s *transactionService) GetAllTransactions(page int, limit int) ([]entity.Transaction, error) {
	db := database.GetPostgres()
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	// Retrieve all transactions from the repository
	transactions, err := s.repo.GetAllTransactions(db, page, limit)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

// GetTransactionByID retrieves a transaction by its ID from the database.
func (s *transactionService) GetTransactionByID(id string) (entity.Transaction, error) {
	db := database.GetPostgres()
	if db == nil {
		return entity.Transaction{}, fmt.Errorf("database connection is nil")
	}

	// Retrieve the transaction by ID from the repository
	transaction, err := s.repo.GetTransactionByID(db, id)
	if err != nil {
		return entity.Transaction{}, err
	}

	return transaction, nil
}

// CreateTransaction creates a new transaction in the database.
// It validates the transaction struct and checks if the ID already exists before creating a new transaction.
func (s *transactionService) CreateTransaction(ctx context.Context, t entity.Transaction) (entity.Transaction, error) {
	db := database.GetPostgres()
	if db == nil {
		return entity.Transaction{}, fmt.Errorf("database connection is nil")
	}

	// Extract the idempotency key and body hash from the context
	meta, ok := metacontext.ExtractIdemCompetencyMeta(ctx)
	if !ok {
		return entity.Transaction{}, fmt.Errorf("idempotency meta not found in context")
	}

	// Set the idempotency cache key in the transaction
	t.IdempotencyCacheKey = meta.Key

	// Validate the transaction struct using the validator
	if err := t.Validate(); err != nil {
		return entity.Transaction{}, err
	}

	createdTransaction := entity.Transaction{}
	err := db.Transaction(func(tx *gorm.DB) error {
		// Check if the transaction is associated with a valid consumer
		if t.ConsumerID == "" {
			return fmt.Errorf("consumer ID is required")
		}

		// Define the consumer repository and service
		consumerRepo := repository.NewConsumerRepository()
		consumerService := NewConsumerService(consumerRepo)

		// Check if the consumer exists
		consumer, err := consumerService.GetConsumerByID(t.ConsumerID)
		if err != nil {
			return err
		}

		if consumer.Status != entity.ConsumerStatusActive {
			return gorm.ErrInvalidData
		}

		// Save the new transaction in the database
		t.Status = entity.TransactionStatusPending // Set initial status to pending
		createdTransaction, err = s.repo.CreateTransaction(tx, t)
		if err != nil {
			return err
		}

		// Save idempotency cache in the database
		idemRepo := repository.NewIdempotencyCacheRepository()
		idemService := NewIdempotencyCacheService(idemRepo)
		if _, err := idemService.CreateIdempotencyCache(ctx, createdTransaction); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return entity.Transaction{}, err
	}

	return createdTransaction, nil
}
