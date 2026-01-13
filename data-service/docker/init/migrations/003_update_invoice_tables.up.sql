-- Migration: Update invoice-related tables to match new naming and structure
-- Version: 003
-- Date: 2026-01-12

-- Rename purchase_invoices to outcome_invoices
ALTER TABLE purchase_invoices RENAME TO outcome_invoices;

-- Rename invoice_details to invoice_items
ALTER TABLE invoice_details RENAME TO invoice_items;

-- Rename customer_invoices to income_invoices
ALTER TABLE customer_invoices RENAME TO income_invoices;

-- Add customer_id column to income_invoices
ALTER TABLE income_invoices ADD COLUMN customer_id UUID REFERENCES customers(id) ON DELETE SET NULL;

-- Add invoice_type to invoice_items to distinguish between outcome and income invoices
ALTER TABLE invoice_items ADD COLUMN IF NOT EXISTS invoice_type VARCHAR(20) NOT NULL DEFAULT 'outcome' CHECK (invoice_type IN ('outcome', 'income'));

-- Update existing records to have the correct invoice_type (assuming they are all outcome for now)
UPDATE invoice_items SET invoice_type = 'outcome' WHERE invoice_type IS NULL OR invoice_type = '';

-- Update comments to reflect new naming
COMMENT ON TABLE outcome_invoices IS 'Expense invoices from suppliers (formerly purchase_invoices)';
COMMENT ON TABLE income_invoices IS 'Income invoices for customers (formerly customer_invoices)';
COMMENT ON TABLE invoice_items IS 'Line items for outcome invoices (formerly invoice_details)';

-- Add any missing indexes if needed
CREATE INDEX IF NOT EXISTS idx_outcome_invoices_supplier_id ON outcome_invoices(supplier_id);
CREATE INDEX IF NOT EXISTS idx_income_invoices_customer_id ON income_invoices(customer_id);
CREATE INDEX IF NOT EXISTS idx_income_invoices_order_id ON income_invoices(order_id);
CREATE INDEX IF NOT EXISTS idx_invoice_items_invoice_id ON invoice_items(invoice_id);