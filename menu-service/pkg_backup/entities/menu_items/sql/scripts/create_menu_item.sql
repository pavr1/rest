INSERT INTO menu_items (name, description, sub_menu_id, price, happy_hour_price, image_url,
                        is_available, preparation_time, menu_types, dietary_tags, allergens, is_alcoholic, display_order)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
RETURNING id, name, description, sub_menu_id, price, item_cost, happy_hour_price, image_url,
          is_available, preparation_time, menu_types, dietary_tags, allergens, is_alcoholic,
          display_order, created_at, updated_at;
