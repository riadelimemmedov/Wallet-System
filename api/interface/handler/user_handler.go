package handler_interface

import "github.com/gin-gonic/gin"

// UserHandler defines the interface for user-related HTTP handlers
type UserHandler interface {
	// CreateUser handles the user creation endpoint
	CreateUser(ctx *gin.Context)

	// GetUser handles retrieving a user by ID
	GetUser(ctx *gin.Context)

	// ListUsers handles retrieving a paginated list of users
	ListUsers(ctx *gin.Context)

	// UpdateUser handles updating a user's information
	UpdateUser(ctx *gin.Context)

	// DeleteUser handles soft-deleting a user (deactivating)
	DeleteUser(ctx *gin.Context)

	// HardDeleteUser handles permanently removing a user
	HardDeleteUser(ctx *gin.Context)
}
