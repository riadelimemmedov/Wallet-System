-- name: CreateUser :one
INSERT INTO users (
    username,
    password_hash,
    email,
    first_name,
    last_name,
    phone_number,
    profile_image_url
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE user_id = $1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY user_id
LIMIT $1 OFFSET $2;

-- name: UpdateUser :one
UPDATE users 
SET 
    username = COALESCE(sqlc.narg('username'), username),
    email = COALESCE(sqlc.narg('email'), email),
    first_name = COALESCE(sqlc.narg('first_name'), first_name),
    last_name = COALESCE(sqlc.narg('last_name'), last_name),
    phone_number = COALESCE(sqlc.narg('phone_number'), phone_number),
    profile_image_url = COALESCE(sqlc.narg('profile_image_url'), profile_image_url),
    is_active = COALESCE(sqlc.narg('is_active'), is_active)
WHERE user_id = sqlc.arg('user_id')
RETURNING *;

-- name: DeleteUser :exec
UPDATE users
SET is_active = false
WHERE user_id = $1;

-- name: HardDeleteUser :exec
DELETE FROM users
WHERE user_id = $1;