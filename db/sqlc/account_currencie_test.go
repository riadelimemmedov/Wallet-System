package db

import (
	"context"
	"database/sql"
	"strconv"
	"testing"

	"github.com/riad/banksystemendtoend/util"
	"github.com/stretchr/testify/require"
)

func createRandomCurrency(t *testing.T) AccountCurrency {
	t.Helper()

	currency_code := util.RandomCurrency()
	formattedValue := strconv.FormatFloat(util.RandomFloat(1.1, 99.99), 'f', -1, 64)
	exchangeRate, err := util.SetNumeric(formattedValue)
	if err != nil {
		require.Error(t, err)
	}

	arg := CreateCurrencyParams{
		CurrencyCode: currency_code,
		CurrencyName: util.RandomCurrencyName(currency_code),
		Symbol: sql.NullString{
			String: util.RandomCurrencySymbol(currency_code),
			Valid:  true,
		},
		ExchangeRate: exchangeRate,
	}
	ctx := context.Background()
	currency, err := testDB.Queries.CreateCurrency(ctx, arg)
	require.NoError(t, err)
	require.NotEmpty(t, currency)

	require.Equal(t, arg.CurrencyCode, currency.CurrencyCode)
	require.Equal(t, arg.CurrencyName, currency.CurrencyName)
	require.Equal(t, arg.Symbol, currency.Symbol)
	require.Equal(t, arg.ExchangeRate, exchangeRate)
	return currency

}

func TestCreateCurrency(t *testing.T) {
	CleanupDB(t, testDB)
	createRandomCurrency(t)
}
