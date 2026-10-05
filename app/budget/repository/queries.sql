-- name: UpsertBudget :one
INSERT INTO budgets (user_id, category_id, amount, period_month)
VALUES ($1, $2, $3, $4)
ON CONFLICT (user_id, category_id, period_month)
DO UPDATE SET amount = EXCLUDED.amount, updated_at = NOW()
RETURNING id, user_id, category_id, amount, period_month, created_at, updated_at;

-- name: ListBudgets :many
-- Pairs each budget with the amount actually spent in that month.
SELECT
    b.id,
    b.category_id,
    c.name AS category_name,
    c.icon AS category_icon,
    c.color AS category_color,
    b.amount::text AS amount,
    b.period_month,
    COALESCE((
        SELECT SUM(t.amount)
        FROM transactions t
        WHERE t.user_id = b.user_id
          AND t.category_id = b.category_id
          AND t.transaction_date >= b.period_month
          AND t.transaction_date < (b.period_month + INTERVAL '1 month')
    ), 0)::text AS spent,
    b.created_at,
    b.updated_at
FROM budgets b
JOIN categories c ON c.id = b.category_id
WHERE b.user_id = $1 AND b.period_month = $2
ORDER BY c.name;

-- name: GetBudgetByID :one
SELECT id, user_id, category_id, amount, period_month, created_at, updated_at
FROM budgets
WHERE id = $1 AND user_id = $2
LIMIT 1;

-- name: UpdateBudget :one
UPDATE budgets
SET amount = $3
WHERE id = $1 AND user_id = $2
RETURNING id, user_id, category_id, amount, period_month, created_at, updated_at;

-- name: DeleteBudget :execrows
DELETE FROM budgets
WHERE id = $1 AND user_id = $2;
