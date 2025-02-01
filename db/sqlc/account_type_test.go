package db

import (
	"context"
	"strings"
	"testing"

	"github.com/riad/banksystemendtoend/util/common"
	"github.com/stretchr/testify/require"
)

// createRandomAccountType creates a new random account type for testing purposes.
// It handles potential duplicate key violations by retrying up to 5 times with different random values.
func createRandomAccountType(t *testing.T) AccountType {
	sqlStore := SetupTestStore(t)

	var accountType AccountType
	var err error

	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		arg := CreateAccountTypeParams{
			AccountType: common.RandomString(15),
			Description: common.RandomString(20),
		}
		accountType, err = sqlStore.Queries.CreateAccountType(context.Background(), arg)
		if err == nil {
			break
		}
		if strings.Contains(err.Error(), "SQLSTATE 23505") {
			if i == maxRetries-1 {
				t.Fatalf("failed to create unique account type after %d attempts: %v", maxRetries, err)
			}
			continue
		}
		require.NoError(t, err)
	}
	require.NotEmpty(t, accountType)
	require.True(t, accountType.IsActive)
	return accountType
}

// TestCreateAccountType verifies the creation of a new account type.
// It first cleans the database and then creates a random account type.
func TestCreateAccountType(t *testing.T) {
	createRandomAccountType(t)
	defer CleanupDB(t)
}

// TestListAccountTypes verifies the listing of account types.
// It creates multiple account types and ensures they can be retrieved correctly.
func TestListAccountTypes(t *testing.T) {
	sqlStore := SetupTestStore(t)

	createdTypes := make(map[string]bool)
	for i := 0; i < 3; i++ {
		accountType := createRandomAccountType(t)
		createdTypes[accountType.AccountType] = true
	}

	accountTypes, err := sqlStore.Queries.ListAccountTypes(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, accountTypes)

	foundTypes := 0
	for _, accountType := range accountTypes {
		require.NotEmpty(t, accountType)
		require.True(t, accountType.IsActive)
		if createdTypes[accountType.AccountType] {
			foundTypes++
		}
	}
	require.True(t, foundTypes > 0)
	require.True(t, foundTypes == 3)

	defer CleanupDB(t)
}

// TestUpdateAccountType verifies the updating of an account type.
// It creates an account type and then updates its type field with a new random value.
func TestUpdateAccountType(t *testing.T) {
	sqlStore := SetupTestStore(t)

	accountType1 := createRandomAccountType(t)

	newType := common.RandomString(15)
	for newType == accountType1.AccountType {
		newType = common.RandomString(15)
	}

	arg := UpdateAccountTypeParams{
		AccountType:   newType,
		AccountType_2: accountType1.AccountType,
	}
	accountType2, err := sqlStore.Queries.UpdateAccountType(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, accountType2)

	require.Equal(t, arg.AccountType, accountType2.AccountType)
	require.Equal(t, accountType1.Description, accountType2.Description)
	require.Equal(t, accountType1.IsActive, accountType2.IsActive)

	defer CleanupDB(t)
}

// TestDeleteAccountType verifies the soft deletion of an account type.
// It creates an account type and then marks it as inactive (soft delete).
func TestDeleteAccountType(t *testing.T) {
	sqlStore := SetupTestStore(t)

	accountType1 := createRandomAccountType(t)

	err := sqlStore.Queries.DeleteAccountType(context.Background(), accountType1.AccountType)
	require.NoError(t, err)

	accountTypes, err := sqlStore.Queries.ListAccountTypes(context.Background())
	require.NoError(t, err)

	for _, accountType := range accountTypes {
		if accountType.AccountType == accountType1.AccountType {
			require.False(t, accountType.IsActive)
		}
	}
	defer CleanupDB(t)
}

// TestHardDeleteAccountType verifies the permanent deletion of an account type.
// It creates an account type and then completely removes it from the database.
func TestHardDeleteAccountType(t *testing.T) {
	sqlStore := SetupTestStore(t)

	accountType1 := createRandomAccountType(t)

	err := sqlStore.Queries.HardDeleteAccountType(context.Background(), accountType1.AccountType)
	require.NoError(t, err)

	accountTypes, err := sqlStore.Queries.ListAccountTypes(context.Background())
	require.NoError(t, err)

	for _, accountType := range accountTypes {
		if accountType.AccountType == accountType1.AccountType {
			require.NotEqual(t, accountType1.AccountType, accountType.AccountType)
		}
	}
	defer CleanupDB(t)
}
