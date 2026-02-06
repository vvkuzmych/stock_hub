-- Down migration: Drop messages table

DROP INDEX IF EXISTS idx_client_id;
DROP INDEX IF EXISTS idx_created_at;
DROP TABLE IF EXISTS messages;
