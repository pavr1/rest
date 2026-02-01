-- Recalculate menu_sub_categories.item_cost as AVG(menu_variants.item_cost) for the sub_category that owns the given menu_variant_id
UPDATE menu_sub_categories
SET item_cost = COALESCE((
    SELECT AVG(mv.item_cost)
    FROM menu_variants mv
    WHERE mv.sub_category_id = menu_sub_categories.id
      AND mv.item_cost IS NOT NULL
), 0),
updated_at = CURRENT_TIMESTAMP
WHERE id = (SELECT sub_category_id FROM menu_variants WHERE id = $1);
