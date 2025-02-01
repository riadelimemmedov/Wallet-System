-- Drop transaction_number from transactions table if exists
ALTER TABLE transactions
DROP COLUMN IF EXISTS transaction_number;