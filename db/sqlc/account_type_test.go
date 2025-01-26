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
//
// Parameters:
//   - t: Testing object for managing test state and assertions
//
// Returns:
//   - AccountType: A newly created account type with random data
//
// The function will fail the test if:
//   - It cannot create a unique account type after maximum retries
//   - Any unexpected database error occurs
//   - The created account type is empty
//   - The created account type is not marked as active
func createRandomAccountType(t *testing.T) AccountType {
	var accountType AccountType
	var err error

	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		arg := CreateAccountTypeParams{
			AccountType: common.RandomAccountType(),
			Description: common.RandomString(20),
		}
		accountType, err = testDB.Queries.CreateAccountType(context.Background(), arg)
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
//
// This test ensures that:
//   - The database can be cleaned successfully
//   - A new account type can be created with random valid data
func TestCreateAccountType(t *testing.T) {
	CleanupDB(t, testDB)
	createRandomAccountType(t)
}

// TestListAccountTypes verifies the listing of account types.
// It creates multiple account types and ensures they can be retrieved correctly.
//
// The test:
//   - Cleans the database
//   - Creates 3 random account types
//   - Retrieves all account types
//   - Verifies that all created account types are present in the list
//   - Ensures each retrieved account type is valid and active
func TestListAccountTypes(t *testing.T) {
	CleanupDB(t, testDB)
	createdTypes := make(map[string]bool)
	for i := 0; i < 3; i++ {
		accountType := createRandomAccountType(t)
		createdTypes[accountType.AccountType] = true
	}

	accountTypes, err := testDB.Queries.ListAccountTypes(context.Background())
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
}

// TestUpdateAccountType verifies the updating of an account type.
// It creates an account type and then updates its type field with a new random value.
//
// The test ensures:
//   - Initial account type can be created
//   - The type field can be updated to a new unique value
//   - Other fields (description, isActive) remain unchanged
//   - The update operation returns the correct updated account type
func TestUpdateAccountType(t *testing.T) {
	CleanupDB(t, testDB)
	accountType1 := createRandomAccountType(t)

	newType := common.RandomAccountType()
	for newType == accountType1.AccountType {
		newType = common.RandomAccountType()
	}

	arg := UpdateAccountTypeParams{
		AccountType:   newType,
		AccountType_2: accountType1.AccountType,
	}
	accountType2, err := testDB.Queries.UpdateAccountType(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, accountType2)

	require.Equal(t, arg.AccountType, accountType2.AccountType)
	require.Equal(t, accountType1.Description, accountType2.Description)
	require.Equal(t, accountType1.IsActive, accountType2.IsActive)
}

// TestDeleteAccountType verifies the soft deletion of an account type.
// It creates an account type and then marks it as inactive (soft delete).
//
// The test verifies:
//   - Initial account type can be created
//   - The account type can be soft deleted
//   - The deleted account type appears as inactive in the list
func TestDeleteAccountType(t *testing.T) {
	CleanupDB(t, testDB)
	accountType1 := createRandomAccountType(t)

	err := testDB.Queries.DeleteAccountType(context.Background(), accountType1.AccountType)
	require.NoError(t, err)

	accountTypes, err := testDB.Queries.ListAccountTypes(context.Background())
	require.NoError(t, err)

	for _, accountType := range accountTypes {
		if accountType.AccountType == accountType1.AccountType {
			require.False(t, accountType.IsActive)
		}
	}
}

// TestHardDeleteAccountType verifies the permanent deletion of an account type.
// It creates an account type and then completely removes it from the database.
//
// The test ensures:
//   - Initial account type can be created
//   - The account type can be permanently deleted
//   - The deleted account type does not appear in the list after deletion
func TestHardDeleteAccountType(t *testing.T) {
	CleanupDB(t, testDB)
	accountType1 := createRandomAccountType(t)

	err := testDB.Queries.HardDeleteAccountType(context.Background(), accountType1.AccountType)
	require.NoError(t, err)

	accountTypes, err := testDB.Queries.ListAccountTypes(context.Background())
	require.NoError(t, err)

	for _, accountType := range accountTypes {
		if accountType.AccountType == accountType1.AccountType {
			require.NotEqual(t, accountType1.AccountType, accountType.AccountType)
		}
	}
}
