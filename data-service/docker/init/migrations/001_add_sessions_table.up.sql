-- Migration: Add sessions table for session management
-- Version: 001
-- Date: 2024

-- Sessions table for storing active user sessions
CREATE TABLE sessions (
    session_id VARCHAR(255) PRIMARY KEY,
    token TEXT NOT NULL,
    staff_id UUID NOT NULL REFERENCES staff(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL
);

-- Index for quick session lookup by staff
CREATE INDEX idx_sessions_staff_id ON sessions(staff_id);

-- Index for cleaning up expired sessions
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);

-- Add session service settings
INSERT INTO settings (service, key, value, description) VALUES
    ('session', 'JWT_SECRET', 'barrest-super-secret-key-change-in-production-2024', 'JWT signing secret key'),
    ('session', 'JWT_EXPIRATION_TIME', '24h', 'JWT token expiration time'),
    ('session', 'SERVER_HOST', '0.0.0.0', 'Session service host'),
    ('session', 'SERVER_PORT', '8087', 'Session service port'),
    ('session', 'LOG_LEVEL', 'INFO', 'Logging level'),
    ('session', 'DB_HOST', 'barrest_postgres', 'Database host'),
    ('session', 'DB_PORT', '5432', 'Database port'),
    ('session', 'DB_USER', 'postgres', 'Database user'),
    ('session', 'DB_PASSWORD', 'postgres123', 'Database password'),
    ('session', 'DB_NAME', 'barrest_db', 'Database name'),
    ('session', 'DB_SSL_MODE', 'disable', 'Database SSL mode')
ON CONFLICT (service, key) DO NOTHING;
