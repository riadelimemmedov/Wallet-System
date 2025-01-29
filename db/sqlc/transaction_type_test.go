package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/riad/banksystemendtoend/util/common"
	"github.com/stretchr/testify/require"
)

// createRandomTransactionType creates a new transaction type with random values for testing.
// It generates random type code and description, creates the record, and validates the creation.
// Returns the created TransactionType.
func createRandomTransactionType(t *testing.T) TransactionType {
	sqlStore := SetupTestStore(t)
	require.NotEmpty(t, sqlStore)

	typeCode := common.RandomString(5)
	description := common.RandomString(30)

	arg := CreateTransactionTypeParams{
		TypeCode:    typeCode,
		Description: description,
	}

	transType, err := sqlStore.Queries.CreateTransactionType(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transType)
	require.Equal(t, typeCode, transType.TypeCode)
	require.Equal(t, description, transType.Description)
	require.True(t, transType.IsActive)

	return transType
}

// TestCreateTransactionType verifies the creation of a new transaction type.
// It cleans the database before testing and validates the creation process.
func TestCreateTransactionType(t *testing.T) {
	createRandomTransactionType(t)
	defer CleanupDB(t)
}

// TestGetTransactionType verifies retrieval of an existing transaction type.
// It creates a test record, retrieves it, and validates all fields match.
func TestGetTransactionType(t *testing.T) {
	sqlStore := SetupTestStore(t)
	require.NotEmpty(t, sqlStore)

	transType1 := createRandomTransactionType(t)
	transType2, err := sqlStore.Queries.GetTransactionType(context.Background(), transType1.TypeCode)

	require.NoError(t, err)
	require.NotEmpty(t, transType2)
	require.Equal(t, transType1.TypeCode, transType2.TypeCode)
	require.Equal(t, transType1.Description, transType2.Description)
	require.Equal(t, transType1.IsActive, transType2.IsActive)
}

// TestUpdateTransactionType verifies updating transaction type fields.
// It creates a test record, updates its description and status, and validates changes.
func TestUpdateTransactionType(t *testing.T) {
	sqlStore := SetupTestStore(t)
	require.NotEmpty(t, sqlStore)

	transType1 := createRandomTransactionType(t)

	arg := UpdateTransactionTypeParams{
		TypeCode:    transType1.TypeCode,
		Description: sql.NullString{String: common.RandomString(20), Valid: true},
		IsActive:    sql.NullBool{Bool: false, Valid: true},
	}
	updatedTransType, err := sqlStore.Queries.UpdateTransactionType(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedTransType)
	require.Equal(t, transType1.TypeCode, updatedTransType.TypeCode)
	require.NotEqual(t, transType1.Description, updatedTransType.Description)
	require.NotEqual(t, transType1.IsActive, updatedTransType.IsActive)

	defer CleanupDB(t)
}

// TestDeleteTransactionType verifies soft deletion of a transaction type.
// It creates a test record, marks it as inactive, and validates the status change.
func TestDeleteTransactionType(t *testing.T) {
	sqlStore := SetupTestStore(t)
	require.NotEmpty(t, sqlStore)

	transType1 := createRandomTransactionType(t)
	err := sqlStore.Queries.DeleteTransactionType(context.Background(), transType1.TypeCode)
	require.NoError(t, err)

	transType2, err := sqlStore.Queries.GetTransactionType(context.Background(), transType1.TypeCode)
	require.NoError(t, err)
	require.NotEmpty(t, transType2)
	require.NotEqual(t, transType1.IsActive, transType2.IsActive)

	defer CleanupDB(t)
}

// TestHardDeleteTransactionType verifies permanent deletion of a transaction type.
// It creates a test record, permanently deletes it, and validates it no longer exists.
func TestHardDeleteTransactionType(t *testing.T) {
	sqlStore := SetupTestStore(t)
	require.NotEmpty(t, sqlStore)

	transType1 := createRandomTransactionType(t)
	err := sqlStore.Queries.HardDeleteTransactionType(context.Background(), transType1.TypeCode)
	require.NoError(t, err)

	transType2, err := sqlStore.Queries.GetTransactionType(context.Background(), transType1.TypeCode)
	require.Error(t, err)
	require.EqualError(t, err, "no rows in result set")
	require.Empty(t, transType2)

	defer CleanupDB(t)
}
