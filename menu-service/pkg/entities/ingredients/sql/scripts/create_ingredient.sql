INSERT INTO menu_item_stock_items (menu_item_id, stock_item_id, quantity)
VALUES ($1, $2, $3)
RETURNING id, menu_item_id, stock_item_id, quantity, created_at, updated_at;
