SELECT id, name, stock_sub_category_id, unit, number_of_units, is_active, created_at, updated_at
FROM stock_variants
ORDER BY name ASC
LIMIT $1 OFFSET $2;
