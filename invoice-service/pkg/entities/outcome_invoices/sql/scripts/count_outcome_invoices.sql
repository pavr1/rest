-- Count outcome invoices with optional filters
SELECT COUNT(*)
FROM outcome_invoices
WHERE ($1 IS NULL OR supplier_id = $1);