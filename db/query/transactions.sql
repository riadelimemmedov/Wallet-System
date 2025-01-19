-- name: CreateTransaction :one
INSERT INTO transactions (
    from_account_id,
    to_account_id,
    type_code,
    amount,
    currency_code,
    exchange_rate,
    status_code,
    description,
    reference_number,
    transaction_date
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
) RETURNING *;

-- name: GetTransaction :one
SELECT * FROM transactions
WHERE transaction_id = $1;

-- name: ListTransactionsByAccount :many
SELECT * FROM transactions
WHERE from_account_id = $1 OR to_account_id = $1
ORDER BY transaction_date DESC
LIMIT $2 OFFSET $3;

-- name: UpdateTransactionStatus :one
UPDATE transactions
SET status_code = $2
WHERE transaction_id = $1
RETURNING *;

-- name: GetTransactionsByDateRange :many
SELECT * FROM transactions
WHERE transaction_date BETWEEN sqlc.arg('start_date') AND sqlc.arg('end_date')
ORDER BY transaction_date DESC;