-- Update an outcome invoice
UPDATE outcome_invoices SET
    supplier_id = COALESCE($2, supplier_id),
    inventory_category_id = COALESCE($3, inventory_category_id),
    inventory_sub_category_id = COALESCE($4, inventory_sub_category_id),
    transaction_date = COALESCE($5, transaction_date),
    due_date = COALESCE($6, due_date),
    subtotal = COALESCE($7, subtotal),
    tax_amount = COALESCE($8, tax_amount),
    discount_amount = COALESCE($9, discount_amount),
    total_amount = COALESCE($10, total_amount),
    image_url = COALESCE($11, image_url),
    notes = COALESCE($12, notes),
    updated_at = NOW()
WHERE id = $1
RETURNING id, updated_at;
