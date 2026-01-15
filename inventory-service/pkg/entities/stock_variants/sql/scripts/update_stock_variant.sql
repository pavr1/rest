UPDATE stock_variants
SET name = COALESCE($2, name),
    unit = COALESCE($3, unit),
    number_of_units = COALESCE($4, number_of_units),
    is_active = COALESCE($5, is_active),
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, name, stock_sub_category_id, unit, number_of_units, is_active, created_at, updated_at;
