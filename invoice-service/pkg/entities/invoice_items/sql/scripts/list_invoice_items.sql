-- List invoice items by invoice ID
SELECT
    id,
    invoice_id,
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
WHERE invoice_id = $1
ORDER BY created_at
