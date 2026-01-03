-- Migration: Add/Update default admin user
-- Version: 003
-- Date: 2024
-- Password: admin (bcrypt hashed)

-- Update existing admin user password or insert if not exists
INSERT INTO staff (username, email, password_hash, first_name, last_name, role)
VALUES (
    'admin',
    'pavr1@hotmail.com',
    '$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrqZ6YQy6lL.PV0YBbK1QK5aQKz/gK',
    'System',
    'Administrator',
    'admin'
)
ON CONFLICT (username) DO UPDATE SET
    password_hash = EXCLUDED.password_hash,
    updated_at = CURRENT_TIMESTAMP;
