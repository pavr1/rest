-- Migration: Add sessions table for session management
-- Version: 001
-- Date: 2024

-- Sessions table for storing active user sessions
-- Keep it minimal - session_id and token only
-- All user/auth info is stored in the JWT token itself
CREATE TABLE IF NOT EXISTS sessions (
    session_id VARCHAR(255) PRIMARY KEY,
    token TEXT NOT NULL
);

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
