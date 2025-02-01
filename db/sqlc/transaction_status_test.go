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
func createRandomTransactionStatus(t *testing.T) (TransactionStatus, error) {
	sqlStore := SetupTestStore(t)

	arg := CreateTransactionStatusParams{
		StatusCode:  common.RandomString(15),
		Description: common.RandomString(30),
	}

	status, err := sqlStore.Queries.CreateTransactionStatus(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, status)
	require.Equal(t, arg.StatusCode, status.StatusCode)
	require.Equal(t, arg.Description, status.Description)
	require.True(t, status.IsActive)
	return status, nil
}

// TestCreateTransactionStatus verifies creation of new transaction status records.
// Cleans database before test execution.
func TestCreateTransactionStatus(t *testing.T) {
	createRandomTransactionStatus(t)
	defer CleanupDB(t)
}

// TestGetTransactionStatus verifies retrieval of transaction status records.
// Creates a test record and validates exact field matching on retrieval.
func TestGetTransactionStatus(t *testing.T) {
	sqlStore := SetupTestStore(t)

	status1, err := createRandomTransactionStatus(t)
	require.NoError(t, err)

	status2, err := sqlStore.Queries.GetTransactionStatus(context.Background(), status1.StatusCode)
	require.NoError(t, err)
	require.NotEmpty(t, status2)
	require.Equal(t, status1.StatusCode, status2.StatusCode)
	require.Equal(t, status1.Description, status2.Description)
	require.Equal(t, status1.IsActive, status2.IsActive)

	defer CleanupDB(t)

}

// TestModifyTransactionStatus verifies status record updates.
// Tests modification of description and active state, ensuring proper field updates.
func TestModifyTransactionStatus(t *testing.T) {
	sqlStore := SetupTestStore(t)

	status1, err := createRandomTransactionStatus(t)
	require.NoError(t, err)

	arg := ModifyTransactionStatusParams{
		StatusCode:  status1.StatusCode,
		Description: sql.NullString{String: common.RandomString(20), Valid: true},
		IsActive:    sql.NullBool{Bool: false, Valid: true},
	}

	updatedStatus, err := sqlStore.Queries.ModifyTransactionStatus(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedStatus)
	require.Equal(t, status1.StatusCode, updatedStatus.StatusCode)
	require.NotEqual(t, status1.Description, updatedStatus.Description)
	require.NotEqual(t, status1.IsActive, updatedStatus.IsActive)

	defer CleanupDB(t)

}

// TestDeleteTransactionStatus verifies soft deletion functionality.
// Confirms status is marked inactive but still retrievable.
func TestDeleteTransactionStatus(t *testing.T) {
	sqlStore := SetupTestStore(t)
	require.NotEmpty(t, sqlStore)

	status1, err := createRandomTransactionStatus(t)
	require.NoError(t, err)

	err = sqlStore.Queries.DeleteTransactionStatus(context.Background(), status1.StatusCode)
	require.NoError(t, err)

	status2, err := sqlStore.Queries.GetTransactionStatus(context.Background(), status1.StatusCode)
	require.NoError(t, err)
	require.NotEmpty(t, status2)
	require.False(t, status2.IsActive)

	defer CleanupDB(t)

}

// TestHardDeleteTransactionStatus verifies permanent deletion.
// Ensures status record is completely removed from database.
func TestHardDeleteTransactionStatus(t *testing.T) {
	sqlStore := SetupTestStore(t)

	status1, err := createRandomTransactionStatus(t)
	require.NoError(t, err)

	err = sqlStore.Queries.HardDeleteTransactionStatus(context.Background(), status1.StatusCode)
	require.NoError(t, err)

	status2, err := sqlStore.Queries.GetTransactionStatus(context.Background(), status1.StatusCode)
	require.Error(t, err)
	require.EqualError(t, err, "no rows in result set")
	require.Empty(t, status2)

	defer CleanupDB(t)

}

// TestGetActiveTransactionStatus verifies bulk active status retrieval.
// Creates multiple test statuses and validates retrieval by status code array.
func TestGetActiveTransactionStatus(t *testing.T) {
	sqlStore := SetupTestStore(t)

	var statusCodes []string
	for i := 0; i < 3; i++ {
		status, err := createRandomTransactionStatus(t)
		require.NoError(t, err)
		statusCodes = append(statusCodes, status.StatusCode)
		require.NotEmpty(t, statusCodes)
	}

	activeStatuses, err := sqlStore.Queries.GetActiveTransactionStatus(context.Background(), statusCodes)
	require.NoError(t, err)
	require.Len(t, activeStatuses, len(statusCodes))

	for _, status := range activeStatuses {
		require.NotEmpty(t, status)
		require.True(t, status.IsActive)
		require.Contains(t, statusCodes, status.StatusCode)
	}
	defer CleanupDB(t)
}
