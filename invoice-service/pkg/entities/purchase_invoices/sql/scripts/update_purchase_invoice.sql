UPDATE purchase_invoices SET
    supplier_name = COALESCE($1, supplier_name),
    invoice_date = COALESCE($2, invoice_date),
    due_date = COALESCE($3, due_date),
    total_amount = COALESCE($4, total_amount),
    status = COALESCE($5, status),
    image_url = COALESCE($6, image_url),
    notes = COALESCE($7, notes),
    updated_at = NOW()
WHERE id = $8
