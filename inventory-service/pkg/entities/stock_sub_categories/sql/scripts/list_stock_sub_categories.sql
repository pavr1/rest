SELECT id, name, description, stock_category_id, display_order, is_active, created_at, updated_at
FROM stock_sub_categories
ORDER BY display_order ASC, name ASC
LIMIT $1 OFFSET $2;
