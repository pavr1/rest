UPDATE sub_menus
SET 
    name = COALESCE($2, name),
    description = COALESCE($3, description),
    category_id = COALESCE($4, category_id),
    image_url = COALESCE($5, image_url),
    item_type = COALESCE($6, item_type),
    display_order = COALESCE($7, display_order),
    is_active = COALESCE($8, is_active)
WHERE id = $1
RETURNING id, name, description, category_id, image_url, item_type, display_order, is_active, created_at, updated_at;
