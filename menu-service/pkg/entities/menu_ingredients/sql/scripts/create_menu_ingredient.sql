-- Create menu ingredient
INSERT INTO menu_ingredients (
    menu_variant_id,
    stock_variant_id,
    quantity,
    is_optional,
    notes
) VALUES ($1, $2, $3, $4, $5)
RETURNING
    id,
    menu_variant_id,
    stock_variant_id,
    quantity,
    is_optional,
    notes,
    created_at,
    updated_at;
