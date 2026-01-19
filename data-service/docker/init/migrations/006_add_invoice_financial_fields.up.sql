-- Remove incorrect fields and add financial/stock variant fields
-- Remove inventory fields from outcome_invoices (belong in invoice_items)
ALTER TABLE outcome_invoices
DROP COLUMN IF EXISTS inventory_category_id,
DROP COLUMN IF EXISTS inventory_sub_category_id;

-- Remove redundant columns from invoice_items
ALTER TABLE invoice_items
DROP COLUMN IF EXISTS stock_variant_id,
DROP COLUMN IF EXISTS invoice_type;

-- Add invoice_id to stock_variants (link to purchase invoice)
ALTER TABLE stock_variants
ADD COLUMN IF NOT EXISTS invoice_id UUID REFERENCES outcome_invoices(id) ON DELETE SET NULL;

-- Add financial fields to outcome_invoices
ALTER TABLE outcome_invoices
ADD COLUMN IF NOT EXISTS due_date DATE,
ADD COLUMN IF NOT EXISTS subtotal DECIMAL(12,2),
ADD COLUMN IF NOT EXISTS tax_amount DECIMAL(12,2),
ADD COLUMN IF NOT EXISTS discount_amount DECIMAL(12,2) DEFAULT 0;
