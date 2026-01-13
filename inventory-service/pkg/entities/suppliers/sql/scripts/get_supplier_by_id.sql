SELECT id, name, contact_name, phone, email, address, created_at, updated_at
FROM suppliers
WHERE id = $1