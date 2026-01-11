INSERT INTO invoice_details (
    invoice_id, stock_item_id, description, quantity, unit_of_measure,
    items_per_unit, unit_price, total_price, expiry_date, batch_number
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING id, invoice_id, stock_item_id, description, quantity,
          unit_of_measure, items_per_unit, unit_price, total_price,
          expiry_date, batch_number, created_at, updated_at
