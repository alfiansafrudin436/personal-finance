-- name: CreateTransaction :one
INSERT INTO transactions (user_id, account_id, category_id, to_account_id, amount, transaction_date, description)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, user_id, account_id, category_id, to_account_id, amount, transaction_date, description, created_at, updated_at;

-- name: ListTransactions :many
SELECT
    t.id,
    t.account_id,
    a.name AS account_name,
    t.to_account_id,
    ta.name AS to_account_name,
    t.category_id,
    c.name AS category_name,
    c.type AS category_type,
    c.icon AS category_icon,
    c.color AS category_color,
    t.amount::text AS amount,
    t.transaction_date,
    t.description,
    t.created_at,
    t.updated_at
FROM transactions t
JOIN accounts a ON a.id = t.account_id
JOIN categories c ON c.id = t.category_id
LEFT JOIN accounts ta ON ta.id = t.to_account_id
WHERE t.user_id = $1
  AND (sqlc.narg('start_date')::date IS NULL OR t.transaction_date >= sqlc.narg('start_date')::date)
  AND (sqlc.narg('end_date')::date IS NULL OR t.transaction_date <= sqlc.narg('end_date')::date)
  AND (sqlc.narg('account_id')::uuid IS NULL OR t.account_id = sqlc.narg('account_id')::uuid OR t.to_account_id = sqlc.narg('account_id')::uuid)
  AND (sqlc.narg('category_id')::uuid IS NULL OR t.category_id = sqlc.narg('category_id')::uuid)
  AND (sqlc.narg('type')::category_type IS NULL OR c.type = sqlc.narg('type')::category_type)
  AND (sqlc.narg('search')::text IS NULL OR t.description ILIKE '%' || sqlc.narg('search')::text || '%')
ORDER BY t.transaction_date DESC, t.created_at DESC
LIMIT sqlc.arg('row_limit') OFFSET sqlc.arg('row_offset');

-- name: CountTransactions :one
SELECT COUNT(*)
FROM transactions t
JOIN categories c ON c.id = t.category_id
WHERE t.user_id = $1
  AND (sqlc.narg('start_date')::date IS NULL OR t.transaction_date >= sqlc.narg('start_date')::date)
  AND (sqlc.narg('end_date')::date IS NULL OR t.transaction_date <= sqlc.narg('end_date')::date)
  AND (sqlc.narg('account_id')::uuid IS NULL OR t.account_id = sqlc.narg('account_id')::uuid OR t.to_account_id = sqlc.narg('account_id')::uuid)
  AND (sqlc.narg('category_id')::uuid IS NULL OR t.category_id = sqlc.narg('category_id')::uuid)
  AND (sqlc.narg('type')::category_type IS NULL OR c.type = sqlc.narg('type')::category_type)
  AND (sqlc.narg('search')::text IS NULL OR t.description ILIKE '%' || sqlc.narg('search')::text || '%');

-- name: GetTransactionByID :one
SELECT
    t.id,
    t.account_id,
    a.name AS account_name,
    t.to_account_id,
    ta.name AS to_account_name,
    t.category_id,
    c.name AS category_name,
    c.type AS category_type,
    c.icon AS category_icon,
    c.color AS category_color,
    t.amount::text AS amount,
    t.transaction_date,
    t.description,
    t.created_at,
    t.updated_at
FROM transactions t
JOIN accounts a ON a.id = t.account_id
JOIN categories c ON c.id = t.category_id
LEFT JOIN accounts ta ON ta.id = t.to_account_id
WHERE t.id = $1 AND t.user_id = $2
LIMIT 1;

-- name: GetTransactionRaw :one
SELECT id, user_id, account_id, category_id, to_account_id, amount::text AS amount, transaction_date, description
FROM transactions
WHERE id = $1 AND user_id = $2
LIMIT 1;

-- name: UpdateTransaction :one
UPDATE transactions
SET account_id = $3,
    category_id = $4,
    to_account_id = $5,
    amount = $6,
    transaction_date = $7,
    description = $8
WHERE id = $1 AND user_id = $2
RETURNING id, user_id, account_id, category_id, to_account_id, amount, transaction_date, description, created_at, updated_at;

-- name: DeleteTransaction :execrows
DELETE FROM transactions
WHERE id = $1 AND user_id = $2;

-- name: AdjustAccountBalance :exec
-- delta is applied as-is, so callers pass a negative value to debit.
UPDATE accounts
SET balance = balance + sqlc.arg('delta')::numeric
WHERE id = $1 AND user_id = $2;

-- name: GetAccountForTransaction :one
SELECT id, name, is_archived
FROM accounts
WHERE id = $1 AND user_id = $2
LIMIT 1;

-- name: GetCategoryForTransaction :one
SELECT id, name, type
FROM categories
WHERE id = $1 AND (user_id = $2 OR user_id IS NULL)
LIMIT 1;
