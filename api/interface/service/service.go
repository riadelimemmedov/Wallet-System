package interface_service

import (
	"context"
	"mime/multipart"

	"github.com/riad/banksystemendtoend/api/dto"
	db "github.com/riad/banksystemendtoend/db/sqlc"
)

// AccountTypeService defines the interface for account type-related business logic
type AccountTypeService interface {
	// CreateAccountType creates a new account type
	CreateAccountType(ctx context.Context, accountType, description string) (db.AccountType, error)

	// GetAccountType retrieves an account type
	GetAccountType(ctx context.Context, accountType string) (db.AccountType, error)

	// ListAccountTypes retrieves all account types
	ListAccountTypes(ctx context.Context) ([]db.AccountType, error)

	// UpdateAccountType updates an account type
	UpdateAccountType(ctx context.Context, accountType, description string, isActive bool) (db.AccountType, error)

	// DeleteAccountType deletes an account type
	DeleteAccountType(ctx context.Context, accountType string) error
}

// S3Service defines the interface for S3 storage operations
type S3Service interface {
	// UploadFile uploads a file to S3 and returns the URL
	UploadFile(ctx context.Context, file *multipart.FileHeader) (string, error)

	// DeleteFile removes a file from S3
	DeleteFile(ctx context.Context, fileURL string) error
}

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
