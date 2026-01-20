-- Create a stock count record when invoice item is created
INSERT INTO stock_count (stock_variant_id, invoice_id, count, unit, purchased_at, is_out)
VALUES ($1, $2, $3, $4, $5, false);
