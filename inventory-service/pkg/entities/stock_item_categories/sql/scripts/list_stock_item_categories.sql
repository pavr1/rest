SELECT id, name, description, created_at, updated_at
FROM stock_item_categories
ORDER BY name ASC
LIMIT $1 OFFSET $2;
