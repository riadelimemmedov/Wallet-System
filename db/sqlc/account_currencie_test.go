package db

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/riad/banksystemendtoend/util"
	"github.com/stretchr/testify/require"
)

// createRandomCurrency creates a new random currency record for testing purposes.
// It handles the entire creation process including parameter generation, database insertion,
// and validation of the created currency.
func createRandomCurrency(t *testing.T) AccountCurrency {
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
func generateCurrencyParams() CreateCurrencyParams {
	currency_code := util.RandomCurrency()
	formattedValue := strconv.FormatFloat(util.RandomFloat(1.1, 99.99), 'f', -1, 64)
	exchangeRate, err := util.SetNumeric(formattedValue)
	if err != nil {
		panic(fmt.Sprintf("failed to set numeric value: %v", err))
	}

	return CreateCurrencyParams{
		CurrencyCode: currency_code,
		CurrencyName: util.RandomCurrencyName(currency_code),
		Symbol: sql.NullString{
			String: util.RandomCurrencySymbol(currency_code),
			Valid:  true,
		},
		ExchangeRate: exchangeRate,
	}
}

// createCurrencyWithRetry attempts to create a currency record with retry logic for handling
// duplicate key errors. It will retry up to maxRetries times before giving up.
// Returns the created currency and any error encountered.
func createCurrencyWithRetry(t *testing.T, arg CreateCurrencyParams) (AccountCurrency, error) {
	const maxRetries = 10
	var accountCurrency AccountCurrency
	var err error

	for attempt := 0; attempt < maxRetries; attempt++ {
		accountCurrency, err = testDB.Queries.CreateCurrency(context.Background(), arg)

		if err == nil {
			return accountCurrency, nil
		}

		if !isDuplicateKeyError(err) {
			return AccountCurrency{}, err
		}

		if attempt == maxRetries-1 {
			return AccountCurrency{}, fmt.Errorf("failed to create unique currency after %d attempts: %w", maxRetries, err)
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
func validateCreatedCurrency(t *testing.T, currency AccountCurrency, arg CreateCurrencyParams) {
	require.NotEmpty(t, currency)
	require.True(t, currency.IsActive)
	require.Equal(t, arg.CurrencyCode, currency.CurrencyCode)
	require.Equal(t, arg.CurrencyName, currency.CurrencyName)
	require.Equal(t, arg.Symbol, currency.Symbol)
}

// TestCreateCurrency verifies that a new currency can be created successfully.
func TestCreateCurrency(t *testing.T) {
	CleanupDB(t, testDB)
	createRandomCurrency(t)
}

// TestGetCurrency verifies that a currency can be retrieved correctly after creation
// and that all fields match the original values.
func TestGetCurrency(t *testing.T) {
	CleanupDB(t, testDB)
	ctx := context.Background()
	currency1 := createRandomCurrency(t)

	currency2, err := testDB.Queries.GetCurrency(ctx, currency1.CurrencyCode)
	require.NoError(t, err)
	require.NotEmpty(t, currency2)

	require.Equal(t, currency1.CurrencyCode, currency2.CurrencyCode)
	require.Equal(t, currency1.CurrencyName, currency2.CurrencyName)
	require.Equal(t, currency1.Symbol, currency2.Symbol)
	require.Equal(t, currency1.ExchangeRate, currency2.ExchangeRate)
	require.Equal(t, currency1.IsActive, currency2.IsActive)
	require.WithinDuration(t, currency1.CreatedAt, currency2.CreatedAt, time.Second)
}

// TestUpdateExchangeRate verifies that a currency's exchange rate can be updated
// successfully and that the update is reflected in the database.
func TestUpdateExchangeRate(t *testing.T) {
	CleanupDB(t, testDB)
	ctx := context.Background()
	currency1 := createRandomCurrency(t)

	formattedValue := strconv.FormatFloat(util.RandomFloat(1.1, 99.99), 'f', -1, 64)
	exchangeRate, _ := util.SetNumeric(formattedValue)

	arg := UpdateExchangeRateParams{
		CurrencyCode: currency1.CurrencyCode,
		ExchangeRate: exchangeRate,
	}

	currency2, err := testDB.Queries.UpdateExchangeRate(ctx, arg)
	require.NoError(t, err)
	require.NotEmpty(t, currency2)

	require.Equal(t, currency1.CurrencyCode, currency2.CurrencyCode)
}

// TestDeleteCurrency verifies that a currency can be soft deleted (marked as inactive)
// and that the deletion is reflected in the database.
func TestDeleteCurrency(t *testing.T) {
	ctx := context.Background()
	currency1 := createRandomCurrency(t)

	err := testDB.Queries.DeleteCurrency(ctx, currency1.CurrencyCode)
	require.NoError(t, err)

	currency2, err := testDB.Queries.GetCurrency(ctx, currency1.CurrencyCode)
	require.NoError(t, err)
	require.False(t, currency2.IsActive)
}
