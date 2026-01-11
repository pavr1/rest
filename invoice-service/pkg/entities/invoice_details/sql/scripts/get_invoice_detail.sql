SELECT id, invoice_id, stock_item_id, description, quantity,
       unit_of_measure, items_per_unit, unit_price, total_price,
       expiry_date, batch_number, created_at, updated_at
FROM invoice_details WHERE id = $1 AND invoice_id = $2
