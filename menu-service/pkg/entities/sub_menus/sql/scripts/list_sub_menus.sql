SELECT 
    sm.id,
    sm.name,
    sm.description,
    sm.category_id,
    mc.name as category_name,
    sm.image_url,
    sm.item_type,
    sm.display_order,
    sm.is_active,
    sm.created_at,
    sm.updated_at
FROM sub_menus sm
LEFT JOIN menu_categories mc ON sm.category_id = mc.id
WHERE 
    ($1::uuid IS NULL OR sm.category_id = $1)
    AND ($2::varchar IS NULL OR sm.item_type = $2)
    AND ($3::boolean IS NULL OR sm.is_active = $3)
ORDER BY sm.display_order, sm.name
LIMIT $4 OFFSET $5;
