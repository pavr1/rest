-- Migration: Simplify sessions table to only store session_id and token
-- Version: 004
-- Date: 2025-01-10
-- Description: Remove staff_id, created_at, expires_at columns since this info is in the JWT token

-- Drop indexes that reference columns we're removing
DROP INDEX IF EXISTS idx_sessions_staff_id;
DROP INDEX IF EXISTS idx_sessions_expires_at;

-- Drop columns that are redundant (info is in JWT token)
ALTER TABLE sessions DROP COLUMN IF EXISTS staff_id;
ALTER TABLE sessions DROP COLUMN IF EXISTS created_at;
ALTER TABLE sessions DROP COLUMN IF EXISTS expires_at;

-- Add index on token for quick lookups
CREATE INDEX IF NOT EXISTS idx_sessions_token ON sessions(token);
