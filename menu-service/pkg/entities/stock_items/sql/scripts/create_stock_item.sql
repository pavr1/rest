INSERT INTO stock_items (name, unit, description, stock_item_category_id)
VALUES ($1, $2, $3, $4)
RETURNING id, name, unit, description, stock_item_category_id, created_at, updated_at;
