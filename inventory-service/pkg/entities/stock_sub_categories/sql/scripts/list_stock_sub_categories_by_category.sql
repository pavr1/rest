SELECT id, name, description, stock_category_id, display_order, is_active, created_at, updated_at
FROM stock_sub_categories
WHERE stock_category_id = $1
ORDER BY display_order ASC, name ASC
LIMIT $2 OFFSET $3;
