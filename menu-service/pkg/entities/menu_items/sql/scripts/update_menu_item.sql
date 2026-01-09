UPDATE menu_items
SET name = COALESCE($2, name),
    description = COALESCE($3, description),
    category_id = COALESCE($4, category_id),
    price = COALESCE($5, price),
    happy_hour_price = COALESCE($6, happy_hour_price),
    image_url = COALESCE($7, image_url),
    is_available = COALESCE($8, is_available),
    item_type = COALESCE($9, item_type),
    menu_types = COALESCE($10, menu_types),
    dietary_tags = COALESCE($11, dietary_tags),
    allergens = COALESCE($12, allergens),
    is_alcoholic = COALESCE($13, is_alcoholic),
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, name, description, category_id, price, item_cost, happy_hour_price, image_url,
          is_available, item_type, menu_types, dietary_tags, allergens, is_alcoholic,
          created_at, updated_at;
