-- Recalculate and set menu_variant.item_cost from sum of (quantity * stock_variants.avg_cost) for its ingredients
UPDATE menu_variants
SET item_cost = COALESCE((
    SELECT SUM(mi.quantity * COALESCE(sv.avg_cost, 0))
    FROM menu_ingredients mi
    LEFT JOIN stock_variants sv ON mi.stock_variant_id = sv.id
    WHERE mi.menu_variant_id = $1
), 0),
updated_at = CURRENT_TIMESTAMP
WHERE id = $1;
