-- Remove financial fields from outcome_invoices table
ALTER TABLE outcome_invoices
DROP COLUMN IF EXISTS discount_amount,
DROP COLUMN IF EXISTS tax_amount,
DROP COLUMN IF EXISTS subtotal,
DROP COLUMN IF EXISTS due_date;
