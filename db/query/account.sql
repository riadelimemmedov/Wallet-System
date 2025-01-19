-- name: CreateAccount :one
INSERT INTO accounts (
    user_id,
    account_number,
    account_type,
    currency_code,
    interest_rate,
    overdraft_limit
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetAccount :one
SELECT * FROM accounts
WHERE account_id = $1;

-- name: GetAccountForUpdate :one
SELECT * FROM accounts
WHERE account_id = $1
FOR NO KEY UPDATE;

-- name: ListAccountsByUser :many
SELECT * FROM accounts
WHERE user_id = $1
ORDER BY account_id;

-- name: UpdateAccountBalance :one
UPDATE accounts 
SET balance = balance + sqlc.arg('amount')
WHERE account_id = sqlc.arg('account_id')
RETURNING *;

-- name: DeleteAccount :exec
UPDATE accounts
SET is_active = false
WHERE account_id = $1;

-- name: HardDeleteAccount :exec
DELETE FROM accounts
WHERE account_id = $1;