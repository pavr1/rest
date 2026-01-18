UPDATE menu_categories
SET name = COALESCE($2, name),
    display_order = COALESCE($3, display_order),
    description = COALESCE($4, description),
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, name, display_order, description, created_at, updated_at;
