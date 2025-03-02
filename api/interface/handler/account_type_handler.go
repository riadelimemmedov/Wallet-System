package handler_interface

import "github.com/gin-gonic/gin"

// AccountTypeHandler defines the interface for account type-related HTTP handlers
type AccountTypeHandler interface {
	// CreateAccountType handles creating a new account type
	CreateAccountType(ctx *gin.Context)

	// GetAccountType handles retrieving an account type
	GetAccountType(ctx *gin.Context)

	// ListAccountTypes handles retrieving all account types
	ListAccountTypes(ctx *gin.Context)

	// UpdateAccountType handles updating an account type
	UpdateAccountType(ctx *gin.Context)

	// DeleteAccountType handles deleting an account type
	DeleteAccountType(ctx *gin.Context)
}
