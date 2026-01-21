SELECT COUNT(*) 
FROM stock_variants sv
JOIN stock_sub_categories ssc ON sv.stock_sub_category_id = ssc.id
WHERE ssc.stock_category_id = $1;
