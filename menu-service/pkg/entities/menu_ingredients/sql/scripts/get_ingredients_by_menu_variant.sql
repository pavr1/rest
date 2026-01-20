-- Get ingredients for a specific menu variant
SELECT
    mi.id,
    mi.menu_variant_id,
    mi.stock_variant_id,
    sv.name as stock_variant_name,
    mi.quantity,
    mi.is_optional,
    mi.notes,
    mi.created_at,
    mi.updated_at
FROM menu_ingredients mi
LEFT JOIN stock_variants sv ON mi.stock_variant_id = sv.id
WHERE mi.menu_variant_id = $1
ORDER BY mi.created_at;
