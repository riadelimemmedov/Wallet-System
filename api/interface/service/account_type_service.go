package service_interface

import (
	"context"

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
