UPDATE invoice_details SET
    stock_item_id = COALESCE($1, stock_item_id),
    description = COALESCE($2, description),
    quantity = COALESCE($3, quantity),
    unit_of_measure = COALESCE($4, unit_of_measure),
    items_per_unit = COALESCE($5, items_per_unit),
    unit_price = COALESCE($6, unit_price),
    total_price = COALESCE($7, total_price),
    expiry_date = COALESCE($8, expiry_date),
    batch_number = COALESCE($9, batch_number),
    updated_at = NOW()
WHERE id = $10 AND invoice_id = $11
