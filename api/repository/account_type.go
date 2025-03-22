package repository

import (
	"context"

	interface_repository "github.com/riad/banksystemendtoend/api/interface/repository"
	db "github.com/riad/banksystemendtoend/db/sqlc"
)

type accountTypeRepository struct {
	store db.Store
}

func NewAccountTypeRepository(store db.Store) interface_repository.AccountTypeRepository {
	return &accountTypeRepository{
		store: store,
	}
}

func (r *accountTypeRepository) CreateAccountType(ctx context.Context, arg db.CreateAccountTypeParams) (db.AccountType, error) {
	return r.store.CreateAccountType(ctx, arg)
}

func (r *accountTypeRepository) GetAccountType(ctx context.Context, accountType string) (db.AccountType, error) {
	return r.store.GetAccountType(ctx, accountType)
}

func (r *accountTypeRepository) ListAccountTypes(ctx context.Context) ([]db.AccountType, error) {
	return r.store.ListAccountTypes(ctx)
}

func (r *accountTypeRepository) UpdateAccountType(ctx context.Context, arg db.UpdateAccountTypeParams) (db.AccountType, error) {
	return r.store.UpdateAccountType(ctx, arg)
}

func (r *accountTypeRepository) DeleteAccountType(ctx context.Context, accountType string) error {
	return r.store.DeleteAccountType(ctx, accountType)
}
