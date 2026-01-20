-- List invoice items for a specific invoice
SELECT
    id,
    invoice_id,
    inventory_category_id,
    inventory_sub_category_id,
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
ORDER BY created_at DESC;
