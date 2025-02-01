package db

import (
	"context"
	"database/sql"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgtype"
	"github.com/riad/banksystemendtoend/util/common"
	"github.com/stretchr/testify/require"
)

// TestCreateTransaction verifies that a new transaction can be created successfully
func createRandomTransaction(t *testing.T) Transaction {
	sqlStore := SetupTestStore(t)

	var transaction Transaction
	var err error

	amount, err := common.RandomNumeric()
	require.NoError(t, err)

	exchangeRate, err := common.RandomNumeric()
	require.NoError(t, err)

	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		arg := CreateTransactionParams{
			FromAccountID: sql.NullInt32{Int32: createRandomAccount(t).AccountID, Valid: true},
			ToAccountID:   sql.NullInt32{Int32: createRandomAccount(t).AccountID, Valid: true},
			TypeCode:      createRandomTransactionType(t).TypeCode,
			Amount:        amount,
			CurrencyCode:  createRandomCurrency(t).CurrencyCode,
			ExchangeRate:  exchangeRate,
			StatusCode: func() string {
				status, _ := createRandomTransactionStatus(t)
				return status.StatusCode
			}(),
			Description:     sql.NullString{String: common.RandomString(20), Valid: true},
			ReferenceNumber: sql.NullString{String: common.RandomString(10), Valid: true},
			TransactionDate: time.Now(),
		}

		transaction, err = sqlStore.CreateTransaction(context.Background(), arg)
		if err == nil {
			break
		}
		if strings.Contains(err.Error(), "SQLSTATE 23505") {
			if i == maxRetries-1 {
				t.Fatalf("failed to create unique transaction after %d attempts: %v", maxRetries, err)
			}
			continue
		}
		require.NoError(t, err)
	}

	require.NotEmpty(t, transaction)
	return transaction
}

// TestCreateTransaction verifies that a new transaction can be created successfully
func TestCreateTransaction(t *testing.T) {
	createRandomTransaction(t)
	defer CleanupDB(t)
}

// TestGetTransaction verifies retrieval of a specific transaction by its ID
func TestGetTransaction(t *testing.T) {
	sqlStore := SetupTestStore(t)

	tx1 := createRandomTransaction(t)
	tx2, err := sqlStore.GetTransaction(context.Background(), tx1.TransactionID)
	require.NoError(t, err)
	require.NotEmpty(t, tx2)
	require.Equal(t, tx1.TransactionID, tx2.TransactionID)

	defer CleanupDB(t)
}

// TestListTransactionsByAccount ensures correct pagination and filtering of transactions by account ID
func TestListTransactionsByAccount(t *testing.T) {
	sqlStore := SetupTestStore(t)

	amount, err := common.RandomNumeric()
	require.NoError(t, err)

	exchangeRate, err := common.RandomNumeric()
	require.NoError(t, err)

	account := createRandomAccount(t)
	for i := 0; i < 5; i++ {
		arg := CreateTransactionParams{
			FromAccountID: sql.NullInt32{Int32: account.AccountID, Valid: true},
			ToAccountID:   sql.NullInt32{Int32: createRandomAccount(t).AccountID, Valid: true},
			TypeCode:      createRandomTransactionType(t).TypeCode,
			Amount:        amount,
			CurrencyCode:  createRandomCurrency(t).CurrencyCode,
			ExchangeRate:  exchangeRate,
			StatusCode: func() string {
				status, _ := createRandomTransactionStatus(t)
				return status.StatusCode
			}(),
			Description:     sql.NullString{String: "Test transaction", Valid: true},
			ReferenceNumber: sql.NullString{String: common.RandomString(10), Valid: true},
			TransactionDate: time.Now(),
		}
		_, err := sqlStore.CreateTransaction(context.Background(), arg)
		require.NoError(t, err)
	}

	arg := ListTransactionsByAccountParams{
		FromAccountID: sql.NullInt32{Int32: account.AccountID, Valid: true},
		Limit:         5,
		Offset:        2,
	}
	transactions, err := sqlStore.ListTransactionsByAccount(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transactions, 3)

	defer CleanupDB(t)
}

