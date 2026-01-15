-- Get outcome invoice by ID
SELECT
    id,
    invoice_number,
    supplier_id,
    transaction_date,
    total_amount,
    image_url,
    notes,
    created_at,
    updated_at
FROM outcome_invoices
WHERE id = $1;
