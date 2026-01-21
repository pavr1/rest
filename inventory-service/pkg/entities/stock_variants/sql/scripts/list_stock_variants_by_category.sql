SELECT sv.id, sv.name, sv.description, sv.stock_sub_category_id, sv.is_active, sv.created_at, sv.updated_at
FROM stock_variants sv
JOIN stock_sub_categories ssc ON sv.stock_sub_category_id = ssc.id
WHERE ssc.stock_category_id = $1
ORDER BY sv.name ASC
LIMIT $2 OFFSET $3;
