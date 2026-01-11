-- Bar-Restaurant Database Schema
-- Database: barrest_db
-- Version: 1.0
-- Based on entities.md specifications

-- Enable required extensions
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Create sequences
CREATE SEQUENCE IF NOT EXISTS existence_reference_seq START 1;
CREATE SEQUENCE IF NOT EXISTS order_number_seq START 1;
CREATE SEQUENCE IF NOT EXISTS invoice_number_seq START 1;

-- =============================================================================
-- CORE ENTITIES
-- =============================================================================

-- 1. Tables
CREATE TABLE tables (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    table_number VARCHAR(10) UNIQUE NOT NULL,
    capacity INTEGER NOT NULL DEFAULT 4,
    status VARCHAR(20) NOT NULL DEFAULT 'available' CHECK (status IN ('available', 'occupied', 'reserved')),
    floor_plan_position JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 2. Customers
CREATE TABLE customers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(20) NOT NULL,
    email VARCHAR(255),
    date_of_birth DATE,
    loyalty_points INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 6. Menu Categories (top level: Drinks, Desserts, etc.)
CREATE TABLE menu_categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) UNIQUE NOT NULL,
    display_order INTEGER NOT NULL DEFAULT 0,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 7. Sub Menus (second level: Smoothies, Sodas, etc. - grouped by category)
CREATE TABLE sub_menus (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    category_id UUID NOT NULL REFERENCES menu_categories(id) ON DELETE RESTRICT,
    image_url VARCHAR(500),
    item_type VARCHAR(20) NOT NULL CHECK (item_type IN ('kitchen', 'bar')),
    display_order INTEGER NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 8. Menu Items (third level: Banana Smoothie, Pineapple Smoothie, etc. - with pricing)
CREATE TABLE menu_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    sub_menu_id UUID NOT NULL REFERENCES sub_menus(id) ON DELETE RESTRICT,
    price DECIMAL(10,2) NOT NULL,
    item_cost DECIMAL(10,2),
    happy_hour_price DECIMAL(10,2),
    image_url VARCHAR(500),
    is_available BOOLEAN NOT NULL DEFAULT true,
    preparation_time INTEGER,
    display_order INTEGER NOT NULL DEFAULT 0,
    menu_types JSONB NOT NULL DEFAULT '["lunch", "dinner"]',
    dietary_tags JSONB,
    allergens JSONB,
    is_alcoholic BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 9. Stock Item Categories
CREATE TABLE stock_item_categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 10. Stock Items
CREATE TABLE stock_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) UNIQUE NOT NULL,
    unit VARCHAR(50) NOT NULL,
    description TEXT,
    stock_item_category_id UUID REFERENCES stock_item_categories(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 11. Menu Item Stock Items (Ingredients - junction table)
CREATE TABLE menu_item_stock_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    menu_item_id UUID NOT NULL REFERENCES menu_items(id) ON DELETE CASCADE,
    stock_item_id UUID NOT NULL REFERENCES stock_items(id) ON DELETE CASCADE,
    quantity DECIMAL(10,2) NOT NULL CHECK (quantity > 0),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(menu_item_id, stock_item_id)
);

-- 10. Suppliers
CREATE TABLE suppliers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) UNIQUE NOT NULL,
    contact_name VARCHAR(255),
    phone VARCHAR(20),
    email VARCHAR(255),
    address TEXT,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 11. Invoices (Purchase Invoices)
