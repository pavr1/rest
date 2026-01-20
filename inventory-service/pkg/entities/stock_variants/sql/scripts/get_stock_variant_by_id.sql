SELECT id, name, description, stock_sub_category_id, is_active, created_at, updated_at
FROM stock_variants
WHERE id = $1;
