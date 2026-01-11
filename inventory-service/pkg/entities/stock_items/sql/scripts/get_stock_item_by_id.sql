SELECT si.id, si.name, si.unit, si.description, si.stock_item_category_id,
       sic.name as category_name, si.created_at, si.updated_at
FROM stock_items si
LEFT JOIN stock_item_categories sic ON si.stock_item_category_id = sic.id
WHERE si.id = $1;
