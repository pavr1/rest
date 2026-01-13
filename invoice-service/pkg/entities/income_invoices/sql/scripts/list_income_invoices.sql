-- List income invoices with optional filters
SELECT
    id,
    order_id,
    payment_id,
    customer_id,
    customer_name,
    customer_tax_id,
    invoice_number,
    invoice_type,
    subtotal,
    tax_amount,
    service_charge,
    total_amount,
    payment_method,
    xml_data,
    digital_signature,
    status,
    generated_at,
    created_at,
    updated_at
FROM income_invoices
WHERE ($1::text IS NULL OR customer_name ILIKE '%' || $1 || '%')
  AND ($2::text IS NULL OR invoice_type = $2)
  AND ($3::text IS NULL OR status = $3)
  AND ($4::text IS NULL OR order_id = $4)
  AND ($5::text IS NULL OR customer_id = $5)
ORDER BY created_at DESC
LIMIT $6 OFFSET $7;