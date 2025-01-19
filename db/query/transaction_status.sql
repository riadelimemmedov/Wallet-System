-- name: CreateTransactionStatus :one
INSERT INTO transaction_status (
    status_code,
    description
) VALUES (
    $1, $2
) RETURNING *;

-- name: GetTransactionStatus :one
SELECT * FROM transaction_status
WHERE status_code = $1;

-- name: ListTransactionStatus :many
SELECT * FROM transaction_status
WHERE is_active = true
ORDER BY status_code;

-- name: ModifyTransactionStatus :one
UPDATE transaction_status 
SET 
    description = COALESCE(sqlc.narg('description'), description),
    is_active = COALESCE(sqlc.narg('is_active'), is_active)
WHERE status_code = sqlc.arg('status_code')
RETURNING *;

-- name: DeleteTransactionStatus :exec
UPDATE transaction_status
SET is_active = false
WHERE status_code = $1;

-- name: HardDeleteTransactionStatus :exec
DELETE FROM transaction_status
WHERE status_code = $1;

-- name: GetActiveTransactionStatus :many
SELECT * FROM transaction_status
WHERE is_active = true
AND status_code = ANY($1::varchar[]);