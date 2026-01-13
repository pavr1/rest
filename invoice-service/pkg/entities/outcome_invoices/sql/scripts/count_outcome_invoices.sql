-- Count outcome invoices with optional filters
SELECT COUNT(*)
FROM outcome_invoices
WHERE ($1::text IS NULL OR supplier_id = $1)
