package service

import (
	"context"
	"database/sql"
	"fmt"

	interface_repository "github.com/riad/banksystemendtoend/api/interface/repository"
	interface_service "github.com/riad/banksystemendtoend/api/interface/service"
	db "github.com/riad/banksystemendtoend/db/sqlc"
	logger "github.com/riad/banksystemendtoend/pkg/log"
	"go.uber.org/zap"
)

type accountTypeService struct {
	repo interface_repository.AccountTypeRepository
}

func NewAccountTypeService(repo interface_repository.AccountTypeRepository) interface_service.AccountTypeService {
	return &accountTypeService{repo: repo}
}

func (s *accountTypeService) CreateAccountType(ctx context.Context, accountType, description string) (db.AccountType, error) {
	arg := db.CreateAccountTypeParams{
		AccountType: accountType,
		Description: description,
	}
	createdAccountType, err := s.repo.CreateAccountType(ctx, arg)
	if err != nil {
		logger.GetLogger().Error("failed to create account type", zap.Error(err))
		return db.AccountType{}, fmt.Errorf("failed to create account type: %w", err)
	}
	return createdAccountType, nil
}

func (s *accountTypeService) GetAccountType(ctx context.Context, accountType string) (db.AccountType, error) {
	accountTypeData, err := s.repo.GetAccountType(ctx, accountType)
	if err != nil {
		if err == sql.ErrNoRows {
			return db.AccountType{}, sql.ErrNoRows
		}
		logger.GetLogger().Error("Failed to get account type", zap.Error(err))
		return db.AccountType{}, fmt.Errorf("failed to get account type: %w", err)
	}
	return accountTypeData, nil
}

func (s *accountTypeService) ListAccountTypes(ctx context.Context) ([]db.AccountType, error) {
	accountTypes, err := s.repo.ListAccountTypes(ctx)
	if err != nil {
		logger.GetLogger().Error("Failed to list account types", zap.Error(err))
		return nil, fmt.Errorf("failed to list account types: %w", err)
	}
	return accountTypes, nil
}

func (s *accountTypeService) UpdateAccountType(ctx context.Context, accountType, description string, isActive bool) (db.AccountType, error) {
	arg := db.UpdateAccountTypeParams{
		AccountType: accountType,
		Description: description,
		IsActive:    isActive,
	}
	updatedAccountType, err := s.repo.UpdateAccountType(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			return db.AccountType{}, sql.ErrNoRows
		}
		logger.GetLogger().Error("Failed to update account type", zap.Error(err))
		return db.AccountType{}, fmt.Errorf("failed to update account type: %w", err)
	}
	return updatedAccountType, nil
}

func (s *accountTypeService) DeleteAccountType(ctx context.Context, accountType string) error {
	err := s.repo.DeleteAccountType(ctx, accountType)
	if err != nil {
		if err == sql.ErrNoRows {
			return sql.ErrNoRows
		}
		logger.GetLogger().Error("Failed to delete account type", zap.Error(err))
		return fmt.Errorf("failed to delete account type: %w", err)
	}
	return nil
}
