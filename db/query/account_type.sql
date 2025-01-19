-- name: ListAccountTypes :many
SELECT * FROM account_types
WHERE is_active = true;

-- name: CreateAccountType :one
INSERT INTO account_types (
  account_type,
  description
) VALUES (
  $1, $2
) RETURNING *;

-- name: UpdateAccountType :one
UPDATE account_types 
SET account_type = $1
WHERE account_type = $2
RETURNING *;

-- name: DeleteAccountType :exec
UPDATE account_types 
SET is_active = false
WHERE account_type = $1;

-- name: HardDeleteAccountType :exec
DELETE FROM account_types
WHERE account_type = $1;
