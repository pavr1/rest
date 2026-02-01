-- Delete menu ingredient (returns menu_variant_id for cost recalc)
DELETE FROM menu_ingredients
WHERE id = $1
RETURNING menu_variant_id;
