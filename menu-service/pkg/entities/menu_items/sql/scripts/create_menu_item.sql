INSERT INTO menu_items (name, description, category_id, price, happy_hour_price, image_url,
                        is_available, item_type, menu_types, dietary_tags, allergens, is_alcoholic)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
RETURNING id, name, description, category_id, price, item_cost, happy_hour_price, image_url,
          is_available, item_type, menu_types, dietary_tags, allergens, is_alcoholic,
          created_at, updated_at;
