-- Create a new outcome invoice
INSERT INTO outcome_invoices (
    invoice_number,
    supplier_id,
    inventory_category_id,
    inventory_sub_category_id,
    transaction_date,
    due_date,
    subtotal,
    tax_amount,
    discount_amount,
    total_amount,
    image_url,
    notes,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, NOW(), NOW()
) RETURNING id, created_at, updated_at;
