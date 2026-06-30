-- name: GetUserByEmail :one
SELECT id, username, email, password_hash
FROM users
WHERE email = $1
LIMIT 1;


-- name: CreateUser :one
INSERT INTO users (email, username, password_hash, created_at)
VALUES ($1, $2, $3, now())
RETURNING id;

-- name: UpdateUserPassword :exec
UPDATE users
SET password_hash = $2,
    reset_password_token = null,
    updated_at = now()
WHERE id = $1;

-- name: UpdateUserResetToken :exec
UPDATE users
SET reset_password_token = $2,
    updated_at = now()
WHERE id = $1;

-- name: GetUserByResetToken :one
SELECT id, username, email
FROM users
WHERE reset_password_token = $1
LIMIT 1;