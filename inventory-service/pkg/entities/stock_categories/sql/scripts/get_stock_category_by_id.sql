SELECT id, name, description, display_order, is_active, created_at, updated_at
FROM stock_categories
WHERE id = $1;
