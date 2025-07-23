package handler

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v9"
	"gorm.io/gorm"

	"github.com/yoanesber/go-idempotency-with-redis/internal/entity"
	"github.com/yoanesber/go-idempotency-with-redis/internal/service"
	httputil "github.com/yoanesber/go-idempotency-with-redis/pkg/util/http-util"
	validation "github.com/yoanesber/go-idempotency-with-redis/pkg/util/validation-util"
)

// This struct defines the TransactionHandler which handles HTTP requests related to transactions.
// It contains a service field of type TransactionService which is used to interact with the transaction data layer.
type TransactionHandler struct {
	Service service.TransactionService
}

// NewTransactionHandler creates a new instance of TransactionHandler.
// It initializes the TransactionHandler struct with the provided TransactionService.
func NewTransactionHandler(transactionService service.TransactionService) *TransactionHandler {
	return &TransactionHandler{Service: transactionService}
}

// GetAllTransactions retrieves all transactions from the database and returns them as JSON.
// @Summary      Get all transactions
// @Description  Get all transactions from the database
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Param        page   query     string  false "Page number (default is 1)"
// @Param        limit  query     string  false "Number of transactions per page (default is 10)"
// @Success      200  {array}   model.HttpResponse for successful retrieval
// @Failure      400  {object}  model.HttpResponse for bad request
// @Failure      404  {object}  model.HttpResponse for not found
// @Failure      500  {object}  model.HttpResponse for internal server error
// @Router       /transactions [get]
func (h *TransactionHandler) GetAllTransactions(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		httputil.BadRequest(c, "Invalid page number", "Page must be a positive integer")
		return
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		httputil.BadRequest(c, "Invalid limit", "Limit must be a positive integer")
		return
	}

	transactions, err := h.Service.GetAllTransactions(page, limit)
	if err != nil {
		httputil.InternalServerError(c, "Failed to retrieve transactions", err.Error())
		return
	}

	if len(transactions) == 0 {
		httputil.NotFound(c, "No transactions found", "No transactions available in the database")
		return
	}

	httputil.Success(c, "All transactions retrieved successfully", transactions)
}

// GetTransactionByID retrieves a transaction by its ID from the database and returns it as JSON.
// @Summary      Get transaction by ID
// @Description  Get a transaction by its ID from the database
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Transaction ID"
// @Success      200  {object}  model.HttpResponse for successful retrieval
// @Failure      400  {object}  model.HttpResponse for bad request
// @Failure      404  {object}  model.HttpResponse for not found
// @Failure      500  {object}  model.HttpResponse for internal server error
// @Router       /transactions/{id} [get]
func (h *TransactionHandler) GetTransactionByID(c *gin.Context) {
	// Parse the ID from the URL parameter
	id := c.Param("id")
	if id == "" {
		httputil.BadRequest(c, "Invalid ID", "ID cannot be empty")
		return
	}

	// Retrieve the transaction by ID from the service
	transaction, err := h.Service.GetTransactionByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			httputil.NotFound(c, "Transaction not found", "No transaction found with the given ID")
			return
		}

		// If the error is not a record not found error, return a generic internal server error
		// This is to avoid exposing internal details of the error
		httputil.InternalServerError(c, "Failed to retrieve transaction", err.Error())
		return
	}

	httputil.Success(c, "Transaction retrieved successfully", transaction)
}

// CreateTransaction creates a new transaction in the database and returns it as JSON.
// @Summary      Create transaction
// @Description  Create a new transaction in the database
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Param        transaction  body      Transaction  true  "Transaction object"
// @Success      201  {object}  model.HttpResponse for successful creation
// @Failure      400  {object}  model.HttpResponse for bad request
// @Failure      500  {object}  model.HttpResponse for internal server error
// @Router       /transactions [post]
func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	// Bind the JSON request body to the Transaction struct
	// This will automatically validate the request body against the struct tags
	var transaction entity.Transaction
	if err := c.ShouldBindJSON(&transaction); err != nil {
		httputil.BadRequest(c, "Invalid request body", err.Error())
		return
	}

	// Create the transaction using the service
	createdTransaction, err := h.Service.CreateTransaction(c.Request.Context(), transaction)
	if err != nil {
		// Check if the error is a validation error
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			httputil.BadRequestMap(c, "Failed to create transaction", validation.FormatValidationErrors(err))
			return
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			httputil.NotFound(c, "Consumer not found", "No consumer found with the given ID")
			return
		}

		if errors.Is(err, gorm.ErrInvalidData) {
			httputil.BadRequest(c, "Invalid transaction data", "Transaction data is invalid, this could be due to missing required fields or incorrect data types")
			return
		}

		// If the error is not a record not found and not a validation error, return a generic internal server error
		// This is to avoid exposing internal details of the error
		httputil.InternalServerError(c, "Failed to create transaction", err.Error())
		return
	}

	httputil.Created(c, "Transaction created successfully", createdTransaction)
}
