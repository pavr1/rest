-- Update an invoice item
UPDATE invoice_items SET
    detail = COALESCE($1, detail),
    count = COALESCE($2, count),
    unit_type = COALESCE($3, unit_type),
    price = COALESCE($4, price),
    items_per_unit = COALESCE($5, items_per_unit),
    total = COALESCE($6, total),
    expiration_date = COALESCE($7, expiration_date),
    updated_at = NOW()
WHERE id = $8
