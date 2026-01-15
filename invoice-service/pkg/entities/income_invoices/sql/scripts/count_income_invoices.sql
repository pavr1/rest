-- Count income invoices with optional filters
SELECT COUNT(*)
FROM income_invoices
WHERE ($1::text IS NULL OR customer_id ILIKE '%' || $1 || '%')
  AND ($2::text IS NULL OR invoice_type = $2)
  AND ($3::text IS NULL OR status = $3)
  AND ($4::text IS NULL OR order_id = $4);
