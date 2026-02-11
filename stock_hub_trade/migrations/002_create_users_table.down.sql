-- Down migration: Drop users table

DROP INDEX IF EXISTS idx_username;
DROP TABLE IF EXISTS users;
