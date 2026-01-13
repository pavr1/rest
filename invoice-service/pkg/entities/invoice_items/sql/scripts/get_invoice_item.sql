-- Get invoice item by ID
SELECT
    id,
    invoice_id,
    invoice_type,
    detail,
    count,
    unit_type,
    price,
    items_per_unit,
    total,
    expiration_date,
    created_at,
    updated_at
FROM invoice_items
WHERE id = $1
