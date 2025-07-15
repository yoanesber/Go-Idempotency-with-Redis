package idempotency

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"

	"github.com/yoanesber/go-idempotency-api/internal/entity"
	metacontext "github.com/yoanesber/go-idempotency-api/pkg/context-data/meta-context"
	hashutil "github.com/yoanesber/go-idempotency-api/pkg/util/hash-util"
	httputil "github.com/yoanesber/go-idempotency-api/pkg/util/http-util"
	redisutil "github.com/yoanesber/go-idempotency-api/pkg/util/redis-util"
)

/**
* Enforce is a middleware function that implements idempotency for HTTP requests.
* It checks if the request has an idempotency key and whether the request has already been processed.
* If the request has already been processed, it returns the cached response.
* If the request has not been processed, it injects the idempotency metadata into the context
* and allows the request to proceed to the handler.
 */
func Enforce() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Read the environment variables
		idemEnabled := os.Getenv("IDEMPOTENCY_ENABLED")
		idemKeyHdr := os.Getenv("IDEMPOTENCY_KEY_HEADER")
		idemPrefix := os.Getenv("IDEMPOTENCY_PREFIX")
		if idemEnabled == "" || idemKeyHdr == "" || idemPrefix == "" {
			httputil.InternalServerError(c, "Internal Server Error", "Idempotency enabled, key header, or prefix environment variables are not set")
			c.Abort()
			return
		}

		// Check if idempotency is enabled
		// If the environment variable is not set to "TRUE", skip the middleware
		if idemEnabled != "TRUE" {
			c.Next()
			return
		}

		// Ensure that the request method is POST, PUT, or DELETE
		if c.Request.Method != "POST" && c.Request.Method != "PUT" && c.Request.Method != "DELETE" {
			httputil.MethodNotAllowed(c, "Method Not Allowed", "Idempotency middleware only supports POST, PUT, or DELETE methods")
			c.Abort()
			return
		}

		// Get the idempotency key from the request header
		// The idempotency key is expected to be provided in the request header
		idemKey := c.GetHeader(idemKeyHdr)
		if idemKey == "" {
			httputil.BadRequest(c, "Bad Request", fmt.Sprintf("Idempotency key header '%s' is required", idemKeyHdr))
			c.Abort()
			return
		}

		// Read the request body
		bodyBytes, err := c.GetRawData()
		if err != nil {
			httputil.InternalServerError(c, "Internal Server Error", "Failed to read request body")
			c.Abort()
			return
		}

		// Restore the request body so it can be read again later in the handler
		// This is necessary because reading the body consumes it, and we need it for further processing
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// Hash the request body to create a unique identifier
		bodyHash, err := hashutil.Hash256Bytes(bodyBytes)
		if err != nil {
			httputil.InternalServerError(c, "Internal Server Error", "Failed to hash request body")
			c.Abort()
			return
		}

		// Check if the request has already been processed
		redisKey := idemPrefix + idemKey
		cachedData, err := redisutil.GetJSON[entity.IdempotencyCache](redisKey)
		if err != nil && err != redis.Nil {
			httputil.InternalServerError(c, "Internal Server Error", err.Error())
			c.Abort()
			return
		}

		if cachedData != nil {
			// If idempotency key exists in Redis with different body hash, return conflict error
			if cachedData.BodyHash != bodyHash {
				httputil.Conflict(c, "Conflict", "Request with the same Idempotency-Key but different body has already been processed")
				c.Abort()
				return
			}

			var respPayload any
			if cachedData.ResponsePayload != "" {
				if err := json.Unmarshal([]byte(cachedData.ResponsePayload), &respPayload); err != nil {
					httputil.InternalServerError(c, "Internal Server Error", "Failed to unmarshal cached response payload")
					c.Abort()
					return
				}
			}

			// If the request has already been processed, return the cached response
			httputil.Success(c, "Request already processed", respPayload)
			c.Abort()
			return
		}

		// Inject the idempotency metadata into the context
		// This metadata will be used later to create or update the idempotency key in the database
		meta := metacontext.IdemCompetencyMeta{
			Key:      idemKey,
			BodyHash: bodyHash,
		}
		ctx := metacontext.InjectIdemCompetencyMeta(c.Request.Context(), meta)

		// Set the new request context with idempotency metadata
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
