-- Rollback: Remove sessions table
-- Version: 001

-- Remove session service settings
DELETE FROM settings WHERE service = 'session';

-- Rollback: Remove sessions table
DROP TABLE IF EXISTS sessions;
