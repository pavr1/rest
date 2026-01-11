-- Rollback: Restore sessions table with staff_id, created_at, expires_at
-- Version: 004

-- Drop token index
DROP INDEX IF EXISTS idx_sessions_token;

-- Add back columns
ALTER TABLE sessions ADD COLUMN IF NOT EXISTS staff_id UUID REFERENCES staff(id) ON DELETE CASCADE;
ALTER TABLE sessions ADD COLUMN IF NOT EXISTS created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE sessions ADD COLUMN IF NOT EXISTS expires_at TIMESTAMP;

-- Recreate indexes
CREATE INDEX IF NOT EXISTS idx_sessions_staff_id ON sessions(staff_id);
CREATE INDEX IF NOT EXISTS idx_sessions_expires_at ON sessions(expires_at);
