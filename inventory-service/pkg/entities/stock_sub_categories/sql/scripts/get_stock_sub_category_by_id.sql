SELECT id, name, description, stock_category_id, display_order, is_active, created_at, updated_at
FROM stock_sub_categories
WHERE id = $1;
