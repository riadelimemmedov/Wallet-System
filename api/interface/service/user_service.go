package service_interface

import (
	"context"

	"github.com/riad/banksystemendtoend/api/dto"
	db "github.com/riad/banksystemendtoend/db/sqlc"
)

// UserService defines the business logic interface for user operations
type UserService interface {
	// CreateUser creates a new user in the system
	CreateUser(ctx context.Context, req dto.CreateUserRequest) (db.User, error)

	// GetUser retrieves a user by ID
	GetUser(ctx context.Context, userID int64) (db.User, error)

	// ListUsers retrieves a paginated list of users
	ListUsers(ctx context.Context, page, pageSize int32) ([]db.User, error)

	// UpdateUser updates a user's information
	UpdateUser(ctx context.Context, userID int64, req dto.UpdateUserRequest) (db.User, error)

	// DeleteUser soft deletes a user (marks as inactive)
	DeleteUser(ctx context.Context, userID int64) error

	// HardDeleteUser permanently removes a user from the system
	HardDeleteUser(ctx context.Context, userID int64) error
}
