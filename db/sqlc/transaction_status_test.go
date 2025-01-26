package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/riad/banksystemendtoend/util/common"
	"github.com/stretchr/testify/require"
)

// createRandomTransactionStatus creates a test transaction status with random values.
// Returns a TransactionStatus object with random status code and description.
// The created status is verified for proper initialization and active state.
func createRandomTransactionStatus(t *testing.T) TransactionStatus {
	arg := CreateTransactionStatusParams{
		StatusCode:  common.RandomString(5),
		Description: common.RandomString(30),
	}

	status, err := testDB.Queries.CreateTransactionStatus(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, status)
	require.Equal(t, arg.StatusCode, status.StatusCode)
	require.Equal(t, arg.Description, status.Description)
	require.True(t, status.IsActive)
	return status
}

// TestCreateTransactionStatus verifies creation of new transaction status records.
// Cleans database before test execution.
func TestCreateTransactionStatus(t *testing.T) {
	CleanupDB(t, testDB)
	createRandomTransactionStatus(t)
}

// TestGetTransactionStatus verifies retrieval of transaction status records.
// Creates a test record and validates exact field matching on retrieval.
func TestGetTransactionStatus(t *testing.T) {
	CleanupDB(t, testDB)

	status1 := createRandomTransactionStatus(t)
	status2, err := testDB.Queries.GetTransactionStatus(context.Background(), status1.StatusCode)

	require.NoError(t, err)
	require.NotEmpty(t, status2)
	require.Equal(t, status1.StatusCode, status2.StatusCode)
	require.Equal(t, status1.Description, status2.Description)
	require.Equal(t, status1.IsActive, status2.IsActive)
}

// TestModifyTransactionStatus verifies status record updates.
// Tests modification of description and active state, ensuring proper field updates.
func TestModifyTransactionStatus(t *testing.T) {
	CleanupDB(t, testDB)

	status1 := createRandomTransactionStatus(t)

	arg := ModifyTransactionStatusParams{
		StatusCode:  status1.StatusCode,
		Description: sql.NullString{String: common.RandomString(20), Valid: true},
		IsActive:    sql.NullBool{Bool: false, Valid: true},
	}

	updatedStatus, err := testDB.Queries.ModifyTransactionStatus(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedStatus)
	require.Equal(t, status1.StatusCode, updatedStatus.StatusCode)
	require.NotEqual(t, status1.Description, updatedStatus.Description)
	require.NotEqual(t, status1.IsActive, updatedStatus.IsActive)
}

// TestDeleteTransactionStatus verifies soft deletion functionality.
// Confirms status is marked inactive but still retrievable.
func TestDeleteTransactionStatus(t *testing.T) {
	CleanupDB(t, testDB)

	status1 := createRandomTransactionStatus(t)
	err := testDB.Queries.DeleteTransactionStatus(context.Background(), status1.StatusCode)
	require.NoError(t, err)

	status2, err := testDB.Queries.GetTransactionStatus(context.Background(), status1.StatusCode)
	require.NoError(t, err)
	require.NotEmpty(t, status2)
	require.False(t, status2.IsActive)
}

// TestHardDeleteTransactionStatus verifies permanent deletion.
// Ensures status record is completely removed from database.
func TestHardDeleteTransactionStatus(t *testing.T) {
	CleanupDB(t, testDB)

	status1 := createRandomTransactionStatus(t)
	err := testDB.Queries.HardDeleteTransactionStatus(context.Background(), status1.StatusCode)
	require.NoError(t, err)

	status2, err := testDB.Queries.GetTransactionStatus(context.Background(), status1.StatusCode)
	require.Error(t, err)
	require.EqualError(t, err, "no rows in result set")
	require.Empty(t, status2)
}

// TestGetActiveTransactionStatus verifies bulk active status retrieval.
// Creates multiple test statuses and validates retrieval by status code array.
func TestGetActiveTransactionStatus(t *testing.T) {
	CleanupDB(t, testDB)

	var statusCodes []string
	for i := 0; i < 3; i++ {
		status := createRandomTransactionStatus(t)
		statusCodes = append(statusCodes, status.StatusCode)
	}

	activeStatuses, err := testDB.Queries.GetActiveTransactionStatus(context.Background(), statusCodes)
	require.NoError(t, err)
	require.Len(t, activeStatuses, len(statusCodes))

	for _, status := range activeStatuses {
		require.NotEmpty(t, status)
		require.True(t, status.IsActive)
		require.Contains(t, statusCodes, status.StatusCode)
	}
}
