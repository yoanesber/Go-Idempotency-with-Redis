package test_consumer

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/yoanesber/go-idempotency-with-redis/internal/entity"
	"github.com/yoanesber/go-idempotency-with-redis/internal/handler"
	"github.com/yoanesber/go-idempotency-with-redis/internal/service"
	"github.com/yoanesber/go-idempotency-with-redis/pkg/customtype"
	httputil "github.com/yoanesber/go-idempotency-with-redis/pkg/util/http-util"
)

func TestGetConsumers(t *testing.T) {
	// Define a mocked repository, service, and handler
	// This will allow us to test the handler without needing a real database connection
	r := NewConsumerMockedRepository()
	s := service.NewConsumerService(r)
	h := handler.NewConsumerHandler(s)

	// Set up the Gin router and the route for getting all consumers
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/api/v1/consumers", h.GetAllConsumers)

	// Create a request to the endpoint and record the response
	req, _ := http.NewRequest("GET", "/api/v1/consumers", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check the response status code and body
	assert.Equal(t, http.StatusOK, w.Code)

	// Unmarshal the response body into a HttpResponse struct
	// This struct is used to standardize the response format
	var httpResponse httputil.HttpResponse
	err := json.Unmarshal(w.Body.Bytes(), &httpResponse)
	assert.NoError(t, err)
	assert.NotEmpty(t, httpResponse.Data)
	assert.Nil(t, httpResponse.Error)
}

func TestGetConsumerByID(t *testing.T) {
	// Define a mocked repository, service, and handler
	// This will allow us to test the handler without needing a real database connection
	r := NewConsumerMockedRepository()
	s := service.NewConsumerService(r)
	h := handler.NewConsumerHandler(s)

	// Set up the Gin router and the route for getting a consumer by ID
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/api/v1/consumers/:id", h.GetConsumerByID)

	// Create a request to the endpoint with a specific consumer ID and record the response
	id := "dummy-id" // Assuming we have a consumer with ID 1 in our mocked repository
	req, _ := http.NewRequest("GET", "/api/v1/consumers/"+id, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check the response status code and body
	assert.Equal(t, http.StatusOK, w.Code)

	// Unmarshal the response body into a HttpResponse struct
	// This struct is used to standardize the response format
	var httpResponse httputil.HttpResponse
	err := json.Unmarshal(w.Body.Bytes(), &httpResponse)
	assert.NoError(t, err)
	assert.NotEmpty(t, httpResponse.Data)
	assert.Nil(t, httpResponse.Error)
	assert.Equal(t, id, httpResponse.Data.(map[string]interface{})["id"].(string), "Expected consumer ID to match")
}

func TestGetConsumerByID_NotFound(t *testing.T) {
	// Define a mocked repository, service, and handler
	// This will allow us to test the handler without needing a real database connection
	r := NewConsumerMockedRepository()
	s := service.NewConsumerService(r)
	h := handler.NewConsumerHandler(s)

	// Set up the Gin router and the route for getting a consumer by ID
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/api/v1/consumers/:id", h.GetConsumerByID)

	// Create a request to the endpoint with a non-existent consumer ID and record the response
	id := "non-existent-id"
	req, _ := http.NewRequest("GET", "/api/v1/consumers/"+id, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check the response status code and body
	assert.Equal(t, http.StatusNotFound, w.Code)

	// Unmarshal the response body into a HttpResponse struct
	var httpResponse httputil.HttpResponse
	err := json.Unmarshal(w.Body.Bytes(), &httpResponse)
	assert.NoError(t, err)
	assert.Empty(t, httpResponse.Data)
	assert.NotNil(t, httpResponse.Error)
}

func TestCreateConsumer(t *testing.T) {
	// Define a mocked repository, service, and handler
	// This will allow us to test the handler without needing a real database connection
	r := NewConsumerMockedRepository()
	s := service.NewConsumerService(r)
	h := handler.NewConsumerHandler(s)

	// Set up the Gin router and the route for creating a consumer
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/api/v1/consumers", h.CreateConsumer)

	// Create a request to the endpoint with a new consumer's data
	// This data will be used to create a new consumer in the mocked repository
	newConsumer := entity.Consumer{
		Fullname:  "John Doe",
		Username:  "johndoe",
		Email:     "john.doe@example.com",
		Phone:     "1234567890",
		Address:   "123 Dummy Street",
		BirthDate: &customtype.Date{Time: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)},
	}

	// Marshal the new consumer data into JSON format and create a request
	// This request will be sent to the endpoint to create a new consumer
	reqBody, _ := json.Marshal(newConsumer)
	req, _ := http.NewRequest("POST", "/api/v1/consumers", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check the response status code and body
	assert.Equal(t, http.StatusCreated, w.Code)

	// Unmarshal the response body into a HttpResponse struct
	var httpResponse httputil.HttpResponse
	err := json.Unmarshal(w.Body.Bytes(), &httpResponse)
	assert.NoError(t, err)
	assert.NotEmpty(t, httpResponse.Data)
	assert.Nil(t, httpResponse.Error)

	// Check if the created consumer's data matches the input data
	createdConsumer := httpResponse.Data.(map[string]interface{})
	assert.Equal(t, newConsumer.Fullname, createdConsumer["fullname"])
	assert.Equal(t, newConsumer.Username, createdConsumer["username"])
	assert.Equal(t, newConsumer.Email, createdConsumer["email"])
}

func TestCreateConsumer_ValidationError(t *testing.T) {
	// Define a mocked repository, service, and handler
	// This will allow us to test the handler without needing a real database connection
	r := NewConsumerMockedRepository()
	s := service.NewConsumerService(r)
	h := handler.NewConsumerHandler(s)

	// Set up the Gin router and the route for creating a consumer
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/api/v1/consumers", h.CreateConsumer)

	// Create a request to the endpoint with invalid consumer data (missing required fields)
	invalidConsumer := entity.Consumer{
		Fullname: "", // Missing fullname
		Username: "johndoe",
	}

	reqBody, _ := json.Marshal(invalidConsumer)
	req, _ := http.NewRequest("POST", "/api/v1/consumers", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check the response status code and body
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Unmarshal the response body into a HttpResponse struct
	var httpResponse httputil.HttpResponse
	err := json.Unmarshal(w.Body.Bytes(), &httpResponse)
	assert.NoError(t, err)
	assert.Empty(t, httpResponse.Data)
	assert.NotNil(t, httpResponse.Error)
}

func TestUpdateConsumerStatus(t *testing.T) {
	// Define a mocked repository, service, and handler
	// This will allow us to test the handler without needing a real database connection
	r := NewConsumerMockedRepository()
	s := service.NewConsumerService(r)
	h := handler.NewConsumerHandler(s)

	// Set up the Gin router and the route for updating a consumer's status
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.PATCH("/api/v1/consumers/:id", h.UpdateConsumerStatus)

	// Create a request to the endpoint with a specific consumer ID and new status
	id := "dummy-id" // Assuming we have a consumer with ID 1 in our mocked repository
	newStatus := "inactive"
	req, _ := http.NewRequest("PATCH", "/api/v1/consumers/"+id+"?status="+newStatus, nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check the response status code and body
	assert.Equal(t, http.StatusOK, w.Code)

	// Unmarshal the response body into a HttpResponse struct
	var httpResponse httputil.HttpResponse
	err := json.Unmarshal(w.Body.Bytes(), &httpResponse)
	assert.NoError(t, err)
	assert.NotEmpty(t, httpResponse.Data)
	assert.Nil(t, httpResponse.Error)
}
