SELECT mis.id, mis.menu_item_id, mis.stock_item_id, si.name as stock_item_name,
       si.unit as stock_item_unit, mis.quantity, mis.created_at, mis.updated_at
FROM menu_item_stock_items mis
JOIN stock_items si ON mis.stock_item_id = si.id
WHERE mis.menu_item_id = $1
ORDER BY si.name ASC;
