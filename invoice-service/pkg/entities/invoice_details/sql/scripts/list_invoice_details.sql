SELECT id, invoice_id, stock_item_id, si.name as stock_item_name,
       description, quantity, unit_of_measure, items_per_unit,
       unit_price, total_price, expiry_date, batch_number,
       id.created_at, id.updated_at
FROM invoice_details id
LEFT JOIN stock_items si ON id.stock_item_id = si.id
WHERE id.invoice_id = $1
ORDER BY id.created_at
