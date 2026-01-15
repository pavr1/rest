-- List invoice items with optional filters
SELECT
    id,
    invoice_id,
    stock_item_id,
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
WHERE ($1 IS NULL OR invoice_id = $1)
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;
