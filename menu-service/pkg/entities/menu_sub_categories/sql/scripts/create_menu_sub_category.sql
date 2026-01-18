INSERT INTO menu_sub_categories (name, description, category_id, item_type, display_order, is_active)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, name, description, category_id, item_type, display_order, is_active, created_at, updated_at;
