INSERT INTO stock_variants (name, description, stock_sub_category_id, is_active)
VALUES ($1, $2, $3, $4)
RETURNING id, name, description, stock_sub_category_id, avg_cost, is_active, created_at, updated_at;
