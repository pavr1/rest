UPDATE menu_sub_categories
SET
    name = COALESCE($2, name),
    description = COALESCE($3, description),
    category_id = COALESCE($4, category_id),
    item_type = COALESCE($5, item_type),
    display_order = COALESCE($6, display_order),
    is_active = COALESCE($7, is_active)
WHERE id = $1
RETURNING id, name, description, category_id, item_type, display_order, is_active, created_at, updated_at;
