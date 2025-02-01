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

-- name: GetTransactionBalance :one
SELECT COALESCE(SUM(amount), 0) as balance
FROM transactions
WHERE from_account_id = $1 OR to_account_id = $1;

-- name: GetTransactionStatement :many
SELECT transaction_id, from_account_id, to_account_id, 
        amount, currency_code, status_code, description,
        transaction_date
FROM transactions
WHERE (from_account_id = $1 OR to_account_id = $1)
    AND transaction_date BETWEEN $2 AND $3
ORDER BY transaction_date DESC;

-- name: GetTransactionsByStatus :many
SELECT *
FROM transactions
WHERE status_code = $1
ORDER BY transaction_date DESC;

-- name: GetTransactionByReference :one
SELECT *
FROM transactions
WHERE reference_number = $1;

-- name: ListAccountTransactions :many
SELECT *
FROM transactions
WHERE from_account_id = $1 OR to_account_id = $1
ORDER BY transaction_date DESC
LIMIT $2 OFFSET $3;

-- name: DeleteTransaction :exec
DELETE FROM transactions 
WHERE transaction_number = $1;

-- name: DeleteAccountTransactions :exec
DELETE FROM transactions
WHERE from_account_id = $1 OR to_account_id = $1;

-- name: DeleteTransactionsByDateRange :exec
DELETE FROM transactions 
WHERE transaction_date BETWEEN sqlc.arg('start_date') AND sqlc.arg('end_date');

-- name: DeleteTransactionsByStatus :exec
DELETE FROM transactions 
WHERE status_code = $1;