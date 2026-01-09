-- Check if sub menu has any menu items
SELECT COUNT(*) as count FROM menu_items WHERE sub_menu_id = $1;
