SELECT id, name, description, stock_sub_category_id, avg_cost, is_active, created_at, updated_at
FROM stock_variants
WHERE is_active = true
ORDER BY name ASC;
