-- Rollback: Revert invoice table updates to original naming
-- Version: 003
-- Date: 2026-01-12

-- Only rollback if the new table names exist
DO $$
BEGIN
    -- Check if new table names exist before attempting rollback
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'outcome_invoices') AND
       EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'income_invoices') AND
       EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'invoice_items') THEN

        -- Remove added indexes
        DROP INDEX IF EXISTS idx_outcome_invoices_supplier_id;
        DROP INDEX IF EXISTS idx_outcome_invoices_transaction_date;
        DROP INDEX IF EXISTS idx_income_invoices_order_id;
        DROP INDEX IF EXISTS idx_income_invoices_payment_id;
        DROP INDEX IF EXISTS idx_invoice_items_invoice_id;

        -- Remove comments
        COMMENT ON TABLE outcome_invoices IS NULL;
        COMMENT ON TABLE income_invoices IS NULL;
        COMMENT ON TABLE invoice_items IS NULL;

        -- Remove the added customer_id column
        IF EXISTS (SELECT 1 FROM information_schema.columns
                   WHERE table_name = 'income_invoices' AND column_name = 'customer_id') THEN
            ALTER TABLE income_invoices DROP COLUMN customer_id;
        END IF;

        -- Remove invoice_type column if it exists
        IF EXISTS (SELECT 1 FROM information_schema.columns
                   WHERE table_name = 'invoice_items' AND column_name = 'invoice_type') THEN
            ALTER TABLE invoice_items DROP COLUMN invoice_type;
        END IF;

        -- Rename tables back to original names
        ALTER TABLE outcome_invoices RENAME TO purchase_invoices;
        ALTER TABLE invoice_items RENAME TO invoice_details;
        ALTER TABLE income_invoices RENAME TO customer_invoices;
    END IF;
END $$;