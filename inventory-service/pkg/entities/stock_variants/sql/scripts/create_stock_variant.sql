INSERT INTO stock_variants (name, stock_sub_category_id, invoice_id, unit, number_of_units, is_active)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, name, stock_sub_category_id, invoice_id, unit, number_of_units, is_active, created_at, updated_at;
