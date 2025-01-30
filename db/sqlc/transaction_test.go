package db

// import (
// 	"context"
// 	"database/sql"
// 	"strings"
// 	"testing"
// 	"time"

// 	"github.com/riad/banksystemendtoend/util/common"
// 	"github.com/stretchr/testify/require"
// )

// func createRandomTransaction(t *testing.T) Transaction {
// 	sqlStore := SetupTestStore(t)
// 	require.NotEmpty(t, sqlStore)

// 	var transaction Transaction
// 	var err error

// 	amount, err := common.RandomNumeric()
// 	require.NoError(t, err)

// 	exchangeRate, err := common.RandomNumeric()
// 	require.NoError(t, err)

// 	maxRetries := 5
// 	for i := 0; i < maxRetries; i++ {
// 		arg := CreateTransactionParams{
// 			FromAccountID: sql.NullInt32{Int32: createRandomAccount(t).AccountID, Valid: true},
// 			ToAccountID:   sql.NullInt32{Int32: createRandomAccount(t).AccountID, Valid: true},
// 			TypeCode:      createRandomTransactionType(t).TypeCode,
// 			Amount:        amount,
// 			CurrencyCode:  createRandomCurrency(t).CurrencyCode,
// 			ExchangeRate:  exchangeRate,
// 			StatusCode: func() string {
// 				status, err := createRandomTransactionStatus(t)
// 				require.NoError(t, err)
// 				return status.StatusCode
// 			}(),
// 			Description:     sql.NullString{String: common.RandomString(20), Valid: true},
// 			ReferenceNumber: sql.NullString{String: common.RandomString(10), Valid: true},
// 			TransactionDate: time.Now(),
// 		}

// 		transaction, err = sqlStore.CreateTransaction(context.Background(), arg)
// 		if err == nil {
// 			break
// 		}
// 		if strings.Contains(err.Error(), "SQLSTATE 23505") {
// 			if i == maxRetries-1 {
// 				t.Fatalf("failed to create unique transaction after %d attempts: %v", maxRetries, err)
// 			}
// 			continue
// 		}
// 		require.NoError(t, err)
// 	}

// 	require.NotEmpty(t, transaction)
// 	return transaction
// }

// func TestCreateTransaction(t *testing.T) {
// 	createRandomTransaction(t)
// 	defer CleanupDB(t)
// }
