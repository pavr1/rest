UPDATE menu_items
SET name = COALESCE($2, name),
    description = COALESCE($3, description),
    sub_menu_id = COALESCE($4, sub_menu_id),
    price = COALESCE($5, price),
    happy_hour_price = COALESCE($6, happy_hour_price),
    image_url = COALESCE($7, image_url),
    is_available = COALESCE($8, is_available),
    preparation_time = COALESCE($9, preparation_time),
    menu_types = COALESCE($10, menu_types),
    dietary_tags = COALESCE($11, dietary_tags),
    allergens = COALESCE($12, allergens),
    is_alcoholic = COALESCE($13, is_alcoholic),
    display_order = COALESCE($14, display_order),
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, name, description, sub_menu_id, price, item_cost, happy_hour_price, image_url,
          is_available, preparation_time, menu_types, dietary_tags, allergens, is_alcoholic,
          display_order, created_at, updated_at;
