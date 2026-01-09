SELECT mi.id, mi.name, mi.description, mi.category_id, mc.name as category_name,
       mi.price, mi.item_cost, mi.happy_hour_price, mi.image_url, mi.is_available,
       mi.item_type, mi.menu_types, mi.dietary_tags, mi.allergens, mi.is_alcoholic,
       mi.created_at, mi.updated_at
FROM menu_items mi
LEFT JOIN menu_categories mc ON mi.category_id = mc.id
WHERE mi.id = $1;
