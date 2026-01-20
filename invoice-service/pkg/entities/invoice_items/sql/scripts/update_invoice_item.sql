-- Update an invoice item
UPDATE invoice_items SET
    inventory_category_id = COALESCE($2, inventory_category_id),
    inventory_sub_category_id = COALESCE($3, inventory_sub_category_id),
    detail = COALESCE($4, detail),
    count = COALESCE($5, count),
    unit_type = COALESCE($6, unit_type),
    price = COALESCE($7, price),
    items_per_unit = COALESCE($8, items_per_unit),
    expiration_date = COALESCE($9, expiration_date),
    updated_at = NOW()
WHERE id = $1
RETURNING id, total, updated_at;
