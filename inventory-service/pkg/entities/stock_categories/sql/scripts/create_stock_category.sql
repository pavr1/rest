INSERT INTO stock_categories (name, description, display_order, is_active)
VALUES ($1, $2, $3, $4)
RETURNING id, name, description, display_order, is_active, created_at, updated_at;
