-- Count invoice items for a specific invoice
SELECT COUNT(*)
FROM invoice_items
WHERE invoice_id = $1;
