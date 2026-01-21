SELECT COUNT(*)
FROM menu_variants mi
LEFT JOIN menu_sub_categories sm ON mi.sub_category_id = sm.id
WHERE ($1::uuid IS NULL OR sm.category_id = $1)
  AND ($2::uuid IS NULL OR mi.sub_category_id = $2)
  AND ($3::boolean IS NULL OR mi.is_available = $3)
  AND ($4::jsonb IS NULL OR mi.menu_types @> $4);
