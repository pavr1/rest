INSERT INTO menu_categories (name, display_order, description)
VALUES ($1, $2, $3)
RETURNING id, name, display_order, description, created_at, updated_at;
