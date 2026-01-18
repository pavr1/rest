-- Get menu ingredient by ID
SELECT
    mi.id,
    mi.menu_variant_id,
    mi.stock_sub_category_id,
    ssc.name as stock_sub_category_name,
    mi.quantity,
    mi.is_optional,
    mi.notes,
    mi.created_at,
    mi.updated_at
FROM menu_ingredients mi
JOIN stock_sub_categories ssc ON mi.stock_sub_category_id = ssc.id
WHERE mi.id = $1;
