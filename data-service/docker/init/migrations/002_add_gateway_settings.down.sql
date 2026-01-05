-- Rollback: Remove gateway service settings
-- Version: 002

DELETE FROM settings WHERE service = 'gateway';

