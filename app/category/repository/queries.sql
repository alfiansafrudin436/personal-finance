-- name: CreateCategory :one
INSERT INTO categories (user_id, name, type, icon, color)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, user_id, name, type, icon, color, created_at, updated_at;

-- name: ListCategories :many
-- Returns the user's own categories plus the global ones (user_id IS NULL).
SELECT id, user_id, name, type, icon, color, created_at, updated_at,
       (user_id IS NULL)::boolean AS is_global
FROM categories
WHERE (user_id = $1 OR user_id IS NULL)
  AND (sqlc.narg('type')::category_type IS NULL OR type = sqlc.narg('type')::category_type)
ORDER BY type, name;

-- name: GetCategoryByID :one
SELECT id, user_id, name, type, icon, color, created_at, updated_at,
       (user_id IS NULL)::boolean AS is_global
FROM categories
WHERE id = $1 AND (user_id = $2 OR user_id IS NULL)
LIMIT 1;

-- name: UpdateCategory :one
-- Scoped to user_id so global categories can never be edited.
UPDATE categories
SET name = $3,
    type = $4,
    icon = $5,
    color = $6
WHERE id = $1 AND user_id = $2
RETURNING id, user_id, name, type, icon, color, created_at, updated_at;

-- name: DeleteCategory :execrows
DELETE FROM categories
WHERE id = $1 AND user_id = $2;

-- name: CountCategoryTransactions :one
SELECT COUNT(*)
FROM transactions
WHERE category_id = $1;
