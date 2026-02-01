SELECT DISTINCT menu_variant_id
FROM menu_ingredients
WHERE stock_variant_id = $1;
