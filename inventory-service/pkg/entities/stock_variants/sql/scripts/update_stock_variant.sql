UPDATE stock_variants
SET name = COALESCE($2, name),
    description = COALESCE($3, description),
    is_active = COALESCE($4, is_active),
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, name, description, stock_sub_category_id, avg_cost, is_active, created_at, updated_at;
