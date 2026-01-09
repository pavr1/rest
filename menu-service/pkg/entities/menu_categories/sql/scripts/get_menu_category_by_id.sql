SELECT id, name, display_order, description, created_at, updated_at
FROM menu_categories
WHERE id = $1;
