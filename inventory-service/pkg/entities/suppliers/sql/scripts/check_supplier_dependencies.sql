SELECT
    (SELECT COUNT(*) FROM purchase_invoices WHERE supplier_id = $1) as purchase_invoice_count,
    (SELECT COUNT(*) FROM outcome_invoices WHERE supplier_id = $1) as outcome_invoice_count