-- Migration Rollback: Revert invoice schema changes
-- Version: 004
-- Date: 2026-01-14

-- Revert income_invoices table changes
-- Rename customer_id back to customer_tax_id
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'income_invoices' AND column_name = 'customer_id') THEN
        ALTER TABLE income_invoices RENAME COLUMN customer_id TO customer_tax_id;
    END IF;
END $$;

-- Add back customer_id UUID FK column
ALTER TABLE income_invoices ADD COLUMN customer_id UUID REFERENCES customers(id) ON DELETE SET NULL;

-- Add back customer_name column
ALTER TABLE income_invoices ADD COLUMN customer_name VARCHAR(255) NOT NULL DEFAULT '';

-- Remove invoice_type column from invoice_items
ALTER TABLE invoice_items DROP COLUMN IF EXISTS invoice_type;

-- Recreate original indexes
DROP INDEX IF EXISTS idx_income_invoices_customer_id;
CREATE INDEX IF NOT EXISTS idx_income_invoices_customer_id ON income_invoices(customer_id);
