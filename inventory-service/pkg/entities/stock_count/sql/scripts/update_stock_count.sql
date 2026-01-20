UPDATE stock_count
SET count = COALESCE($2, count),
    unit = COALESCE($3, unit),
    is_out = COALESCE($4, is_out),
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, stock_variant_id, invoice_id, count, unit, purchased_at, is_out, created_at, updated_at;
