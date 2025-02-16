package middleware

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// APIKey struct holds the actual API key.
type APIKey struct {
	apiKey string
}

// NewAPIKey creates a new instance of APIKey by reading the API key from the environment variable.
func NewAPIKey() (*APIKey, error) {
	apiKey := os.Getenv("API_KEY")

	if apiKey == "" {
		return nil, fmt.Errorf("API_KEY environment variable is not set")
	}
	return &APIKey{apiKey: apiKey}, nil
}

// ValidateAPIKey is a middleware function that checks if the provided API key in the request
func (apk *APIKey) ValidateAPIKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")

		if apiKey == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "API key is missing",
				"code":  "MISSING_API_KEY",
			})
			return
		}

		if apiKey != apk.apiKey {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid API key",
				"code":  "INVALID_API_KEY",
			})
			return
		}
		c.Next()
	}
}
