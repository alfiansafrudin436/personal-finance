-- name: GetUserByEmail :one
SELECT id, username, email, password_hash
FROM users
WHERE email = $1
LIMIT 1;
