SELECT id, name, contact_name, phone, email, address, created_at, updated_at
FROM suppliers
WHERE ($1::text IS NULL OR name ILIKE '%' || $1 || '%')
  AND ($2::text IS NULL OR email ILIKE '%' || $2 || '%')
  AND ($3::text IS NULL OR phone ILIKE '%' || $3 || '%')
ORDER BY name
LIMIT $4 OFFSET $5