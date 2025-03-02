package repository_interface

import (
	"context"

	db "github.com/riad/banksystemendtoend/db/sqlc"
)

type AccountTypeRepository interface {
	// CreateAccountType creates a new account type
	CreateAccountType(ctx context.Context, arg db.CreateAccountTypeParams) (db.AccountType, error)

	// GetAccountType retrieves an account type by ID
	GetAccountType(ctx context.Context, accountType string) (db.AccountType, error)

	// ListAccountTypes retrieves all account types
	ListAccountTypes(ctx context.Context) ([]db.AccountType, error)

	// UpdateAccountType updates an account type
	UpdateAccountType(ctx context.Context, arg db.UpdateAccountTypeParams) (db.AccountType, error)

	// DeleteAccountType deletes an account type
	DeleteAccountType(ctx context.Context, accountType string) error
}
