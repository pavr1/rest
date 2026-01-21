UPDATE stock_count
SET count = COALESCE($2, count),
    unit = COALESCE($3, unit),
    unit_price = COALESCE($4, unit_price),
    cost_per_portion = COALESCE($5, cost_per_portion),
    is_out = COALESCE($6, is_out),
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, stock_variant_id, invoice_id, count, unit, unit_price, cost_per_portion, purchased_at, is_out, created_at, updated_at;
