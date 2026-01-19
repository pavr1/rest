-- Create a new invoice item
INSERT INTO invoice_items (
    invoice_id,
    inventory_category_id,
    inventory_sub_category_id,
    detail,
    count,
    unit_type,
    price,
    items_per_unit,
    expiration_date
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING id, total, created_at, updated_at;
