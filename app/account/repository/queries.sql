-- name: CreateAccount :one
INSERT INTO accounts (user_id, name, type, balance, currency)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, user_id, name, type, balance, currency, is_archived, created_at, updated_at;

-- name: ListAccounts :many
SELECT id, user_id, name, type, balance, currency, is_archived, created_at, updated_at
FROM accounts
WHERE user_id = $1
  AND (sqlc.narg('include_archived')::boolean IS TRUE OR is_archived = FALSE)
ORDER BY is_archived, name;

-- name: GetAccountByID :one
SELECT id, user_id, name, type, balance, currency, is_archived, created_at, updated_at
FROM accounts
WHERE id = $1 AND user_id = $2
LIMIT 1;

-- name: UpdateAccount :one
UPDATE accounts
SET name = $3,
    type = $4,
    currency = $5
WHERE id = $1 AND user_id = $2
RETURNING id, user_id, name, type, balance, currency, is_archived, created_at, updated_at;

-- name: SetAccountArchived :exec
UPDATE accounts
SET is_archived = $3
WHERE id = $1 AND user_id = $2;

-- name: DeleteAccount :execrows
DELETE FROM accounts
WHERE id = $1 AND user_id = $2;

-- name: CountAccountTransactions :one
SELECT COUNT(*)
FROM transactions
WHERE user_id = $1 AND (account_id = $2 OR to_account_id = $2);

-- name: GetAccountsSummary :one
SELECT
    COALESCE(SUM(balance), 0)::text AS total_balance,
    COUNT(*) AS total_accounts
FROM accounts
WHERE user_id = $1 AND is_archived = FALSE;
