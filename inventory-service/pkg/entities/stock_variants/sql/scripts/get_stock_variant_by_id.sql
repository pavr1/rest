SELECT id, name, stock_sub_category_id, invoice_id, unit, number_of_units, is_active, created_at, updated_at
FROM stock_variants
WHERE id = $1;
