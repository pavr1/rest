UPDATE suppliers
SET name = COALESCE($2, name),
    contact_name = COALESCE($3, contact_name),
    phone = COALESCE($4, phone),
    email = COALESCE($5, email),
    address = COALESCE($6, address),
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, name, contact_name, phone, email, address, created_at, updated_at