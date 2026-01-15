-- Update an outcome invoice
UPDATE outcome_invoices SET
    supplier_id = COALESCE($2, supplier_id),
    transaction_date = COALESCE($3, transaction_date),
    total_amount = COALESCE($4, total_amount),
    image_url = COALESCE($5, image_url),
    notes = COALESCE($6, notes),
    updated_at = NOW()
WHERE id = $1
RETURNING id, updated_at;