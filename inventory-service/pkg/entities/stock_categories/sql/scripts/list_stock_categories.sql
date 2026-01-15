SELECT id, name, description, display_order, is_active, created_at, updated_at
FROM stock_categories
ORDER BY display_order ASC, name ASC
LIMIT $1 OFFSET $2;
