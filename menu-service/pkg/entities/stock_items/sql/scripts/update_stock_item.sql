UPDATE stock_items
SET name = COALESCE($2, name),
    unit = COALESCE($3, unit),
    description = COALESCE($4, description),
    stock_item_category_id = COALESCE($5, stock_item_category_id),
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, name, unit, description, stock_item_category_id, created_at, updated_at;