CREATE TABLE purchase_invoices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_number VARCHAR(100) UNIQUE NOT NULL,
    supplier_id UUID REFERENCES suppliers(id) ON DELETE SET NULL,
    transaction_date DATE NOT NULL,
    total_amount DECIMAL(12,2),
    image_url VARCHAR(500),
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 12. Invoice Details
CREATE TABLE invoice_details (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_id UUID NOT NULL REFERENCES purchase_invoices(id) ON DELETE CASCADE,
    stock_item_id UUID REFERENCES stock_items(id) ON DELETE SET NULL,
    detail TEXT NOT NULL,
    count DECIMAL(10,2) NOT NULL CHECK (count > 0),
    unit_type VARCHAR(50) NOT NULL,
    price DECIMAL(10,2) NOT NULL CHECK (price > 0),
    items_per_unit INTEGER NOT NULL DEFAULT 1 CHECK (items_per_unit > 0),
    total DECIMAL(12,2) GENERATED ALWAYS AS (count * price) STORED,
    expiration_date DATE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 13. Existences
CREATE TABLE existences (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    existence_reference_code INTEGER UNIQUE NOT NULL DEFAULT nextval('existence_reference_seq'),
    stock_item_id UUID NOT NULL REFERENCES stock_items(id) ON DELETE CASCADE,
    invoice_detail_id UUID NOT NULL REFERENCES invoice_details(id) ON DELETE CASCADE,
    units_purchased DECIMAL(10,2) NOT NULL,
    units_available DECIMAL(10,2) NOT NULL,
    unit_type VARCHAR(50) NOT NULL,
    items_per_unit INTEGER NOT NULL CHECK (items_per_unit > 0),
    cost_per_item DECIMAL(10,2) GENERATED ALWAYS AS (cost_per_unit / items_per_unit) STORED,
    cost_per_unit DECIMAL(10,2) NOT NULL,
    total_purchase_cost DECIMAL(12,2) GENERATED ALWAYS AS (units_purchased * cost_per_unit) STORED,
    remaining_value DECIMAL(12,2) GENERATED ALWAYS AS (units_available * cost_per_unit) STORED,
    expiration_date DATE,
    income_margin_percentage DECIMAL(5,2) DEFAULT 30.00,
    income_margin_amount DECIMAL(10,2) DEFAULT 0.00,
    minimum_price DECIMAL(10,2) DEFAULT 0.00,
    maximum_price DECIMAL(10,2),
    final_price DECIMAL(10,2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT check_final_price_range CHECK (final_price IS NULL OR (final_price >= minimum_price AND (maximum_price IS NULL OR final_price <= maximum_price)))
);

-- 14. Customer Favorites
CREATE TABLE customer_favorites (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id UUID NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
    menu_item_id UUID NOT NULL REFERENCES menu_items(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(customer_id, menu_item_id)
);

-- =============================================================================
-- BUSINESS LOGIC ENTITIES
-- =============================================================================

-- 17. Staff
CREATE TABLE staff (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(100) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    role VARCHAR(30) NOT NULL CHECK (role IN ('waiter', 'bartender', 'chef', 'manager', 'admin', 'dj_karaoke_operator')),
    is_active BOOLEAN NOT NULL DEFAULT true,
    last_login_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 18. Sessions (User authentication sessions)
-- Simplified: only session_id and token stored (other info in JWT token)
CREATE TABLE sessions (
    session_id VARCHAR(255) PRIMARY KEY,
    token TEXT NOT NULL
);

-- 24. Promotions
CREATE TABLE promotions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    promotion_type VARCHAR(30) NOT NULL CHECK (promotion_type IN ('happy_hour', 'daily_deal', 'seasonal', 'birthday', 'loyalty_reward')),
    discount_type VARCHAR(20) NOT NULL CHECK (discount_type IN ('percentage', 'fixed_amount')),
    discount_value DECIMAL(10,2) NOT NULL,
    points_required INTEGER,
    recurrence_type VARCHAR(20) NOT NULL DEFAULT 'none' CHECK (recurrence_type IN ('none', 'daily', 'weekly', 'monthly')),
    start_date DATE NOT NULL,
    end_date DATE,
    from_time TIME NOT NULL,
    to_time TIME NOT NULL,
    days_of_week JSONB,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 3. Orders
CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_number VARCHAR(50) UNIQUE NOT NULL DEFAULT 'ORD-' || nextval('order_number_seq'),
    table_id UUID NOT NULL REFERENCES tables(id) ON DELETE RESTRICT,
    customer_id UUID REFERENCES customers(id) ON DELETE SET NULL,
    customer_name VARCHAR(255),
    status VARCHAR(20) NOT NULL DEFAULT 'created' CHECK (status IN ('created', 'preparing', 'paid', 'cancelled')),
    payment_status VARCHAR(20) NOT NULL DEFAULT 'unpaid' CHECK (payment_status IN ('unpaid', 'requested', 'paid', 'refunded')),
    payment_requested_at TIMESTAMP,
    paid_at TIMESTAMP,
    promotion_id UUID REFERENCES promotions(id) ON DELETE SET NULL,
    total_amount DECIMAL(10,2) DEFAULT 0.00,
    promotion_discount_amount DECIMAL(10,2) DEFAULT 0.00,
    tax_amount DECIMAL(10,2) DEFAULT 0.00,
    service_charge DECIMAL(10,2) DEFAULT 0.00,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    confirmed_at TIMESTAMP
);

-- 4. Order Items
CREATE TABLE order_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    menu_item_id UUID NOT NULL REFERENCES menu_items(id) ON DELETE RESTRICT,
    quantity INTEGER NOT NULL DEFAULT 1,
    unit_price DECIMAL(10,2) NOT NULL,
    subtotal DECIMAL(10,2) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'requested' CHECK (status IN ('requested', 'preparing', 'ready', 'delivered', 'lost')),
    order_type VARCHAR(20) CHECK (order_type IN ('for_here', 'to_go')),
    special_instructions TEXT,
    lost_reason VARCHAR(255),
    fault_type VARCHAR(20) CHECK (fault_type IN ('restaurant_fault', 'customer_fault')),
    lost_cost DECIMAL(10,2),
    requested_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    preparing_at TIMESTAMP,
    ready_at TIMESTAMP,
    delivered_at TIMESTAMP,
    lost_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 15. Payments
CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE RESTRICT,
    amount DECIMAL(10,2) NOT NULL,
    payment_method VARCHAR(20) NOT NULL CHECK (payment_method IN ('cash', 'credit_card', 'debit_card', 'apple_pay', 'google_pay')),
    payment_status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (payment_status IN ('pending', 'processing', 'completed', 'failed', 'refunded')),
    exact_change_amount DECIMAL(10,2),
    cash_amount_provided DECIMAL(10,2),
    tip_amount DECIMAL(10,2) DEFAULT 0.00,
    processed_by_staff_id UUID REFERENCES staff(id) ON DELETE SET NULL,
    transaction_id VARCHAR(255),
    processed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 16. Customer Invoices (Receipts)
CREATE TABLE customer_invoices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE RESTRICT,
    payment_id UUID REFERENCES payments(id) ON DELETE SET NULL,
    invoice_number VARCHAR(50) UNIQUE NOT NULL DEFAULT 'INV-' || nextval('invoice_number_seq'),
    invoice_type VARCHAR(20) NOT NULL CHECK (invoice_type IN ('sales', 'credit_note')),
    customer_name VARCHAR(255) NOT NULL,
    customer_tax_id VARCHAR(50),
    subtotal DECIMAL(10,2) NOT NULL,
    tax_amount DECIMAL(10,2) NOT NULL,
    service_charge DECIMAL(10,2) DEFAULT 0.00,
    total_amount DECIMAL(10,2) NOT NULL,
    payment_method VARCHAR(50) NOT NULL,
    xml_data TEXT,
    digital_signature TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'generated', 'sent', 'cancelled')),
    generated_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 18. Request Notifications
CREATE TABLE request_notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    request_type VARCHAR(30) NOT NULL CHECK (request_type IN ('payment', 'assistance', 'refill', 'issue_report', 'special_request')),
    table_id UUID NOT NULL REFERENCES tables(id) ON DELETE CASCADE,
    order_id UUID REFERENCES orders(id) ON DELETE SET NULL,
    customer_name VARCHAR(255),
    message TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'in_progress', 'completed', 'cancelled')),
    handled_by_staff_id UUID REFERENCES staff(id) ON DELETE SET NULL,
    handled_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 19. Karaoke Song Requests
