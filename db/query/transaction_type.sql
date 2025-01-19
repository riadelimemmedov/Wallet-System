-- name: CreateTransactionType :one
INSERT INTO transaction_types (
    type_code,
    description
) VALUES (
    $1, $2
) RETURNING *;

-- name: GetTransactionType :one
SELECT * FROM transaction_types
WHERE type_code = $1;

-- name: ListTransactionTypes :many
SELECT * FROM transaction_types
WHERE is_active = true
ORDER BY type_code;

-- name: UpdateTransactionType :one
UPDATE transaction_types 
SET 
    description = COALESCE(sqlc.narg('description'), description),
    is_active = COALESCE(sqlc.narg('is_active'), is_active)
WHERE type_code = sqlc.arg('type_code')
RETURNING *;

-- name: DeleteTransactionType :exec
UPDATE transaction_types
SET is_active = false
WHERE type_code = $1;

-- name: HardDeleteTransactionType :exec
DELETE FROM transaction_types
WHERE type_code = $1;