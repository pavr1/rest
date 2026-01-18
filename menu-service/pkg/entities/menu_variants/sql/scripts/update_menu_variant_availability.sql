UPDATE menu_variants
SET is_available = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, name, description, sub_category_id, price, item_cost, happy_hour_price, image_url,
          is_available, preparation_time, menu_types, dietary_tags, allergens, is_alcoholic,
          display_order, created_at, updated_at;
