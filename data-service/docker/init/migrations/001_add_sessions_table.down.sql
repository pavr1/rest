-- Rollback: Remove sessions table
-- Version: 001

-- Remove session service settings
DELETE FROM settings WHERE service = 'session';

-- Drop indexes
DROP INDEX IF EXISTS idx_sessions_expires_at;
DROP INDEX IF EXISTS idx_sessions_staff_id;

-- Drop sessions table
DROP TABLE IF EXISTS sessions;
