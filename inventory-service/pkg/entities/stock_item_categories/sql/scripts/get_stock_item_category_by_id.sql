SELECT id, name, description, created_at, updated_at
FROM stock_item_categories
WHERE id = $1;
