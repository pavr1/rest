INSERT INTO sub_menus (name, description, category_id, image_url, item_type, display_order, is_active)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, name, description, category_id, image_url, item_type, display_order, is_active, created_at, updated_at;
