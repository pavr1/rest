INSERT INTO stock_sub_categories (name, description, stock_category_id, display_order, is_active)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, name, description, stock_category_id, display_order, is_active, created_at, updated_at;
