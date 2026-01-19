-- Get invoice item by ID
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
WHERE id = $1;