// TestUpdateTransactionStatus verifies that a transaction's status can be updated correctly
func TestUpdateTransactionStatus(t *testing.T) {
	sqlStore := SetupTestStore(t)

	tx := createRandomTransaction(t)
	newStatus, _ := createRandomTransactionStatus(t)

	arg := UpdateTransactionStatusParams{
		TransactionID: tx.TransactionID,
		StatusCode:    newStatus.StatusCode,
	}

	updatedTx, err := sqlStore.UpdateTransactionStatus(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, newStatus.StatusCode, updatedTx.StatusCode)
	require.Equal(t, tx.TransactionID, updatedTx.TransactionID)

	defer CleanupDB(t)
}

// TestGetTransactionsByDateRange checks filtering of transactions within a specified date range
func TestGetTransactionsByDateRange(t *testing.T) {
	sqlStore := SetupTestStore(t)

	for i := 0; i < 3; i++ {
		createRandomTransaction(t)
	}

	arg := GetTransactionsByDateRangeParams{
		StartDate: time.Now().Add(-24 * time.Hour),
		EndDate:   time.Now().Add(24 * time.Hour),
	}

	transactions, er := sqlStore.GetTransactionsByDateRange(context.Background(), arg)
	require.NoError(t, er)
	require.NotEmpty(t, transactions)
	require.Len(t, transactions, 3)

	for _, transaction := range transactions {
		require.True(t, transaction.TransactionDate.After(arg.StartDate) && transaction.TransactionDate.Before(arg.EndDate))
	}

	defer CleanupDB(t)
}

// TestGetTransactionBalance ensures correct calculation of transaction balance for an account
func TestGetTransactionBalance(t *testing.T) {
	sqlStore := SetupTestStore(t)

	account := createRandomAccount(t)

	var expectedBalance pgtype.Numeric
	expectedBalance.Set("0.00")

	amount, err := common.RandomNumeric()
	require.NoError(t, err)

	exchangeRate, err := common.RandomNumeric()
	require.NoError(t, err)

	for i := 0; i < 5; i++ {
		arg := CreateTransactionParams{
			FromAccountID: sql.NullInt32{Int32: account.AccountID, Valid: true},
			ToAccountID:   sql.NullInt32{Int32: createRandomAccount(t).AccountID, Valid: true},
			TypeCode:      createRandomTransactionType(t).TypeCode,
			Amount:        amount,
			CurrencyCode:  createRandomCurrency(t).CurrencyCode,
			ExchangeRate:  exchangeRate,
			StatusCode: func() string {
				status, _ := createRandomTransactionStatus(t)
				return status.StatusCode
			}(),
			Description:     sql.NullString{String: "Test transaction", Valid: true},
			ReferenceNumber: sql.NullString{String: common.RandomString(10), Valid: true},
			TransactionDate: time.Now(),
		}
		_, err := sqlStore.CreateTransaction(context.Background(), arg)
		require.NoError(t, err)
		expectedBalance.Int.Add(expectedBalance.Int, amount.Int)
	}

	balance, err := sqlStore.GetTransactionBalance(context.Background(), sql.NullInt32{Int32: account.AccountID, Valid: true})
	require.NoError(t, err)
	require.Equal(t, expectedBalance, balance)

	defer CleanupDB(t)
}

