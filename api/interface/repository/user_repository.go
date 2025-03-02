package repository_interface

import (
	"context"
	"time"

	db "github.com/riad/banksystemendtoend/db/sqlc"
)

// UserRepository defines the interface for user-related database operations
type UserRepository interface {
	// CreateUser creates a new user with the given parameters
	CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error)

	// GetUser retrieves a user by their ID
	GetUser(ctx context.Context, userID int64) (db.User, error)

	//ListUsers retrieves a list of users with the given limit and offset
	ListUsers(ctx context.Context, limit, offset int32) ([]db.User, error)

	// UpdateUser updates a user with the given parameters
	UpdateUser(ctx context.Context, arg db.UpdateUserParams) (db.User, error)

	// DeleteUser deactivates a user with the given ID
	DeleteUser(ctx context.Context, userID int64) error

	// HardDeleteUser deletes a user with the given ID
	HardDeleteUser(ctx context.Context, userID int64) error

	// UpdateLastLogin updates the last login time of a user with the given ID
	UpdateLastLogin(ctx context.Context, userID int64, time time.Time) error
}
