-- Update an outcome invoice
UPDATE outcome_invoices SET
    supplier_id = COALESCE($2, supplier_id),
    transaction_date = COALESCE($3, transaction_date),
    due_date = COALESCE($4, due_date),
    subtotal = COALESCE($5, subtotal),
    tax_amount = COALESCE($6, tax_amount),
    discount_amount = COALESCE($7, discount_amount),
    total_amount = COALESCE($8, total_amount),
    image_url = COALESCE($9, image_url),
    notes = COALESCE($10, notes),
    updated_at = NOW()
WHERE id = $1
RETURNING id, updated_at;
