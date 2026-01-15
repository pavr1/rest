-- Count invoice items with optional filters
SELECT COUNT(*)
FROM invoice_items
WHERE ($1 IS NULL OR invoice_id = $1);
