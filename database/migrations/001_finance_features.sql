-- 001_finance_features.sql
-- Upgrades an existing database to the schema in database/tables/.
-- Safe to run more than once.

BEGIN;

-- Enum types (database/tables/00_types.sql)
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'account_type') THEN
        CREATE TYPE account_type AS ENUM ('cash', 'bank', 'e_wallet', 'credit_card', 'investment');
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'category_type') THEN
        CREATE TYPE category_type AS ENUM ('income', 'expense', 'transfer');
    END IF;
END $$;

ALTER TABLE accounts   ADD COLUMN IF NOT EXISTS currency    VARCHAR(3) NOT NULL DEFAULT 'IDR';
ALTER TABLE accounts   ADD COLUMN IF NOT EXISTS is_archived BOOLEAN    NOT NULL DEFAULT FALSE;
ALTER TABLE categories ADD COLUMN IF NOT EXISTS color       VARCHAR(7);

CREATE TABLE IF NOT EXISTS budgets (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category_id UUID NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    amount NUMERIC(15,2) NOT NULL CHECK (amount > 0),
    period_month DATE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_budgets_user_category_period UNIQUE(user_id, category_id, period_month)
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_categories_global_name_type
ON categories(name, type) WHERE user_id IS NULL;

CREATE INDEX IF NOT EXISTS idx_budgets_user_period   ON budgets(user_id, period_month);
CREATE INDEX IF NOT EXISTS idx_transactions_to_account ON transactions(to_account_id);

DROP TRIGGER IF EXISTS trg_budgets_updated_at ON budgets;
CREATE TRIGGER trg_budgets_updated_at
BEFORE UPDATE ON budgets
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

COMMIT;
