-- name: GetPeriodSummary :one
-- Transfers are excluded: moving money between own accounts is not income or expense.
SELECT
    COALESCE(SUM(CASE WHEN c.type = 'income'  THEN t.amount ELSE 0 END), 0)::text AS total_income,
    COALESCE(SUM(CASE WHEN c.type = 'expense' THEN t.amount ELSE 0 END), 0)::text AS total_expense,
    COUNT(*) AS transaction_count
FROM transactions t
JOIN categories c ON c.id = t.category_id
WHERE t.user_id = $1
  AND t.transaction_date >= sqlc.arg('start_date')
  AND t.transaction_date <= sqlc.arg('end_date');

-- name: GetSpendingByCategory :many
SELECT
    c.id AS category_id,
    c.name AS category_name,
    c.icon AS category_icon,
    c.color AS category_color,
    SUM(t.amount)::text AS total,
    COUNT(*) AS transaction_count
FROM transactions t
JOIN categories c ON c.id = t.category_id
WHERE t.user_id = $1
  AND c.type = sqlc.arg('type')::category_type
  AND t.transaction_date >= sqlc.arg('start_date')
  AND t.transaction_date <= sqlc.arg('end_date')
GROUP BY c.id, c.name, c.icon, c.color
ORDER BY SUM(t.amount) DESC;

-- name: GetMonthlyTrend :many
-- One row per month in the requested window, including months with no activity.
SELECT
    to_char(m.month, 'YYYY-MM') AS period,
    COALESCE(SUM(CASE WHEN c.type = 'income'  THEN t.amount ELSE 0 END), 0)::text AS total_income,
    COALESCE(SUM(CASE WHEN c.type = 'expense' THEN t.amount ELSE 0 END), 0)::text AS total_expense
FROM generate_series(sqlc.arg('start_month')::date, sqlc.arg('end_month')::date, INTERVAL '1 month') AS m(month)
LEFT JOIN transactions t
       ON t.user_id = $1
      AND t.transaction_date >= m.month
      AND t.transaction_date < (m.month + INTERVAL '1 month')
LEFT JOIN categories c ON c.id = t.category_id
GROUP BY m.month
ORDER BY m.month;

-- name: GetAccountBalances :many
SELECT id, name, type, balance::text AS balance, currency
FROM accounts
WHERE user_id = $1 AND is_archived = FALSE
ORDER BY balance DESC;

-- name: GetTotalBalance :one
SELECT COALESCE(SUM(balance), 0)::text AS total_balance
FROM accounts
WHERE user_id = $1 AND is_archived = FALSE;

-- name: GetRecentTransactions :many
SELECT
    t.id,
    a.name AS account_name,
    c.name AS category_name,
    c.type AS category_type,
    c.icon AS category_icon,
    c.color AS category_color,
    t.amount::text AS amount,
    t.transaction_date,
    t.description
FROM transactions t
JOIN accounts a ON a.id = t.account_id
JOIN categories c ON c.id = t.category_id
WHERE t.user_id = $1
ORDER BY t.transaction_date DESC, t.created_at DESC
LIMIT sqlc.arg('row_limit');
