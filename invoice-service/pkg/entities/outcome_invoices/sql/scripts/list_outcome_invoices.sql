-- List outcome invoices
-- pvillalobos -> revisit later about adding NULL suppliers for filtering
SELECT
    id,
    invoice_number,
    supplier_id,
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
FROM outcome_invoices
ORDER BY transaction_date DESC
LIMIT $1 OFFSET $2;
