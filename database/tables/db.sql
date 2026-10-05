CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    username VARCHAR(50) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    reset_password_token VARCHAR(255) NULL
);

CREATE TABLE accounts (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    user_id UUID NOT NULL,
    name VARCHAR(50) NOT NULL,
    type account_type NOT NULL,
    balance NUMERIC(15,2) NOT NULL DEFAULT 0.00,
    currency VARCHAR(3) NOT NULL DEFAULT 'IDR',
    is_archived BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_accounts_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    CONSTRAINT uq_accounts_user_name
        UNIQUE(user_id, name)
);

CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    user_id UUID NULL,
    name VARCHAR(50) NOT NULL,
    type category_type NOT NULL,
    icon VARCHAR(50),
    color VARCHAR(7),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_categories_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    CONSTRAINT uq_categories_user_name_type
        UNIQUE(user_id, name, type)
);

-- UNIQUE(user_id, ...) does not constrain global categories because NULL user_id
-- never conflicts in Postgres; this partial index covers them.
CREATE UNIQUE INDEX uq_categories_global_name_type
ON categories(name, type)
WHERE user_id IS NULL;

CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT uuidv7(),

    user_id UUID NOT NULL,
    account_id UUID NOT NULL,
    category_id UUID NOT NULL,
    to_account_id UUID NULL,

    amount NUMERIC(15,2) NOT NULL,
    transaction_date DATE NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_transactions_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_transactions_account
        FOREIGN KEY (account_id)
        REFERENCES accounts(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_transactions_category
        FOREIGN KEY (category_id)
        REFERENCES categories(id)
        ON DELETE RESTRICT,

    CONSTRAINT fk_transactions_to_account
        FOREIGN KEY (to_account_id)
        REFERENCES accounts(id)
        ON DELETE SET NULL
);

CREATE INDEX idx_transactions_user_date
ON transactions(user_id, transaction_date);

CREATE INDEX idx_transactions_category
ON transactions(category_id);

CREATE INDEX idx_transactions_account
ON transactions(account_id);

CREATE INDEX idx_accounts_user
ON accounts(user_id);

CREATE TABLE budgets (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    user_id UUID NOT NULL,
    category_id UUID NOT NULL,
    amount NUMERIC(15,2) NOT NULL,
    period_month DATE NOT NULL, -- always the 1st day of the budgeted month
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_budgets_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_budgets_category
        FOREIGN KEY (category_id)
        REFERENCES categories(id)
        ON DELETE CASCADE,

    CONSTRAINT uq_budgets_user_category_period
        UNIQUE(user_id, category_id, period_month),

    CONSTRAINT ck_budgets_amount_positive
        CHECK (amount > 0)
);

CREATE INDEX idx_budgets_user_period
ON budgets(user_id, period_month);

CREATE INDEX idx_transactions_to_account
ON transactions(to_account_id);
