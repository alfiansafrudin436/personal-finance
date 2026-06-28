-- name: GetUserByEmailOrUsername :one
SELECT id, username, email, password_hash
FROM users
WHERE email = $1 OR username = $1
LIMIT 1;
