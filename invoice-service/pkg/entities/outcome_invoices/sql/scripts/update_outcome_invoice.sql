-- Update an outcome invoice
UPDATE outcome_invoices SET
    supplier_id = COALESCE($2, supplier_id),
    inventory_category_id = COALESCE($3, inventory_category_id),
    inventory_sub_category_id = COALESCE($4, inventory_sub_category_id),
    transaction_date = COALESCE($5, transaction_date),
    total_amount = COALESCE($6, total_amount),
    image_url = COALESCE($7, image_url),
    notes = COALESCE($8, notes),
    updated_at = NOW()
WHERE id = $1
RETURNING id, updated_at;
