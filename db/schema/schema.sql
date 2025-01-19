--! Users table - stores user information
CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    phone_number VARCHAR(20),
    is_active BOOLEAN NOT NULL DEFAULT true,
    last_login TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT users_email_check CHECK (
        email ~* '^[A-Za-z0-9._+%-]+@[A-Za-z0-9.-]+[.][A-Za-z]+$'
    )
);

CREATE INDEX idx_users_composite ON users(user_id, username, email, is_active);

--! Account Types Reference Table
CREATE TABLE account_types (
    account_type VARCHAR(50) PRIMARY KEY,
    description TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

--! Account Currencies Reference Table
CREATE TABLE account_currencies (
    currency_code VARCHAR(3) PRIMARY KEY,
    currency_name VARCHAR(50) NOT NULL,
    symbol VARCHAR(5),
    is_active BOOLEAN NOT NULL DEFAULT true,
    exchange_rate DECIMAL(10, 6),
    last_updated_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

--! Accounts table - stores account information
CREATE TABLE accounts (
    account_id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(user_id),
    account_number VARCHAR(20) NOT NULL UNIQUE,
    account_type VARCHAR(50) NOT NULL REFERENCES account_types(account_type),
    balance DECIMAL(15, 2) NOT NULL DEFAULT 0.00,
    currency_code VARCHAR(3) NOT NULL REFERENCES account_currencies(currency_code),
    interest_rate DECIMAL(5, 2),
    overdraft_limit DECIMAL(15, 2) DEFAULT 0.00,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT accounts_balance_check CHECK (balance >= - overdraft_limit)
);

CREATE INDEX idx_accounts_composite ON accounts(account_id, user_id);

--! Transaction Types Reference Table
CREATE TABLE transaction_types (
    type_code VARCHAR(50) PRIMARY KEY,
    description TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

--! Transaction Status Reference Table
CREATE TABLE transaction_status (
    status_code VARCHAR(50) PRIMARY KEY,
    description TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

--! Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

--! Transactions table - stores all transactions
CREATE TABLE transactions (
    transaction_id SERIAL PRIMARY KEY,
    transaction_number UUID DEFAULT uuid_generate_v4() NOT NULL UNIQUE,
    from_account_id INTEGER REFERENCES accounts(account_id),
    to_account_id INTEGER REFERENCES accounts(account_id),
    type_code VARCHAR(50) NOT NULL REFERENCES transaction_types(type_code),
    amount DECIMAL(15, 2) NOT NULL,
    currency_code VARCHAR(3) NOT NULL REFERENCES account_currencies(currency_code),
    exchange_rate DECIMAL(10, 6),
    status_code VARCHAR(50) NOT NULL REFERENCES transaction_status(status_code),
    description TEXT,
    reference_number VARCHAR(50) UNIQUE,
    transaction_date TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT transactions_amount_check CHECK (amount > 0),
    CONSTRAINT transactions_accounts_check CHECK (
        (from_account_id IS NOT NULL OR to_account_id IS NOT NULL) AND
        (from_account_id != to_account_id OR (from_account_id IS NULL OR to_account_id IS NULL))
    )
);

--! Add index for UUID
CREATE INDEX idx_transactions_number ON transactions(transaction_number);
CREATE INDEX idx_transactions_accounts ON transactions(from_account_id, to_account_id);
CREATE INDEX idx_transactions_date ON transactions(transaction_date);

--! Audit Trail table - tracks all important changes
CREATE TABLE audit_trail (
    audit_id SERIAL PRIMARY KEY,
    table_name VARCHAR(50) NOT NULL,
    record_id INTEGER NOT NULL,
    action VARCHAR(10) NOT NULL,
    old_values JSONB,
    new_values JSONB,
    user_id INTEGER REFERENCES users(user_id),
    ip_address VARCHAR(45),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_audit_trail_composite ON audit_trail(table_name, record_id);


-- Insert basic reference data
-- INSERT INTO
--     account_types (account_type, description)
-- VALUES
--     ('SAVINGS', 'Standard savings account'),
--     ('CHECKING', 'Standard checking account'),
--     ('FIXED_DEPOSIT', 'Fixed deposit account'),
--     ('MONEY_MARKET', 'Money market account');

-- INSERT INTO
--     account_currencies (
--         currency_code,
--         currency_name,
--         symbol,
--         exchange_rate
--     )
-- VALUES
--     ('USD', 'US Dollar', '$', 1.0),
--     ('EUR', 'Euro', '€', 1.08),
--     ('GBP', 'British Pound', '£', 1.26),
--     ('JPY', 'Japanese Yen', '¥', 0.0067);

-- INSERT INTO
--     transaction_types (type_code, description)
-- VALUES
--     ('DEPOSIT', 'Cash or check deposit'),
--     ('WITHDRAWAL', 'Cash withdrawal'),
--     ('TRANSFER', 'Transfer between accounts'),
--     ('PAYMENT', 'Bill payment'),
--     ('FEE', 'Service fee'),
--     ('INTEREST', 'Interest credit');

-- INSERT INTO
--     transaction_status (status_code, description)
-- VALUES
--     ('PENDING', 'Transaction is pending'),
--     (
--         'COMPLETED',
--         'Transaction completed successfully'
--     ),
--     ('FAILED', 'Transaction failed'),
--     ('REVERSED', 'Transaction was reversed'),
--     ('CANCELLED', 'Transaction was cancelled');


-- -- First, create a test user
-- INSERT INTO users (username, password_hash, email, first_name, last_name, phone_number) 
-- VALUES 
-- ('john_doe', '$2a$10$vI8aWBnW3fID.ZQ4/zo1G.q1lRps.9cGLcZEiGDMVr5yUP1KUOYTa', 'john@example.com', 'John', 'Doe', '+1234567890'),
-- ('jane_smith', '$2a$10$vI8aWBnW3fID.ZQ4/zo1G.q1lRps.9cGLcZEiGDMVr5yUP1KUOYTa', 'jane@example.com', 'Jane', 'Smith', '+1234567891');

-- -- Create accounts for both users
-- INSERT INTO accounts (user_id, account_number, account_type, balance, currency_code, interest_rate) 
-- VALUES 
-- -- John's accounts
-- (1, 'ACC100001', 'SAVINGS', 1000.00, 'USD', 2.5),
-- (1, 'ACC100002', 'CHECKING', 500.00, 'USD', 0.5),
-- -- Jane's accounts
-- (2, 'ACC100003', 'SAVINGS', 1500.00, 'EUR', 2.0),
-- (2, 'ACC100004', 'CHECKING', 750.00, 'EUR', 0.5);

-- -- Create some transactions
-- INSERT INTO transactions (
--     transaction_number,
--     from_account_id,
--     to_account_id,
--     type_code,
--     amount,
--     currency_code,
--     exchange_rate,
--     status_code,
--     description,
--     reference_number,
--     transaction_date
-- ) VALUES 
-- -- Deposit to John's savings
-- (uuid_generate_v4(), NULL, 1, 'DEPOSIT', 1000.00, 'USD', 1.0, 'COMPLETED', 'Initial deposit', 'REF001', CURRENT_TIMESTAMP),

-- -- Transfer from John's savings to checking
-- (uuid_generate_v4(), 1, 2, 'TRANSFER', 200.00, 'USD', 1.0, 'COMPLETED', 'Transfer to checking', 'REF002', CURRENT_TIMESTAMP),

-- -- Jane's deposit
-- (uuid_generate_v4(), NULL, 3, 'DEPOSIT', 1500.00, 'EUR', 1.08, 'COMPLETED', 'Initial deposit', 'REF003', CURRENT_TIMESTAMP),

-- -- International transfer from Jane to John
-- (uuid_generate_v4(), 3, 1, 'TRANSFER', 500.00, 'EUR', 1.08, 'COMPLETED', 'International transfer', 'REF004', CURRENT_TIMESTAMP),

-- -- Failed transaction
-- (uuid_generate_v4(), 2, 3, 'TRANSFER', 5000.00, 'USD', 1.0, 'FAILED', 'Insufficient funds', 'REF005', CURRENT_TIMESTAMP);

-- -- Create some audit entries
-- INSERT INTO audit_trail (table_name, record_id, action, old_values, new_values, user_id, ip_address)
-- VALUES 
-- ('accounts', 1, 'UPDATE', 
--   '{"balance": 1000.00}'::jsonb,
--   '{"balance": 800.00}'::jsonb,
--   1, '192.168.1.1'),
-- ('transactions', 1, 'INSERT',
--   NULL,
--   '{"amount": 1000.00, "type": "DEPOSIT"}'::jsonb,
--   1, '192.168.1.1'),
-- ('users', 2, 'UPDATE',
--   '{"email": "jane@old.com"}'::jsonb,
--   '{"email": "jane@example.com"}'::jsonb,
--   2, '192.168.1.2');


--! Create trigger function for updating timestamps
CREATE
OR REPLACE FUNCTION update_updated_at_column() RETURNS TRIGGER AS $ $ BEGIN NEW.updated_at = CURRENT_TIMESTAMP;

RETURN NEW;

END;

$ $ language 'plpgsql';

--! Create triggers for updating timestamps
CREATE TRIGGER update_users_updated_at BEFORE
UPDATE
    ON users FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_accounts_updated_at BEFORE
UPDATE
    ON accounts FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_transactions_updated_at BEFORE
UPDATE
    ON transactions FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();