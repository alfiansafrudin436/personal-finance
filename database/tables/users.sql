CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE accounts (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    name VARCHAR(50) NOT NULL,
    type account_type NOT NULL,
    balance DECIMAL(15, 2) DEFAULT 0.00,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    user_id INT NULL,
    name VARCHAR(50) NOT NULL,
    type category_type NOT NULL,
    icon VARCHAR(50) NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    account_id INT NOT NULL,
    category_id INT NOT NULL,
    to_account_id INT NULL, -- Digunakan hanya jika tipe transaksi adalah transfer
    amount DECIMAL(15, 2) NOT NULL,
    transaction_date DATE NOT NULL,
    description TEXT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE RESTRICT,
    FOREIGN KEY (to_account_id) REFERENCES accounts(id) ON DELETE SET NULL
);

-- Index untuk mempercepat pencarian transaksi berdasarkan user dan rentang tanggal (Sangat penting untuk halaman dashboard/history)
CREATE INDEX idx_transactions_user_date 
ON transactions(user_id, transaction_date);

-- Index untuk mempercepat filter laporan keuangan berdasarkan kategori
CREATE INDEX idx_transactions_category 
ON transactions(category_id);

-- Index untuk melihat mutasi saldo per rekening/dompet tertentu
CREATE INDEX idx_transactions_account 
ON transactions(account_id);

-- Index untuk mencari dompet milik user tertentu dengan cepat
CREATE INDEX idx_accounts_user 
ON accounts(user_id);