CREATE TABLE karaoke_song_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    table_id UUID NOT NULL REFERENCES tables(id) ON DELETE CASCADE,
    customer_name VARCHAR(255) NOT NULL,
    song_title VARCHAR(255) NOT NULL,
    artist_name VARCHAR(255) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'queued' CHECK (status IN ('queued', 'playing', 'completed', 'skipped', 'removed')),
    position_in_queue INTEGER NOT NULL,
    skip_reason TEXT,
    skipped_by_staff_id UUID REFERENCES staff(id) ON DELETE SET NULL,
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 20. Karaoke Song Library
CREATE TABLE karaoke_song_library (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    song_title VARCHAR(255) NOT NULL,
    artist_name VARCHAR(255) NOT NULL,
    duration_estimate INTEGER,
    is_available BOOLEAN NOT NULL DEFAULT true,
    popularity_count INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 21. Loyalty Points Transactions
CREATE TABLE loyalty_points_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id UUID NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
    order_id UUID REFERENCES orders(id) ON DELETE SET NULL,
    transaction_type VARCHAR(20) NOT NULL CHECK (transaction_type IN ('earned', 'spent', 'expired', 'bonus', 'referral')),
    points INTEGER NOT NULL,
    description VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- =============================================================================
-- ADVANCED FEATURES ENTITIES (FUTURE)
-- =============================================================================

-- 23. Reservations
CREATE TABLE reservations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    table_id UUID NOT NULL REFERENCES tables(id) ON DELETE CASCADE,
    customer_id UUID REFERENCES customers(id) ON DELETE SET NULL,
    customer_name VARCHAR(255) NOT NULL,
    customer_phone VARCHAR(20) NOT NULL,
    reservation_date DATE NOT NULL,
    reservation_time TIME NOT NULL,
    party_size INTEGER NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'confirmed' CHECK (status IN ('confirmed', 'cancelled', 'completed', 'no_show')),
    special_requests TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 25. Reviews
CREATE TABLE reviews (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id UUID REFERENCES customers(id) ON DELETE SET NULL,
    order_id UUID REFERENCES orders(id) ON DELETE SET NULL,
    menu_item_id UUID REFERENCES menu_items(id) ON DELETE SET NULL,
    rating INTEGER NOT NULL CHECK (rating >= 1 AND rating <= 5),
    comment TEXT,
    is_moderated BOOLEAN NOT NULL DEFAULT false,
    is_approved BOOLEAN,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- =============================================================================
-- SETTINGS TABLE FOR CENTRALIZED CONFIGURATION
-- =============================================================================

CREATE TABLE settings (
    setting_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    service VARCHAR(50) NOT NULL,
    key VARCHAR(100) NOT NULL,
    value TEXT NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(service, key)
);

-- =============================================================================
-- INDEXES FOR PERFORMANCE
-- =============================================================================

-- Tables indexes
CREATE INDEX idx_tables_status ON tables(status);

-- Customers indexes
CREATE INDEX idx_customers_customer_id ON customers(customer_id);
CREATE INDEX idx_customers_phone ON customers(phone);

-- Sub Menus indexes
CREATE INDEX idx_sub_menus_category ON sub_menus(category_id);
CREATE INDEX idx_sub_menus_active ON sub_menus(is_active);
CREATE INDEX idx_sub_menus_item_type ON sub_menus(item_type);

-- Menu Items indexes
CREATE INDEX idx_menu_items_sub_menu ON menu_items(sub_menu_id);
CREATE INDEX idx_menu_items_available ON menu_items(is_available);

-- Stock Items indexes
CREATE INDEX idx_stock_items_category ON stock_items(stock_item_category_id);

-- Menu Item Stock Items indexes
CREATE INDEX idx_menu_item_stock_items_menu ON menu_item_stock_items(menu_item_id);
CREATE INDEX idx_menu_item_stock_items_stock ON menu_item_stock_items(stock_item_id);

-- Existences indexes
CREATE INDEX idx_existences_stock_item ON existences(stock_item_id);
CREATE INDEX idx_existences_reference_code ON existences(existence_reference_code);
CREATE INDEX idx_existences_available ON existences(units_available);
CREATE INDEX idx_existences_expiration ON existences(expiration_date);

-- Orders indexes
CREATE INDEX idx_orders_table ON orders(table_id);
CREATE INDEX idx_orders_customer ON orders(customer_id);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_payment_status ON orders(payment_status);
CREATE INDEX idx_orders_created_at ON orders(created_at);

-- Order Items indexes
CREATE INDEX idx_order_items_order ON order_items(order_id);
CREATE INDEX idx_order_items_menu_item ON order_items(menu_item_id);
CREATE INDEX idx_order_items_status ON order_items(status);

-- Payments indexes
CREATE INDEX idx_payments_order ON payments(order_id);
CREATE INDEX idx_payments_status ON payments(payment_status);

-- Request Notifications indexes
CREATE INDEX idx_request_notifications_table ON request_notifications(table_id);
CREATE INDEX idx_request_notifications_status ON request_notifications(status);
CREATE INDEX idx_request_notifications_type_status ON request_notifications(request_type, status);

-- Karaoke Song Requests indexes
CREATE INDEX idx_karaoke_requests_table ON karaoke_song_requests(table_id);
CREATE INDEX idx_karaoke_requests_status ON karaoke_song_requests(status);
CREATE INDEX idx_karaoke_requests_queue ON karaoke_song_requests(status, position_in_queue);

-- Sessions (auth) indexes
CREATE INDEX idx_sessions_token ON sessions(token);

-- Promotions indexes
CREATE INDEX idx_promotions_type ON promotions(promotion_type);
CREATE INDEX idx_promotions_active ON promotions(is_active);
CREATE INDEX idx_promotions_dates ON promotions(start_date, end_date);

-- Settings indexes
CREATE INDEX idx_settings_service ON settings(service);
CREATE INDEX idx_settings_key ON settings(key);

-- =============================================================================
-- TRIGGERS FOR AUTOMATIC UPDATES
-- =============================================================================

-- Update timestamps trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Apply triggers to all tables with updated_at
CREATE TRIGGER update_tables_updated_at BEFORE UPDATE ON tables 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_customers_updated_at BEFORE UPDATE ON customers 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_menu_categories_updated_at BEFORE UPDATE ON menu_categories 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_sub_menus_updated_at BEFORE UPDATE ON sub_menus 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_menu_items_updated_at BEFORE UPDATE ON menu_items 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_stock_item_categories_updated_at BEFORE UPDATE ON stock_item_categories 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_stock_items_updated_at BEFORE UPDATE ON stock_items 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_menu_item_stock_items_updated_at BEFORE UPDATE ON menu_item_stock_items 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_suppliers_updated_at BEFORE UPDATE ON suppliers 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_purchase_invoices_updated_at BEFORE UPDATE ON purchase_invoices 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_invoice_details_updated_at BEFORE UPDATE ON invoice_details 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_existences_updated_at BEFORE UPDATE ON existences 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_staff_updated_at BEFORE UPDATE ON staff 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_promotions_updated_at BEFORE UPDATE ON promotions 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_orders_updated_at BEFORE UPDATE ON orders 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_order_items_updated_at BEFORE UPDATE ON order_items 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_payments_updated_at BEFORE UPDATE ON payments 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_customer_invoices_updated_at BEFORE UPDATE ON customer_invoices 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_request_notifications_updated_at BEFORE UPDATE ON request_notifications 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_karaoke_song_requests_updated_at BEFORE UPDATE ON karaoke_song_requests 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_karaoke_song_library_updated_at BEFORE UPDATE ON karaoke_song_library 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_reservations_updated_at BEFORE UPDATE ON reservations 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_reviews_updated_at BEFORE UPDATE ON reviews 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_settings_updated_at BEFORE UPDATE ON settings 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- =============================================================================
-- DEFAULT DATA
-- =============================================================================

-- Insert default settings for all services
INSERT INTO settings (service, key, value, description) VALUES
-- General Database Settings
('General', 'DB_HOST', 'barrest_postgres', 'Database host address'),
('General', 'DB_PORT', '5432', 'Database port number'),
('General', 'DB_USER', 'postgres', 'Database username'),
('General', 'DB_PASSWORD', 'postgres123', 'Database password'),
('General', 'DB_NAME', 'barrest_db', 'Database name'),
('General', 'DB_SSL_MODE', 'disable', 'Database SSL mode'),
('General', 'LOG_LEVEL', 'info', 'Logging level'),
('General', 'JWT_SECRET', 'barrest-super-secret-jwt-key-change-in-production', 'JWT signing secret'),
('General', 'JWT_EXPIRATION_TIME', '30m', 'JWT token expiration time'),
('General', 'DEFAULT_TAX_RATE', '13.0', 'Default tax rate (Costa Rica IVA)'),
('General', 'DEFAULT_SERVICE_RATE', '10.0', 'Default service charge rate'),

-- Gateway Service Settings
('Gateway', 'SERVER_PORT', '8082', 'Port for the gateway service'),
('Gateway', 'SESSION_SERVICE_URL', 'http://barrest_session:8081', 'URL for session service'),
('Gateway', 'ORDERS_SERVICE_URL', 'http://barrest_orders:8083', 'URL for orders service'),
('Gateway', 'MENU_SERVICE_URL', 'http://barrest_menu:8087', 'URL for menu service'),
('Gateway', 'INVENTORY_SERVICE_URL', 'http://barrest_inventory:8084', 'URL for inventory service'),
('Gateway', 'PAYMENT_SERVICE_URL', 'http://barrest_payment:8088', 'URL for payment service'),
('Gateway', 'CUSTOMER_SERVICE_URL', 'http://barrest_customer:8089', 'URL for customer service'),

-- Session Service Settings
('Session', 'SERVER_PORT', '8081', 'Port for the session service'),

-- Orders Service Settings
('Orders', 'SERVER_PORT', '8083', 'Port for the orders service'),

-- Menu Service Settings
('Menu', 'SERVER_PORT', '8087', 'Port for the menu service'),

-- Inventory Service Settings
('Inventory', 'SERVER_PORT', '8084', 'Port for the inventory service'),

-- Payment Service Settings
('Payment', 'SERVER_PORT', '8088', 'Port for the payment service'),

-- Customer Service Settings
('Customer', 'SERVER_PORT', '8089', 'Port for the customer service'),

-- Data Service Settings
('Data', 'DATA_SERVER_PORT', '8086', 'Port for the data service'),
('Data', 'DB_MAX_OPEN_CONNS', '25', 'Maximum number of open connections'),
('Data', 'DB_MAX_IDLE_CONNS', '5', 'Maximum number of idle connections'),
('Data', 'DB_CONN_MAX_LIFETIME', '5m', 'Maximum lifetime of connections'),
('Data', 'DB_CONN_MAX_IDLE_TIME', '5m', 'Maximum idle time of connections'),
('Data', 'DB_CONNECT_TIMEOUT', '10s', 'Database connection timeout'),
('Data', 'DB_QUERY_TIMEOUT', '30s', 'Database query timeout'),

-- UI Service Settings
('UI', 'UI_PORT', '3000', 'Port for the UI service'),
('UI', 'GATEWAY_URL', 'http://barrest_gateway:8082', 'Gateway service URL');

-- Insert default admin user (password: admin)
INSERT INTO staff (username, email, password_hash, first_name, last_name, role) VALUES
('admin', 'admin@barrest.com', '$2a$10$o4Pv9FXpT5HNIaPRS7U.xuWj2b8EyfuGp6EhGKByB8d3vdGNkgYYq', 'System', 'Administrator', 'admin');

-- Insert default menu categories
INSERT INTO menu_categories (name, display_order, description) VALUES
('Appetizers', 1, 'Starters and appetizers'),
('Main Course', 2, 'Main dishes'),
('Desserts', 3, 'Sweet treats and desserts'),
('Beer', 4, 'Beer selection'),
('Cocktails', 5, 'Mixed drinks and cocktails'),
('Wine', 6, 'Wine selection'),
('Soft Drinks', 7, 'Non-alcoholic beverages'),
('Snacks', 8, 'Bar snacks');

-- Insert default stock item categories
INSERT INTO stock_item_categories (name, description) VALUES
('Meat', 'Meat products'),
('Vegetables', 'Fresh vegetables'),
('Dairy', 'Milk, cheese, cream, eggs'),
('Beverages', 'Drinks and beverages'),
('Spices', 'Spices and seasonings'),
('Grains', 'Rice, pasta, bread'),
('Alcohol', 'Alcoholic beverages'),
('Condiments', 'Sauces and condiments'),
('Frozen', 'Frozen products'),
('Dry Goods', 'Dry ingredients');

-- Insert sample tables
INSERT INTO tables (table_number, capacity, status) VALUES
('T1', 2, 'available'),
('T2', 2, 'available'),
('T3', 4, 'available'),
('T4', 4, 'available'),
('T5', 4, 'available'),
('T6', 6, 'available'),
('T7', 6, 'available'),
('T8', 8, 'available'),
('BAR1', 1, 'available'),
('BAR2', 1, 'available'),
('BAR3', 1, 'available'),
('BAR4', 1, 'available');

-- =============================================================================
-- END OF SCHEMA
-- =============================================================================

