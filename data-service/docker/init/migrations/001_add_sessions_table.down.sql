-- Rollback: Remove sessions table
-- Version: 001

-- Remove session service settings
DELETE FROM settings WHERE service = 'session';

-- Drop indexes
DROP INDEX IF EXISTS idx_sessions_expires_at;
DROP INDEX IF EXISTS idx_sessions_staff_id;

-- Remove added columns (but keep the table structure minimal)
ALTER TABLE sessions DROP COLUMN IF EXISTS expires_at;
ALTER TABLE sessions DROP COLUMN IF EXISTS created_at;
ALTER TABLE sessions DROP COLUMN IF EXISTS staff_id;
