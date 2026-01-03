-- Rollback: Revert admin password to original (admin123)
-- Version: 003

UPDATE staff 
SET password_hash = '$2a$12$04xNgahgyY9qDgv7goYUVenjgTHF7.ei9GVkp.uYixLs.ebrJxw6u',
    updated_at = CURRENT_TIMESTAMP
WHERE username = 'admin';
