UPDATE stock_count
SET is_out = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, stock_variant_id, invoice_id, count, unit, purchased_at, is_out, created_at, updated_at;
