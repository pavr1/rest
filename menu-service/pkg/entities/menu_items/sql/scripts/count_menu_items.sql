SELECT COUNT(*)
FROM menu_items mi
WHERE ($1::uuid IS NULL OR mi.category_id = $1)
  AND ($2::boolean IS NULL OR mi.is_available = $2)
  AND ($3::jsonb IS NULL OR mi.menu_types = $3);
