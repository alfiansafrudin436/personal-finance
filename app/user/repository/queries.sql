-- name: GetAllUsers :many
SELECT id, username, email, is_active, created_at
FROM users;

-- name: GetUserByEmail :one
SELECT id, username, email, password_hash
FROM users
WHERE email = $1
LIMIT 1;

-- name: GetUserByID :one
SELECT id, username, email
FROM users
WHERE id = $1
LIMIT 1;

-- name: UpdateUserUsername :exec
UPDATE users
SET username = $2,
    updated_at = now()
WHERE id = $1;

-- name: DeactivateUser :exec
UPDATE users
SET is_active = false,
    updated_at = now()
WHERE id = $1;