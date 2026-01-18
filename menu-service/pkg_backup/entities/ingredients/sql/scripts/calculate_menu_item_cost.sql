-- Calculate the cost of a menu item based on its ingredients
-- Uses the average unit price from available existences for each stock item
SELECT 
    mis.stock_item_id,
    si.name as stock_item_name,
    mis.quantity,
    COALESCE(
        (SELECT AVG(id.price / id.items_per_unit)
         FROM existences e
         JOIN invoice_details id ON e.invoice_detail_id = id.id
         WHERE e.stock_item_id = mis.stock_item_id
           AND e.units_available > 0
        ), 0
    ) as unit_cost
FROM menu_item_stock_items mis
JOIN stock_items si ON mis.stock_item_id = si.id
WHERE mis.menu_item_id = $1;