// TestGetTransactionStatement verifies generation of transaction statements for a given date range
func TestGetTransactionStatement(t *testing.T) {
	sqlStore := SetupTestStore(t)

	account := createRandomAccount(t)

	amount, err := common.RandomNumeric()
	require.NoError(t, err)

	exchangeRate, err := common.RandomNumeric()
	require.NoError(t, err)

	for i := 0; i < 5; i++ {
		arg := CreateTransactionParams{
			FromAccountID: sql.NullInt32{Int32: account.AccountID, Valid: true},
			ToAccountID:   sql.NullInt32{Int32: createRandomAccount(t).AccountID, Valid: true},
			TypeCode:      createRandomTransactionType(t).TypeCode,
			Amount:        amount,
			CurrencyCode:  createRandomCurrency(t).CurrencyCode,
			ExchangeRate:  exchangeRate,
			StatusCode: func() string {
				status, _ := createRandomTransactionStatus(t)
				return status.StatusCode
			}(),
			Description:     sql.NullString{String: "Test transaction", Valid: true},
			ReferenceNumber: sql.NullString{String: common.RandomString(10), Valid: true},
			TransactionDate: time.Now(),
		}
		_, err := sqlStore.CreateTransaction(context.Background(), arg)
		require.NoError(t, err)
	}

	arg := GetTransactionStatementParams{
		FromAccountID:     sql.NullInt32{Int32: account.AccountID, Valid: true},
		TransactionDate:   time.Now().Add(-24 * time.Hour),
		TransactionDate_2: time.Now().Add(24 * time.Hour),
	}
	require.NotEmpty(t, arg)

	statement, err := sqlStore.GetTransactionStatement(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, statement)
	require.Len(t, statement, 5)

	for _, transaction := range statement {
		require.NotEmpty(t, transaction)
		require.True(t, transaction.TransactionDate.After(arg.TransactionDate) && transaction.TransactionDate.Before(arg.TransactionDate_2))
	}

	defer CleanupDB(t)
}

// TestGetTransactionsByStatus checks filtering of transactions by their status code
func TestGetTransactionsByStatus(t *testing.T) {
	sqlStore := SetupTestStore(t)

	status, err := createRandomTransactionStatus(t)
	require.NoError(t, err)

	amount, err := common.RandomNumeric()
	require.NoError(t, err)

	exchangeRate, err := common.RandomNumeric()
	require.NoError(t, err)

	for i := 0; i < 5; i++ {
		arg := CreateTransactionParams{
			FromAccountID:   sql.NullInt32{Int32: createRandomAccount(t).AccountID, Valid: true},
			ToAccountID:     sql.NullInt32{Int32: createRandomAccount(t).AccountID, Valid: true},
			TypeCode:        createRandomTransactionType(t).TypeCode,
			Amount:          amount,
			CurrencyCode:    createRandomCurrency(t).CurrencyCode,
			ExchangeRate:    exchangeRate,
			StatusCode:      status.StatusCode,
			Description:     sql.NullString{String: "Test transaction", Valid: true},
			ReferenceNumber: sql.NullString{String: common.RandomString(10), Valid: true},
			TransactionDate: time.Now(),
		}
		_, err := sqlStore.CreateTransaction(context.Background(), arg)
		require.NoError(t, err)
	}

	transactions, err := sqlStore.GetTransactionsByStatus(context.Background(), status.StatusCode)
	require.NoError(t, err)
	require.NotEmpty(t, transactions)

	for _, transaction := range transactions {
		require.Equal(t, status.StatusCode, transaction.StatusCode)
	}

	defer CleanupDB(t)
}

// TestGetTransactionByReference verifies retrieval of a transaction using its reference number
func TestGetTransactionByReference(t *testing.T) {
	store := SetupTestStore(t)

	tx1 := createRandomTransaction(t)
	tx2, err := store.GetTransactionByReference(context.Background(), tx1.ReferenceNumber)

	require.NoError(t, err)
	require.Equal(t, tx1.TransactionID, tx2.TransactionID)

	defer CleanupDB(t)
}

// TestDeleteTransaction ensures a single transaction can be deleted successfully
func TestDeleteTransaction(t *testing.T) {
	sqlStore := SetupTestStore(t)

	tx := createRandomTransaction(t)

	err := sqlStore.DeleteTransaction(context.Background(), tx.TransactionNumber)
	require.NoError(t, err)

	_, err = sqlStore.GetTransaction(context.Background(), tx.TransactionID)
	require.Error(t, err)
	require.EqualError(t, err, "no rows in result set")

	defer CleanupDB(t)
}

