-- Create menu ingredient
INSERT INTO menu_ingredients (
    menu_variant_id,
    stock_variant_id,
    menu_sub_category_id,
    quantity,
    is_optional,
    notes
) VALUES ($1, $2, $3, $4, $5, $6)
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
