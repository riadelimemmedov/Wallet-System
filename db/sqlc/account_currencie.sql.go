// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: account_currencie.sql

package db

import (
	"context"
	"database/sql"

	"github.com/jackc/pgtype"
)

const createCurrency = `-- name: CreateCurrency :one
INSERT INTO account_currencies (
    currency_code,
    currency_name,
    symbol,
    exchange_rate,
    last_updated_at
) VALUES (
    $1,
    $2,
    $3,
    $4,
    CURRENT_TIMESTAMP
) RETURNING currency_code, currency_name, symbol, is_active, exchange_rate, last_updated_at, created_at, updated_at
`

type CreateCurrencyParams struct {
	CurrencyCode string         `json:"currency_code"`
	CurrencyName string         `json:"currency_name"`
	Symbol       sql.NullString `json:"symbol"`
	ExchangeRate pgtype.Numeric `json:"exchange_rate"`
}

func (q *Queries) CreateCurrency(ctx context.Context, arg CreateCurrencyParams) (AccountCurrency, error) {
	row := q.db.QueryRow(ctx, createCurrency,
		arg.CurrencyCode,
		arg.CurrencyName,
		arg.Symbol,
		arg.ExchangeRate,
	)
	var i AccountCurrency
	err := row.Scan(
		&i.CurrencyCode,
		&i.CurrencyName,
		&i.Symbol,
		&i.IsActive,
		&i.ExchangeRate,
		&i.LastUpdatedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteCurrency = `-- name: DeleteCurrency :exec
UPDATE account_currencies
SET is_active = false
WHERE currency_code = $1
`

func (q *Queries) DeleteCurrency(ctx context.Context, currencyCode string) error {
	_, err := q.db.Exec(ctx, deleteCurrency, currencyCode)
	return err
}

const getCurrency = `-- name: GetCurrency :one
SELECT currency_code, currency_name, symbol, is_active, exchange_rate, last_updated_at, created_at, updated_at FROM account_currencies
WHERE currency_code = $1
`

func (q *Queries) GetCurrency(ctx context.Context, currencyCode string) (AccountCurrency, error) {
	row := q.db.QueryRow(ctx, getCurrency, currencyCode)
	var i AccountCurrency
	err := row.Scan(
		&i.CurrencyCode,
		&i.CurrencyName,
		&i.Symbol,
		&i.IsActive,
		&i.ExchangeRate,
		&i.LastUpdatedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const hardDeleteCurrency = `-- name: HardDeleteCurrency :exec
DELETE FROM account_currencies
WHERE currency_code = $1
`

func (q *Queries) HardDeleteCurrency(ctx context.Context, currencyCode string) error {
	_, err := q.db.Exec(ctx, hardDeleteCurrency, currencyCode)
	return err
}

const listCurrencies = `-- name: ListCurrencies :many
SELECT currency_code, currency_name, symbol, is_active, exchange_rate, last_updated_at, created_at, updated_at FROM account_currencies
WHERE is_active = true
`

func (q *Queries) ListCurrencies(ctx context.Context) ([]AccountCurrency, error) {
	rows, err := q.db.Query(ctx, listCurrencies)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []AccountCurrency{}
	for rows.Next() {
		var i AccountCurrency
		if err := rows.Scan(
			&i.CurrencyCode,
			&i.CurrencyName,
			&i.Symbol,
			&i.IsActive,
			&i.ExchangeRate,
			&i.LastUpdatedAt,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateExchangeRate = `-- name: UpdateExchangeRate :one
UPDATE account_currencies
SET 
    exchange_rate = $1,
    last_updated_at = CURRENT_TIMESTAMP
WHERE currency_code = $2
RETURNING currency_code, currency_name, symbol, is_active, exchange_rate, last_updated_at, created_at, updated_at
`

type UpdateExchangeRateParams struct {
	ExchangeRate pgtype.Numeric `json:"exchange_rate"`
	CurrencyCode string         `json:"currency_code"`
}

func (q *Queries) UpdateExchangeRate(ctx context.Context, arg UpdateExchangeRateParams) (AccountCurrency, error) {
	row := q.db.QueryRow(ctx, updateExchangeRate, arg.ExchangeRate, arg.CurrencyCode)
	var i AccountCurrency
	err := row.Scan(
		&i.CurrencyCode,
		&i.CurrencyName,
		&i.Symbol,
		&i.IsActive,
		&i.ExchangeRate,
		&i.LastUpdatedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
