-- Get income invoice by ID
SELECT
    id,
    order_id,
    payment_id,
    customer_id,
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
WHERE id = $1;
