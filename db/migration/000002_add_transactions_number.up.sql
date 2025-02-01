-- Add the uuid-ossp extension if not already added
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";


-- Add transaction_number column
ALTER TABLE transactions
ADD COLUMN transaction_number UUID DEFAULT uuid_generate_v4() NOT NULL UNIQUE;