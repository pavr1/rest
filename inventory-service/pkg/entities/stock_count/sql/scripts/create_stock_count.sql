INSERT INTO stock_count (stock_variant_id, invoice_id, count, unit, purchased_at, is_out)
VALUES ($1, $2, $3, $4, $5, false)
RETURNING id, stock_variant_id, invoice_id, count, unit, purchased_at, is_out, created_at, updated_at;
