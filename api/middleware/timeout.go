package middleware

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

// TimeOut is a middleware function that sets a timeout for handling requests.
func TimeOut(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)

		done := make(chan bool, 1)

		go func() {
			c.Next()
			done <- true
		}()

		select {
		case <-ctx.Done():
			if ctx.Err() == context.DeadlineExceeded {
				c.AbortWithStatusJSON(408, gin.H{"error": "Request timeout"})
			}
		case <-done:
			return
		}
	}
}
