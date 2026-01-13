INSERT INTO suppliers (name, contact_name, phone, email, address)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, name, contact_name, phone, email, address, created_at, updated_at