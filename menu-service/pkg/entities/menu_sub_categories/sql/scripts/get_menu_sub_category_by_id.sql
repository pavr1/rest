SELECT
    sm.id,
    sm.name,
    sm.description,
    sm.category_id,
    mc.name as category_name,
    sm.item_type,
    sm.display_order,
    sm.is_active,
    sm.created_at,
    sm.updated_at
FROM menu_sub_categories sm
LEFT JOIN menu_categories mc ON sm.category_id = mc.id
WHERE sm.id = $1;
