-- Update an invoice item
UPDATE invoice_items SET
    stock_variant_id = COALESCE($2, stock_variant_id),
    detail = COALESCE($3, detail),
    count = COALESCE($4, count),
    unit_type = COALESCE($5, unit_type),
    price = COALESCE($6, price),
    items_per_unit = COALESCE($7, items_per_unit),
    expiration_date = COALESCE($8, expiration_date),
    updated_at = NOW()
WHERE id = $1
RETURNING id, invoice_id, stock_variant_id, detail, count, unit_type, price, items_per_unit, total, expiration_date, created_at, updated_at;
