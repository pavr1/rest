UPDATE menu_items
SET is_available = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, name, description, category_id, price, item_cost, happy_hour_price, image_url,
          is_available, item_type, menu_types, dietary_tags, allergens, is_alcoholic,
          created_at, updated_at;
