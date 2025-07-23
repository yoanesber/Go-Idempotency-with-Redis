package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"

	"github.com/yoanesber/go-idempotency-with-redis/internal/service"
	httputil "github.com/yoanesber/go-idempotency-with-redis/pkg/util/http-util"
)

// This struct defines the DataRedisHandler which handles HTTP requests related to Redis data.
// It contains a service field of type DataRedisService which is used to interact with the Redis data layer.
type DataRedisHandler struct {
	Service service.DataRedisService
}

// NewDataRedisHandler creates a new instance of DataRedisHandler.
// It initializes the DataRedisHandler struct with the provided DataRedisService.
func NewDataRedisHandler(dataRedisService service.DataRedisService) *DataRedisHandler {
	return &DataRedisHandler{Service: dataRedisService}
}

// GetStringValue retrieves a string value from Redis by its key and returns it as JSON.
// @Summary      Get string value from Redis
// @Description  Get a string value from Redis by its key
// @Tags         dataredis
// @Accept       json
// @Produce      json
// @Param        key   path      string  true  "Redis key"
// @Success      200  {object}  HttpResponse for successful retrieval
// @Failure      400  {object}  HttpResponse for bad request
// @Failure      404  {object}  HttpResponse for not found
// @Failure      500  {object}  HttpResponse for internal server error
// @Router       /dataredis/string/{key} [get]
func (h *DataRedisHandler) GetStringValue(c *gin.Context) {
	// Parse the key from the URL parameter
	key := c.Param("key")
	if key == "" {
		httputil.BadRequest(c, "Invalid key", "Key cannot be empty")
		return
	}

	// Call the service to get the string value from Redis
	value, err := h.Service.GetStringValue(key)
	if err == redis.Nil {
		httputil.NotFound(c, "Value not found", "Key does not exist in Redis")
		return
	}

	if err != nil {
		httputil.InternalServerError(c, "Failed to get string value", err.Error())
		return
	}

	// Check if the value is empty
	if value == "" {
		httputil.NotFound(c, "Value not found", "Value is empty")
		return
	}

	// Return the string value as JSON
	httputil.Success(c, "String value retrieved successfully", value)
}

// GetJSONValue retrieves a JSON value from Redis by its key and returns it as JSON.
// @Summary      Get JSON value from Redis
// @Description  Get a JSON value from Redis by its key
// @Tags         dataredis
// @Accept       json
// @Produce      json
// @Param        key   path      string  true  "Redis key"
// @Success      200  {object}  HttpResponse for successful retrieval
// @Failure      400  {object}  HttpResponse for bad request
// @Failure      404  {object}  HttpResponse for not found
// @Failure      500  {object}  HttpResponse for internal server error
// @Router       /dataredis/json/{key} [get]
func (h *DataRedisHandler) GetJSONValue(c *gin.Context) {
	// Parse the key from the URL parameter
	key := c.Param("key")
	if key == "" {
		httputil.BadRequest(c, "Invalid key", "Key cannot be empty")
		return
	}

	// Call the service to get the JSON value from Redis
	value, err := h.Service.GetJSONValue(key)
	if err == redis.Nil {
		httputil.NotFound(c, "Value not found", "Key does not exist in Redis")
		return
	}

	if err != nil {
		httputil.InternalServerError(c, "Failed to get JSON value", err.Error())
		return
	}

	// Check if the value is empty
	if value == nil {
		httputil.NotFound(c, "Value not found", "Value is empty")
		return
	}

	// Return the JSON value as JSON
	httputil.Success(c, "JSON value retrieved successfully", value)
}
