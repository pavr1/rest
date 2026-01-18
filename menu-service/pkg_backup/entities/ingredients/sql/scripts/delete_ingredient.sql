DELETE FROM menu_item_stock_items
WHERE menu_item_id = $1 AND stock_item_id = $2;
