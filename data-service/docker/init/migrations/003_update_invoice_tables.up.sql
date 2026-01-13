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

-- Update the foreign key constraint for invoice_items to reference outcome_invoices instead of purchase_invoices
-- First, we need to drop and recreate the foreign key constraint
ALTER TABLE invoice_items DROP CONSTRAINT IF EXISTS invoice_details_invoice_id_fkey;
ALTER TABLE invoice_items ADD CONSTRAINT invoice_items_invoice_id_fkey
    FOREIGN KEY (invoice_id) REFERENCES outcome_invoices(id) ON DELETE CASCADE;

-- Update comments to reflect new naming
COMMENT ON TABLE outcome_invoices IS 'Expense invoices from suppliers (formerly purchase_invoices)';
COMMENT ON TABLE income_invoices IS 'Income invoices for customers (formerly customer_invoices)';
COMMENT ON TABLE invoice_items IS 'Line items for outcome invoices (formerly invoice_details)';

-- Add any missing indexes if needed
CREATE INDEX IF NOT EXISTS idx_outcome_invoices_supplier_id ON outcome_invoices(supplier_id);
CREATE INDEX IF NOT EXISTS idx_income_invoices_customer_id ON income_invoices(customer_id);
CREATE INDEX IF NOT EXISTS idx_income_invoices_order_id ON income_invoices(order_id);
CREATE INDEX IF NOT EXISTS idx_invoice_items_invoice_id ON invoice_items(invoice_id);