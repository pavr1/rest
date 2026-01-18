SELECT COUNT(*) 
FROM menu_sub_categories
WHERE 
    ($1::uuid IS NULL OR category_id = $1)
    AND ($2::varchar IS NULL OR item_type = $2)
    AND ($3::boolean IS NULL OR is_active = $3);
