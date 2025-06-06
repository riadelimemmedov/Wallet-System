-- name: CreateCurrency :one
INSERT INTO account_currencies (
    currency_code,
    currency_name,
    symbol,
    exchange_rate,
    last_updated_at
) VALUES (
    sqlc.arg('currency_code'),
    sqlc.arg('currency_name'),
    sqlc.arg('symbol'),
    sqlc.arg('exchange_rate'),
    CURRENT_TIMESTAMP
) RETURNING *;

-- name: ListCurrencies :many
SELECT * FROM account_currencies
WHERE is_active = true;

-- name: UpdateExchangeRate :one
UPDATE account_currencies
SET 
    exchange_rate = sqlc.arg('exchange_rate'),
    last_updated_at = CURRENT_TIMESTAMP
WHERE currency_code = sqlc.arg('currency_code')
RETURNING *;

-- name: GetCurrency :one
SELECT * FROM account_currencies
WHERE currency_code = $1;

-- name: DeleteCurrency :exec
UPDATE account_currencies
SET is_active = false
WHERE currency_code = $1;

-- name: HardDeleteCurrency :exec
DELETE FROM account_currencies
WHERE currency_code = $1;