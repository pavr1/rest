SELECT id, name, description, stock_sub_category_id, is_active, created_at, updated_at
FROM stock_variants
ORDER BY name ASC
LIMIT $1 OFFSET $2;
