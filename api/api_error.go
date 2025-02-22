package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Error constants
var (
	ErrUserExists        = errors.New("user already exists")
	ErrAccountExists     = errors.New("account already exists")
	ErrTransactionFailed = errors.New("transaction failed")
)

// HandleCreateUserAccountError handles errors that occur when creating a user account
func HandleCreateUserAccountError(c *gin.Context, err error) {
	switch err {
	case ErrUserExists:
		c.JSON(http.StatusConflict, gin.H{"error": "Username or email already exists"})
	case ErrAccountExists:
		c.JSON(http.StatusConflict, gin.H{"error": "Account already exists"})
	case ErrTransactionFailed:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction failed"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An unexpected error occurred"})
	}
}

// ErrorResponse returns a JSON response for an error
func ErrorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
