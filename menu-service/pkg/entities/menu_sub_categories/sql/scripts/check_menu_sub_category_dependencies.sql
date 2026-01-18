-- Check if sub-category has any menu variants
SELECT COUNT(*) as count FROM menu_variants WHERE sub_category_id = $1;
