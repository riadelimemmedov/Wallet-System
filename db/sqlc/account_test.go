package db

import (
	"context"
	"strings"
	"testing"

	"github.com/riad/banksystemendtoend/util/common"
	"github.com/stretchr/testify/require"
)

// - TestCreateAccount: Verifies random account creation
func createRandomAccount(t *testing.T) Account {
	sqlStore := SetupTestStore(t)
	require.NotEmpty(t, sqlStore)

	var account Account
	var err error

	overdraftLimit, err := common.RandomNumeric()
	require.NoError(t, err)

	interestRate, err := common.RandomInterestRate()
	require.NoError(t, err)

	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		arg := CreateAccountParams{
			UserID:         createRandomUser(t).UserID,
			AccountNumber:  common.RandomString(15),
			AccountType:    createRandomAccountType(t).AccountType,
			CurrencyCode:   createRandomCurrency(t).CurrencyCode,
			InterestRate:   interestRate,
			OverdraftLimit: overdraftLimit,
		}
		account, err = sqlStore.Queries.CreateAccount(context.Background(), arg)
		if err == nil {
			break
		}
		if strings.Contains(err.Error(), "SQLSTATE 23505") {
			if i == maxRetries-1 {
				t.Fatalf("failed to create unique account after %d attempts: %v", maxRetries, err)
			}
			continue
		}
		require.NoError(t, err)
	}

	require.NotEmpty(t, account)
	return account
}

// - TestCreateAccount: Verifies random account creation
func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
	defer CleanupDB(t)
}

// - TestGetAccount: Ensures account retrieval matches created account
func TestGetAccount(t *testing.T) {
	sqlStore := SetupTestStore(t)
	account1 := createRandomAccount(t)

	account2, err := sqlStore.Queries.GetAccount(context.Background(), account1.AccountID)
	require.NoError(t, err)
	require.Equal(t, account1, account2)

	defer CleanupDB(t)
}

// - TestUpdateAccountBalance: Checks balance updates work correctly
func TestUpdateAccountBalance(t *testing.T) {
	sqlStore := SetupTestStore(t)
	account1 := createRandomAccount(t)

	amount, err := common.RandomNumeric()
	require.NoError(t, err)

	arg := UpdateAccountBalanceParams{
		AccountID: account1.AccountID,
		Amount:    amount,
	}

	account2, err := sqlStore.Queries.UpdateAccountBalance(context.Background(), arg)
	require.NoError(t, err)
	require.NotEqual(t, account1.Balance, account2.Balance)

	defer CleanupDB(t)
}

// - TestDeleteAccount: Verifies soft delete (IsActive flag)
func TestDeleteAccount(t *testing.T) {
	sqlStore := SetupTestStore(t)
	account1 := createRandomAccount(t)

	err := sqlStore.Queries.DeleteAccount(context.Background(), account1.AccountID)
	require.NoError(t, err)

	account2, err := sqlStore.Queries.GetAccount(context.Background(), account1.AccountID)
	require.NoError(t, err)
	require.False(t, account2.IsActive)

	defer CleanupDB(t)
}

// - TestHardDeleteAccount: Confirms permanent account delet
func TestHardDeleteAccount(t *testing.T) {
	sqlStore := SetupTestStore(t)
	account1 := createRandomAccount(t)

	err := sqlStore.Queries.HardDeleteAccount(context.Background(), account1.AccountID)
	require.NoError(t, err)

	accounts, err := sqlStore.Queries.ListAccountsByUser(context.Background(), account1.UserID)
	require.NoError(t, err)

	for _, account := range accounts {
		require.NotEqual(t, account1.AccountID, account.AccountID)
	}

	defer CleanupDB(t)
}
