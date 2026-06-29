-- name: GetUserByEmail :one
SELECT id, username, email, password_hash
FROM users
WHERE email = $1
LIMIT 1;


-- name: CreateUser :one
INSERT INTO users (email, username, password_hash, created_at)
VALUES ($1, $2, $3, now())
RETURNING id;