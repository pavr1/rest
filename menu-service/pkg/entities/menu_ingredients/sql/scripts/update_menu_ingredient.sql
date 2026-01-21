-- Update menu ingredient
UPDATE menu_ingredients
SET
    quantity = COALESCE($2, quantity),
    is_optional = COALESCE($3, is_optional),
    notes = COALESCE($4, notes),
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING
    id,
    menu_variant_id,
    stock_variant_id,
    menu_sub_category_id,
    quantity,
    is_optional,
    notes,
    created_at,
    updated_at;
