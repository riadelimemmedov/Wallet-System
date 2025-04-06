package repository

import (
	"context"

	interface_repository "github.com/riad/banksystemendtoend/api/interface/repository"
	db "github.com/riad/banksystemendtoend/db/sqlc"
	"github.com/riad/banksystemendtoend/pkg/cache"
)

type accountTypeRepository struct {
	store     db.Store
	cacheable *CacheableRepository
}

func NewAccountTypeRepository(store db.Store, cacheService *cache.Service) interface_repository.AccountTypeRepository {
	//? Create a dedicated cache service for account types
	accountTypeCache := cache.NewService(
		cacheService.GetRedisClient(),
		"account_type",
		cacheService.GetDefaultTTL(),
	)

	return &accountTypeRepository{
		store:     store,
		cacheable: NewCacheableRepository(accountTypeCache),
	}
}

func (r *accountTypeRepository) CreateAccountType(ctx context.Context, arg db.CreateAccountTypeParams) (db.AccountType, error) {
	result, err := r.store.CreateAccountType(ctx, arg)
	if err != nil {
		return db.AccountType{}, err
	}
	//Invalidate the cache
	r.cacheable.InvalidateCache(ctx, "")
	return result, nil

}

func (r *accountTypeRepository) GetAccountType(ctx context.Context, accountType string) (db.AccountType, error) {
	var result db.AccountType

	err := r.cacheable.GetCached(ctx, accountType, &result, func() (interface{}, error) {
		return r.store.GetAccountType(ctx, accountType)
	})
	return result, err
}

func (r *accountTypeRepository) ListAccountTypes(ctx context.Context) ([]db.AccountType, error) {
	var result []db.AccountType

	err := r.cacheable.GetCached(ctx, "list_account_types", &result, func() (interface{}, error) {
		return r.store.ListAccountTypes(ctx)
	})
	return result, err
}

func (r *accountTypeRepository) UpdateAccountType(ctx context.Context, arg db.UpdateAccountTypeParams) (db.AccountType, error) {
	result, err := r.store.UpdateAccountType(ctx, arg)
	if err != nil {
		return db.AccountType{}, err
	}
	// Invalidate the cache
	r.cacheable.InvalidateCache(ctx, "")
	return result, nil
}

func (r *accountTypeRepository) DeleteAccountType(ctx context.Context, accountType string) error {
	err := r.store.DeleteAccountType(ctx, accountType)
	if err != nil {
		return err
	}
	// Invalidate the cache
	r.cacheable.InvalidateCache(ctx, "")
	return nil
}
