package transaction

import (
	"context"
	"database/sql"
	"fmt"
	"math/big"

	"github.com/jackc/pgtype"
	db "github.com/riad/banksystemendtoend/db/sqlc"
	"github.com/riad/banksystemendtoend/util/config"
	setup "github.com/riad/banksystemendtoend/util/db"
)

var ctx = context.Background()

// NegateNumeric creates a new pgtype.Numeric with negated value while preserving other properties
func NegateNumeric(num pgtype.Numeric) pgtype.Numeric {
	negated := num
	negated.Int = new(big.Int).Neg(num.Int)
	return negated
}

// CreateTransactionStatus creates a new transaction status or returns an existing one.
func CreateTransactionStatus(transactionStatus string) (db.TransactionStatus, error) {
	store, err := db.GetSQLStore(setup.GetStore())
	if err != nil {
		return db.TransactionStatus{}, fmt.Errorf("failed to get SQL store: %v", err)
	}

	existingStatus, err := store.Queries.GetTransactionStatus(ctx, transactionStatus)
	if err == nil && existingStatus.StatusCode == transactionStatus {
		return existingStatus, nil
	}

	descriptions := map[string]string{
		config.TransactionStatuses.PENDING:   "Transaction is being processed",
		config.TransactionStatuses.COMPLETED: "Transaction completed successfully",
		config.TransactionStatuses.FAILED:    "Transaction failed to process",
		config.TransactionStatuses.CANCELLED: "Transaction was cancelled",
		config.TransactionStatuses.REVERSED:  "Transaction was reversed",
	}

	arg := db.CreateTransactionStatusParams{
		StatusCode: transactionStatus,
		Description: func() string {
			description, ok := descriptions[transactionStatus]
			if !ok {
				return config.TransactionStatuses.PENDING
			}
			return description
		}(),
	}

	status, err := store.Queries.CreateTransactionStatus(ctx, arg) // Use the passed context here
	if err != nil {
		return db.TransactionStatus{}, fmt.Errorf("failed to create transaction status: %v", err)
	}
	return status, nil
}

// CreateCurrencyCode creates a new currency or returns an existing one.
func CreateCurrencyCode(currencyCode string) (db.AccountCurrency, error) {
	store, err := db.GetSQLStore(setup.GetStore())
	if err != nil {
		return db.AccountCurrency{}, fmt.Errorf("failed to get SQL store: %v", err)
	}

	existingCurrency, err := store.Queries.GetCurrency(context.Background(), currencyCode)
	if err == nil && existingCurrency.CurrencyCode == currencyCode {
		return existingCurrency, nil
	}

	var currencyName, symbol sql.NullString
	var rate sql.NullFloat64

	currencies := map[string]struct {
		name   string
		symbol string
		rate   float64
	}{
		config.TransactionCurrencies.USD.CODE: {
			name:   config.TransactionCurrencies.USD.NAME,
			symbol: config.TransactionCurrencies.USD.SYMBOL,
			rate:   config.TransactionCurrencies.USD.RATE,
		},
		config.TransactionCurrencies.EUR.CODE: {
			name:   config.TransactionCurrencies.EUR.NAME,
			symbol: config.TransactionCurrencies.EUR.SYMBOL,
			rate:   config.TransactionCurrencies.EUR.RATE,
		},
		config.TransactionCurrencies.GBP.CODE: {
			name:   config.TransactionCurrencies.GBP.NAME,
			symbol: config.TransactionCurrencies.GBP.SYMBOL,
			rate:   config.TransactionCurrencies.GBP.RATE,
		},
		config.TransactionCurrencies.JPY.CODE: {
			name:   config.TransactionCurrencies.JPY.NAME,
			symbol: config.TransactionCurrencies.JPY.SYMBOL,
			rate:   config.TransactionCurrencies.JPY.RATE,
		},
	}

	if currency, ok := currencies[currencyCode]; ok {
		currencyName = sql.NullString{String: currency.name, Valid: true}
		symbol = sql.NullString{String: currency.symbol, Valid: true}
		rate = sql.NullFloat64{Float64: currency.rate, Valid: true}
	} else {
		currencyName = sql.NullString{String: config.TransactionCurrencies.USD.NAME, Valid: true}
		symbol = sql.NullString{String: config.TransactionCurrencies.USD.SYMBOL, Valid: true}
		rate = sql.NullFloat64{Float64: config.TransactionCurrencies.USD.RATE, Valid: true}
	}

	var exchangeRate pgtype.Numeric
	err = exchangeRate.Set(rate.Float64)
	if err != nil {
		return db.AccountCurrency{}, fmt.Errorf("failed to convert exchange rate: %v", err)
	}

	arg := db.CreateCurrencyParams{
		CurrencyCode: currencyCode,
		CurrencyName: currencyName.String,
		Symbol:       sql.NullString{String: symbol.String, Valid: true},
		ExchangeRate: exchangeRate,
	}

	currency, err := store.Queries.CreateCurrency(context.Background(), arg)
	if err != nil {
		return db.AccountCurrency{}, fmt.Errorf("failed to create currency: %v", err)
	}

	return currency, nil
}

// CreateTransactionType creates a new transaction type or returns an existing one.
func CreateTransactionType(transactionType string) (db.TransactionType, error) {
	store, err := db.GetSQLStore(setup.GetStore())
	if err != nil {
		return db.TransactionType{}, fmt.Errorf("failed to get SQL store: %v", err)
	}

	existingType, err := store.Queries.GetTransactionType(context.Background(), transactionType)

	if err == nil && existingType.TypeCode == transactionType {
		return existingType, nil
	}

	descriptions := map[string]string{
		config.TransactionTypes.TRANSFER:   "Transfer of funds between accounts",
		config.TransactionTypes.DEPOSIT:    "Deposit of funds into account",
		config.TransactionTypes.WITHDRAWAL: "Withdrawal of funds from account",
		config.TransactionTypes.PAYMENT:    "Payment for goods or services",
		config.TransactionTypes.REFUND:     "Refund of previous transaction",
		config.TransactionTypes.ADJUSTMENT: "Account balance adjustment",
		config.TransactionTypes.FEE:        "Service or transaction fee",
		config.TransactionTypes.INTEREST:   "Interest earned or charged",
	}

	arg := db.CreateTransactionTypeParams{
		TypeCode: transactionType,
		Description: func() string {
			description, ok := descriptions[transactionType]
			if !ok {
				return config.TransactionTypes.TRANSFER
			}
			return description
		}(),
	}

	transType, err := store.Queries.CreateTransactionType(context.Background(), arg)
	if err != nil {
		return db.TransactionType{}, fmt.Errorf("failed to create transaction type: %v", err)
	}
	return transType, nil
}
