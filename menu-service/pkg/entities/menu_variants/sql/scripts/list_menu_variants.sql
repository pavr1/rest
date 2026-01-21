SELECT mi.id, mi.name, mi.description, mi.sub_category_id, sm.name as sub_category_name,
       sm.item_type, mi.price, mi.item_cost, mi.happy_hour_price, mi.image_url, mi.is_available,
       mi.preparation_time, mi.menu_types, mi.dietary_tags, mi.allergens, mi.is_alcoholic,
       mi.display_order, mi.created_at, mi.updated_at
FROM menu_variants mi
LEFT JOIN menu_sub_categories sm ON mi.sub_category_id = sm.id
WHERE ($1::uuid IS NULL OR sm.category_id = $1)
  AND ($2::uuid IS NULL OR mi.sub_category_id = $2)
  AND ($3::boolean IS NULL OR mi.is_available = $3)
  AND ($4::jsonb IS NULL OR mi.menu_types @> $4)
ORDER BY mi.display_order, mi.name ASC
LIMIT $5 OFFSET $6;
