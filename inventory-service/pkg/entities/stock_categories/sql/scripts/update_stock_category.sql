UPDATE stock_categories
SET name = COALESCE($2, name),
    description = COALESCE($3, description),
    display_order = COALESCE($4, display_order),
    is_active = COALESCE($5, is_active),
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, name, description, display_order, is_active, created_at, updated_at;
