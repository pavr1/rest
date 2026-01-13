-- List outcome invoices with optional filters
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
WHERE ($1::text IS NULL OR supplier_id = $1)
ORDER BY transaction_date DESC
LIMIT $2 OFFSET $3
