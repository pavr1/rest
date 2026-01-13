-- Rollback: Revert invoice table updates to original naming
-- Version: 003
-- Date: 2026-01-12

-- Remove added indexes
DROP INDEX IF EXISTS idx_outcome_invoices_supplier_id;
DROP INDEX IF EXISTS idx_income_invoices_customer_id;
DROP INDEX IF EXISTS idx_income_invoices_order_id;
DROP INDEX IF EXISTS idx_invoice_items_invoice_id;

-- Remove comments
COMMENT ON TABLE outcome_invoices IS NULL;
COMMENT ON TABLE income_invoices IS NULL;
COMMENT ON TABLE invoice_items IS NULL;

-- Remove the updated foreign key constraint and restore the original
ALTER TABLE invoice_items DROP CONSTRAINT IF EXISTS invoice_items_invoice_id_fkey;
ALTER TABLE invoice_items ADD CONSTRAINT invoice_details_invoice_id_fkey
    FOREIGN KEY (invoice_id) REFERENCES purchase_invoices(id) ON DELETE CASCADE;

-- Remove the added customer_id column
ALTER TABLE income_invoices DROP COLUMN IF EXISTS customer_id;

-- Rename tables back to original names
ALTER TABLE outcome_invoices RENAME TO purchase_invoices;
ALTER TABLE invoice_items RENAME TO invoice_details;
ALTER TABLE income_invoices RENAME TO customer_invoices;