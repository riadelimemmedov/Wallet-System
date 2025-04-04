package db

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	db "github.com/riad/banksystemendtoend/db/sqlc"
	"github.com/riad/banksystemendtoend/util/common"
	"github.com/stretchr/testify/require"
)

// createRandomCurrency creates a new random currency record for testing purposes.
// It handles the entire creation process including parameter generation, database insertion,
// and validation of the created currency.
func createRandomCurrency(t *testing.T) db.AccountCurrency {
	t.Helper()
	arg := generateCurrencyParams()

	accountCurrency, err := createCurrencyWithRetry(t, arg)
	require.NoError(t, err)

	validateCreatedCurrency(t, accountCurrency, arg)
	return accountCurrency
}

// generateCurrencyParams creates a new set of random currency parameters for testing.
// It generates random values for currency code, exchange rate, name, and symbol.
// Panics if the exchange rate numeric conversion fails.
func generateCurrencyParams() db.CreateCurrencyParams {
	currency_code := common.RandomString(3)
	formattedValue := strconv.FormatFloat(common.RandomFloat(1.1, 99.99), 'f', -1, 64)
	exchangeRate, err := common.SetNumeric(formattedValue)
	if err != nil {
		panic(fmt.Sprintf("failed to set numeric value: %v", err))
	}

	return db.CreateCurrencyParams{
		CurrencyCode: currency_code,
		CurrencyName: common.RandomString(3),
		Symbol: sql.NullString{
			String: common.RandomString(3),
			Valid:  true,
		},
		ExchangeRate: exchangeRate,
	}
}

// createCurrencyWithRetry attempts to create a currency record with retry logic for handling
// duplicate key errors. It will retry up to maxRetries times before giving up.
// Returns the created currency and any error encountered.
func createCurrencyWithRetry(t *testing.T, arg db.CreateCurrencyParams) (db.AccountCurrency, error) {
	sqlStore := SetupTestStore(t)

	const maxRetries = 10
	var accountCurrency db.AccountCurrency
	var err error

	for attempt := 0; attempt < maxRetries; attempt++ {
		accountCurrency, err = sqlStore.Queries.CreateCurrency(context.Background(), arg)

		if err == nil {
			return accountCurrency, nil
		}

		if !isDuplicateKeyError(err) {
			return db.AccountCurrency{}, err
		}

		if attempt == maxRetries-1 {
			return db.AccountCurrency{}, fmt.Errorf("failed to create unique currency after %d attempts: %w", maxRetries, err)
		}
	}

	return accountCurrency, err
}

// isDuplicateKeyError checks if the given error is a PostgreSQL duplicate key error.
// Returns true if the error contains the specific PostgreSQL error code for duplicate keys.
func isDuplicateKeyError(err error) bool {
	return strings.Contains(err.Error(), "SQLSTATE 23505")
}

// validateCreatedCurrency validates that a created currency record matches its input parameters
// and meets all required conditions (non-empty, active, matching fields).
func validateCreatedCurrency(t *testing.T, currency db.AccountCurrency, arg db.CreateCurrencyParams) {
	require.NotEmpty(t, currency)
	require.True(t, currency.IsActive)
	require.Equal(t, arg.CurrencyCode, currency.CurrencyCode)
	require.Equal(t, arg.CurrencyName, currency.CurrencyName)
	require.Equal(t, arg.Symbol, currency.Symbol)
}

// TestCreateCurrency verifies that a new currency can be created successfully.
func TestCreateCurrency(t *testing.T) {
	createRandomCurrency(t)
	defer CleanupDB(t)
}

// TestGetCurrency verifies that a currency can be retrieved correctly after creation
// and that all fields match the original values.
func TestGetCurrency(t *testing.T) {
	sqlStore := SetupTestStore(t)

	currency1 := createRandomCurrency(t)

	currency2, err := sqlStore.Queries.GetCurrency(ctx, currency1.CurrencyCode)
	require.NoError(t, err)
	require.NotEmpty(t, currency2)

	require.Equal(t, currency1.CurrencyCode, currency2.CurrencyCode)
	require.Equal(t, currency1.CurrencyName, currency2.CurrencyName)
	require.Equal(t, currency1.Symbol, currency2.Symbol)
	require.Equal(t, currency1.ExchangeRate, currency2.ExchangeRate)
	require.Equal(t, currency1.IsActive, currency2.IsActive)
	require.WithinDuration(t, currency1.CreatedAt, currency2.CreatedAt, time.Second)

	defer CleanupDB(t)
}

// TestUpdateExchangeRate verifies that a currency's exchange rate can be updated
// successfully and that the update is reflected in the database.
func TestUpdateExchangeRate(t *testing.T) {
	sqlStore := SetupTestStore(t)

	currency1 := createRandomCurrency(t)

	formattedValue := strconv.FormatFloat(common.RandomFloat(1.1, 99.99), 'f', -1, 64)
	exchangeRate, _ := common.SetNumeric(formattedValue)

	arg := db.UpdateExchangeRateParams{
		CurrencyCode: currency1.CurrencyCode,
		ExchangeRate: exchangeRate,
	}

	currency2, err := sqlStore.Queries.UpdateExchangeRate(ctx, arg)
	require.NoError(t, err)
	require.NotEmpty(t, currency2)

	require.Equal(t, currency1.CurrencyCode, currency2.CurrencyCode)

	defer CleanupDB(t)
}

// TestDeleteCurrency verifies that a currency can be soft deleted (marked as inactive)
// and that the deletion is reflected in the database.
func TestDeleteCurrency(t *testing.T) {
	sqlStore := SetupTestStore(t)

	currency1 := createRandomCurrency(t)

	err := sqlStore.Queries.DeleteCurrency(ctx, currency1.CurrencyCode)
	require.NoError(t, err)

	currency2, err := sqlStore.Queries.GetCurrency(ctx, currency1.CurrencyCode)
	require.NoError(t, err)
	require.False(t, currency2.IsActive)

	defer CleanupDB(t)
}