// TestDeleteAccountTransactions verifies deletion of all transactions associated with an account
func TestDeleteAccountTransactions(t *testing.T) {
	sqlStore := SetupTestStore(t)

	var err error

	status, err := createRandomTransactionStatus(t)
	require.NoError(t, err)

	amount, err := common.RandomNumeric()
	require.NoError(t, err)

	exchangeRate, err := common.RandomNumeric()
	require.NoError(t, err)

	account := createRandomAccount(t)

	for i := 0; i < 5; i++ {
		arg := CreateTransactionParams{
			FromAccountID:   sql.NullInt32{Int32: account.AccountID, Valid: true},
			ToAccountID:     sql.NullInt32{Int32: createRandomAccount(t).AccountID, Valid: true},
			TypeCode:        createRandomTransactionType(t).TypeCode,
			Amount:          amount,
			CurrencyCode:    createRandomCurrency(t).CurrencyCode,
			ExchangeRate:    exchangeRate,
			StatusCode:      status.StatusCode,
			Description:     sql.NullString{String: "Test transaction", Valid: true},
			ReferenceNumber: sql.NullString{String: common.RandomString(10), Valid: true},
			TransactionDate: time.Now(),
		}
		_, err := sqlStore.CreateTransaction(context.Background(), arg)
		require.NoError(t, err)
	}

	err = sqlStore.DeleteAccountTransactions(context.Background(), sql.NullInt32{Int32: account.AccountID, Valid: true})
	require.NoError(t, err)

	transactions, err := sqlStore.ListAccountTransactions(context.Background(), ListAccountTransactionsParams{
		FromAccountID: sql.NullInt32{Int32: account.AccountID, Valid: true},
		Limit:         10,
		Offset:        0,
	})
	require.NoError(t, err)
	require.Empty(t, transactions)

	defer CleanupDB(t)
}

// TestDeleteTransactionsByDateRange ensures transactions within a date range are deleted correctly
func TestDeleteTransactionsByDateRange(t *testing.T) {
	// Create a new SQLStore instance
	sqlStore := SetupTestStore(t)

	startDate := time.Now().Add(-24 * time.Hour)
	endDate := time.Now().Add(24 * time.Hour)

	for i := 0; i < 10; i++ {
		createRandomTransaction(t)
	}

	err := sqlStore.DeleteTransactionsByDateRange(context.Background(), DeleteTransactionsByDateRangeParams{
		StartDate: startDate,
		EndDate:   endDate,
	})
	require.NoError(t, err)

	transactions, err := sqlStore.GetTransactionsByDateRange(context.Background(), GetTransactionsByDateRangeParams{
		StartDate: startDate,
		EndDate:   endDate,
	})

	require.NoError(t, err)
	require.Empty(t, transactions)

	defer CleanupDB(t)
}

// TestDeleteTransactionsByStatus verifies deletion of all transactions with a specific status
func TestDeleteTransactionsByStatus(t *testing.T) {
	sqlStore := SetupTestStore(t)

	var err error

	status, err := createRandomTransactionStatus(t)
	require.NoError(t, err)

	amount, err := common.RandomNumeric()
	require.NoError(t, err)

	exchangeRate, err := common.RandomNumeric()
	require.NoError(t, err)

	for i := 0; i < 5; i++ {
		arg := CreateTransactionParams{
			FromAccountID:   sql.NullInt32{Int32: createRandomAccount(t).AccountID, Valid: true},
			ToAccountID:     sql.NullInt32{Int32: createRandomAccount(t).AccountID, Valid: true},
			TypeCode:        createRandomTransactionType(t).TypeCode,
			Amount:          amount,
			CurrencyCode:    createRandomCurrency(t).CurrencyCode,
			ExchangeRate:    exchangeRate,
			StatusCode:      status.StatusCode,
			Description:     sql.NullString{String: "Test transaction", Valid: true},
			ReferenceNumber: sql.NullString{String: common.RandomString(10), Valid: true},
			TransactionDate: time.Now(),
		}
		_, err := sqlStore.CreateTransaction(context.Background(), arg)
		require.NoError(t, err)
	}

	err = sqlStore.DeleteTransactionsByStatus(context.Background(), status.StatusCode)
	require.NoError(t, err)

	transactions, err := sqlStore.GetTransactionsByStatus(context.Background(), status.StatusCode)
	require.NoError(t, err)
	require.Empty(t, transactions)

	defer CleanupDB(t)
}
