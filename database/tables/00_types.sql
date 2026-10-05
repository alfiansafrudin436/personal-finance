-- Enum types. Must be created before the tables in db.sql reference them.

CREATE TYPE account_type AS ENUM (
    'cash',
    'bank',
    'e_wallet',
    'credit_card',
    'investment'
);

CREATE TYPE category_type AS ENUM (
    'income',
    'expense',
    'transfer'
);
