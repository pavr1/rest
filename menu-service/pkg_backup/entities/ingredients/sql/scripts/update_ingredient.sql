UPDATE menu_item_stock_items
SET quantity = $3, updated_at = CURRENT_TIMESTAMP
WHERE menu_item_id = $1 AND stock_item_id = $2
RETURNING id, menu_item_id, stock_item_id, quantity, created_at, updated_at;
