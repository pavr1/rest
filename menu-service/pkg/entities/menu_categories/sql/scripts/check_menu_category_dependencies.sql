SELECT COUNT(*) FROM menu_variants mv
JOIN menu_sub_categories msc ON mv.sub_category_id = msc.id
WHERE msc.category_id = $1;
