-- Get ingredients for a specific menu variant
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
WHERE mi.menu_variant_id = $1
ORDER BY mi.created_at;
