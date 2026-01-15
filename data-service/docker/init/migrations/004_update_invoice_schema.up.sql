-- Migration: Update invoice schema for new structure
-- Version: 004
-- Date: 2026-01-14
--
-- Changes:
-- 1. Update invoice_items table to include invoice_type field
-- 2. Update income_invoices table: remove customer_name and customer_id FK, rename customer_tax_id to customer_id

-- Add invoice_type column to invoice_items if it doesn't exist
ALTER TABLE invoice_items ADD COLUMN IF NOT EXISTS invoice_type VARCHAR(20) NOT NULL DEFAULT 'outcome';

-- Add check constraint for invoice_type
ALTER TABLE invoice_items DROP CONSTRAINT IF EXISTS invoice_items_invoice_type_check;
ALTER TABLE invoice_items ADD CONSTRAINT invoice_items_invoice_type_check CHECK (invoice_type IN ('income', 'outcome'));

-- Update existing records to have correct invoice_type
UPDATE invoice_items SET invoice_type = 'outcome' WHERE invoice_type NOT IN ('income', 'outcome');

-- Update income_invoices table structure
-- Remove customer_name column
ALTER TABLE income_invoices DROP COLUMN IF EXISTS customer_name;

-- Remove customer_id UUID FK column
ALTER TABLE income_invoices DROP COLUMN IF EXISTS customer_id;

-- Rename customer_tax_id to customer_id
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'income_invoices' AND column_name = 'customer_tax_id') THEN
        ALTER TABLE income_invoices RENAME COLUMN customer_tax_id TO customer_id;
    END IF;
END $$;

-- Add comment to clarify customer_id field
COMMENT ON COLUMN income_invoices.customer_id IS 'Customer tax ID (Cédula) - primary customer identifier';

-- Update indexes
DROP INDEX IF EXISTS idx_income_invoices_customer_id;
CREATE INDEX IF NOT EXISTS idx_income_invoices_customer_id ON income_invoices(customer_id);
