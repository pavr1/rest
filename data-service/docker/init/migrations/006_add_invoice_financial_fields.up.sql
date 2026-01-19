-- Add financial fields to outcome_invoices table
ALTER TABLE outcome_invoices
ADD COLUMN IF NOT EXISTS due_date DATE,
ADD COLUMN IF NOT EXISTS subtotal DECIMAL(12,2),
ADD COLUMN IF NOT EXISTS tax_amount DECIMAL(12,2),
ADD COLUMN IF NOT EXISTS discount_amount DECIMAL(12,2) DEFAULT 0;
