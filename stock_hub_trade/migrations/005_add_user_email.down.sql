-- Remove email field from users table

DROP INDEX IF EXISTS idx_user_email;
ALTER TABLE users DROP COLUMN IF EXISTS email;
