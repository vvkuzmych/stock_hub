-- Add email field to users table

ALTER TABLE users ADD COLUMN IF NOT EXISTS email TEXT;

-- Create index on email for faster lookups
CREATE INDEX IF NOT EXISTS idx_user_email ON users(email);

-- Add unique constraint (optional, uncomment if emails should be unique)
-- ALTER TABLE users ADD CONSTRAINT users_email_key UNIQUE (email);
