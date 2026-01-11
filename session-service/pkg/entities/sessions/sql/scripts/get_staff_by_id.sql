SELECT id, username, email, first_name, last_name, role, is_active, last_login_at, created_at, updated_at
FROM staff
WHERE id = $1 AND is_active = true
