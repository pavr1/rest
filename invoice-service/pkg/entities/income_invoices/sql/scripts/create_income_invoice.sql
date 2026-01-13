-- Create a new income invoice
INSERT INTO income_invoices (
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
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, NOW(), NOW()
) RETURNING id, created_at, updated_at;