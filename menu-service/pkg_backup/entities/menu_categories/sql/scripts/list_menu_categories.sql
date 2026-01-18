SELECT id, name, display_order, description, created_at, updated_at
FROM menu_categories
ORDER BY display_order ASC
LIMIT $1 OFFSET $2;
