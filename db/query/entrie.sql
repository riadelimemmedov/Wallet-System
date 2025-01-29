-- name: CreateEntry :one
INSERT INTO entries (
    account_id,
    amount
) VALUES (
    $1, $2
) RETURNING *;

-- name: GetEntry :one
SELECT * FROM entries
WHERE id = $1;

-- name: ListEntries :many
SELECT * FROM entries
WHERE account_id = $1
ORDER BY id DESC
LIMIT $2 OFFSET $3;

-- name: GetAccountBalance :one
SELECT COALESCE(SUM(amount), 0) as balance
FROM entries
WHERE account_id = $1;

-- name: GetAccountStatement :many
SELECT created_at, amount 
FROM entries
WHERE account_id = $1
AND created_at BETWEEN $2 AND $3
ORDER BY created_at DESC;

-- name: HardDeleteEntries :exec
DELETE FROM entries
WHERE account_id = $1;
