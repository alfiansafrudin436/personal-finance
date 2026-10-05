-- 002_seed_default_categories.sql
-- Global categories (user_id IS NULL) available to every user.
-- Safe to run more than once.

INSERT INTO categories (user_id, name, type, icon, color)
VALUES
    (NULL, 'Salary',          'income',   'wallet',        '#16a34a'),
    (NULL, 'Bonus',           'income',   'gift',          '#22c55e'),
    (NULL, 'Investment',      'income',   'trending-up',   '#0ea5e9'),
    (NULL, 'Other Income',    'income',   'plus-circle',   '#14b8a6'),
    (NULL, 'Food & Drink',    'expense',  'utensils',      '#f97316'),
    (NULL, 'Transport',       'expense',  'car',           '#6366f1'),
    (NULL, 'Housing',         'expense',  'home',          '#8b5cf6'),
    (NULL, 'Utilities',       'expense',  'zap',           '#eab308'),
    (NULL, 'Health',          'expense',  'heart-pulse',   '#ef4444'),
    (NULL, 'Education',       'expense',  'book-open',     '#3b82f6'),
    (NULL, 'Entertainment',   'expense',  'clapperboard',  '#ec4899'),
    (NULL, 'Shopping',        'expense',  'shopping-bag',  '#d946ef'),
    (NULL, 'Other Expense',   'expense',  'minus-circle',  '#64748b'),
    (NULL, 'Transfer',        'transfer', 'arrow-right-left', '#0891b2')
ON CONFLICT DO NOTHING;
