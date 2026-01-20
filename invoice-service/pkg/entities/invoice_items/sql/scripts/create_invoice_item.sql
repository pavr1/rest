-- Create a new invoice item
INSERT INTO invoice_items (
    invoice_id,
    stock_variant_id,
    detail,
    count,
    unit_type,
    price,
    items_per_unit,
    expiration_date
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING id, invoice_id, stock_variant_id, detail, count, unit_type, price, items_per_unit, total, expiration_date, created_at, updated_at;
