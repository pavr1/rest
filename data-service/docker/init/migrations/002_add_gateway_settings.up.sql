-- Migration: Add gateway service settings
-- Version: 002
-- Date: 2024

-- Add gateway service settings
INSERT INTO settings (service, key, value, description) VALUES
    ('gateway', 'SERVER_HOST', '0.0.0.0', 'Gateway service host'),
    ('gateway', 'SERVER_PORT', '8082', 'Gateway service port'),
    ('gateway', 'LOG_LEVEL', 'INFO', 'Logging level'),
    ('gateway', 'SESSION_SERVICE_URL', 'http://barrest_session_service:8087', 'Session service URL'),
    ('gateway', 'DATA_SERVICE_URL', 'http://barrest_data_service:8086', 'Data service URL'),
    ('gateway', 'CORS_ALLOWED_ORIGINS', '*', 'CORS allowed origins'),
    ('gateway', 'CORS_ALLOWED_METHODS', 'GET,POST,PUT,DELETE,OPTIONS', 'CORS allowed methods'),
    ('gateway', 'CORS_ALLOWED_HEADERS', 'Content-Type,Authorization', 'CORS allowed headers')
ON CONFLICT (service, key) DO NOTHING;
