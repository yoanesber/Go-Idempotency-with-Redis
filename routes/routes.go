package routes

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"

	"github.com/yoanesber/go-idempotency-api/internal/handler"
	"github.com/yoanesber/go-idempotency-api/internal/repository"
	"github.com/yoanesber/go-idempotency-api/internal/service"
	"github.com/yoanesber/go-idempotency-api/pkg/middleware/headers"
	"github.com/yoanesber/go-idempotency-api/pkg/middleware/idempotency"
	"github.com/yoanesber/go-idempotency-api/pkg/middleware/logging"
	request_filter "github.com/yoanesber/go-idempotency-api/pkg/middleware/request-filter"
	httputil "github.com/yoanesber/go-idempotency-api/pkg/util/http-util"
)

// SetupRouter initializes the router and sets up the routes for the application.
func SetupRouter() *gin.Engine {
	// Create a new Gin router instance
	r := gin.Default()

	// Set up middleware for the router
	// Middleware is used to handle cross-cutting concerns such as logging, security, and request ID generation
	r.Use(
		headers.SecurityHeaders(),
		headers.CorsHeaders(),
		headers.ContentType(),
		request_filter.DetectParameterPollution(),
		logging.RequestLogger(),
		gzip.Gzip(gzip.DefaultCompression),
	)

	// Set up the API version 1 routes
	v1 := r.Group("/api/v1")
	{
		// Routes for consumer management
		// These routes handle CRUD operations for consumers
		consumerGroup := v1.Group("/consumers")
		{
			// Initialize the transaction repository and service
			// This is where the actual implementation of the repository and service would be used
			r := repository.NewConsumerRepository()
			s := service.NewConsumerService(r)

			// Initialize the transaction handler with the service
			// This handler handles the HTTP requests and responses for transaction-related operations
			h := handler.NewConsumerHandler(s)

			// Define the routes for transaction management
			// These routes handle CRUD operations for transactions
			// The GET methods are accessible to both admin and user roles
			consumerGroup.GET("", h.GetAllConsumers)
			consumerGroup.GET("/:id", h.GetConsumerByID)
			consumerGroup.GET("/active", h.GetActiveConsumers)
			consumerGroup.GET("/inactive", h.GetInactiveConsumers)
			consumerGroup.GET("/suspended", h.GetSuspendedConsumers)

			// The POST and PUT methods are restricted to admin users only
			consumerGroup.POST("", h.CreateConsumer)
			consumerGroup.PATCH("/:id", h.UpdateConsumerStatus)
		}

		// Routes for transaction management
		// These routes handle CRUD operations for transactions
		trxGroup := v1.Group("/transactions")
		{
			// Initialize the transaction repository and service
			// This is where the actual implementation of the repository and service would be used
			r := repository.NewTransactionRepository()
			s := service.NewTransactionService(r)

			// Initialize the transaction handler with the service
			// This handler handles the HTTP requests and responses for transaction-related operations
			h := handler.NewTransactionHandler(s)

			// Define the routes for transaction management
			// These routes handle CRUD operations for transactions
			trxGroup.GET("", h.GetAllTransactions)
			trxGroup.GET("/:id", h.GetTransactionByID)

			// The POST and PUT methods are restricted to admin users only
			trxGroup.POST("", idempotency.Enforce(), h.CreateTransaction)
		}
	}

	// NoRoute handler for undefined routes
	// This handler will be called when no other route matches the request
	r.NoRoute(func(c *gin.Context) {
		httputil.NotFound(c, "Not Found", "The requested resource was not found")
	})

	// NoMethod handler for unsupported HTTP methods
	// This handler will be called when a request method is not allowed for the requested resource
	r.NoMethod(func(c *gin.Context) {
		httputil.MethodNotAllowed(c, "Method Not Allowed", "The requested method is not allowed for this resource")
	})

	return r
}
