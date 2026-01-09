SELECT COUNT(*)
FROM stock_items si
WHERE ($1::uuid IS NULL OR si.stock_item_category_id = $1);
