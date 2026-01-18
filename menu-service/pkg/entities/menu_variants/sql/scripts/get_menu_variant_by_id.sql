SELECT mi.id, mi.name, mi.description, mi.sub_category_id, sm.name as sub_menu_name,
       sm.item_type, mi.price, mi.item_cost, mi.happy_hour_price, mi.image_url, mi.is_available,
       mi.preparation_time, mi.menu_types, mi.dietary_tags, mi.allergens, mi.is_alcoholic,
       mi.display_order, mi.created_at, mi.updated_at
FROM menu_variants mi
LEFT JOIN sub_menus sm ON mi.sub_category_id = sm.id
WHERE mi.id = $1;
