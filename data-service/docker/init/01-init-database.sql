-- Bar-Restaurant Database Schema
-- Database: barrest_db
-- Version: 1.0
-- Based on entities.md specifications

-- Enable required extensions
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Create sequences
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

-- 7. Menu Sub-Categories (second level: Smoothies, Sodas, etc. - grouped by category)
CREATE TABLE menu_sub_categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    category_id UUID NOT NULL REFERENCES menu_categories(id) ON DELETE RESTRICT,
    item_type VARCHAR(20) NOT NULL CHECK (item_type IN ('kitchen', 'bar')),
    display_order INTEGER NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 8. Menu Variants (third level: Banana Smoothie, Pineapple Smoothie, etc. - with pricing)
CREATE TABLE menu_variants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    sub_category_id UUID NOT NULL REFERENCES menu_sub_categories(id) ON DELETE RESTRICT,
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

-- 9. Stock Categories
CREATE TABLE stock_categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) UNIQUE NOT NULL,
    description TEXT,
    display_order INTEGER NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 10. Stock Sub-Categories
CREATE TABLE stock_sub_categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    stock_category_id UUID NOT NULL REFERENCES stock_categories(id) ON DELETE CASCADE,
    display_order INTEGER NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 11. Stock Variants (simplified - actual counts are tracked in stock_count table)
CREATE TABLE stock_variants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    stock_sub_category_id UUID NOT NULL REFERENCES stock_sub_categories(id) ON DELETE CASCADE,
    avg_cost DECIMAL(10,2) DEFAULT 0,  -- Average cost per portion across active stock counts
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 12. Menu Ingredients (links menu items to stock variants they require)
CREATE TABLE menu_ingredients (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    menu_variant_id UUID NOT NULL REFERENCES menu_variants(id) ON DELETE CASCADE,
    stock_variant_id UUID REFERENCES stock_variants(id) ON DELETE RESTRICT,
    menu_sub_category_id UUID REFERENCES menu_sub_categories(id) ON DELETE RESTRICT,
    quantity DECIMAL(10,2) NOT NULL CHECK (quantity > 0),
    is_optional BOOLEAN NOT NULL DEFAULT false,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    -- Ensure exactly one of stock_variant_id or menu_sub_category_id is provided
    CONSTRAINT chk_ingredient_type CHECK (
        (stock_variant_id IS NOT NULL AND menu_sub_category_id IS NULL) OR
        (stock_variant_id IS NULL AND menu_sub_category_id IS NOT NULL)
    )
);

-- Menu Ingredients indexes
CREATE INDEX idx_menu_ingredients_stock_variant ON menu_ingredients(stock_variant_id);
CREATE INDEX idx_menu_ingredients_menu_sub_category ON menu_ingredients(menu_sub_category_id);

-- 13. Suppliers
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

-- 14. Outcome Invoices (Supplier Purchase Invoices)
CREATE TABLE outcome_invoices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_number VARCHAR(100) UNIQUE NOT NULL,
    supplier_id UUID REFERENCES suppliers(id) ON DELETE SET NULL,
    transaction_date DATE NOT NULL,
    due_date DATE,
    subtotal DECIMAL(12,2),
    tax_amount DECIMAL(12,2),
    discount_amount DECIMAL(12,2) DEFAULT 0,
    total_amount DECIMAL(12,2),
    image_url VARCHAR(500),
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 15. Stock Count (inventory tracking - links stock variants to purchases)
CREATE TABLE stock_count (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    stock_variant_id UUID NOT NULL REFERENCES stock_variants(id) ON DELETE CASCADE,
    invoice_id UUID REFERENCES outcome_invoices(id) ON DELETE CASCADE,  -- Nullable for manual stock counts
    count DECIMAL(10,2) NOT NULL CHECK (count > 0),
    unit VARCHAR(50) NOT NULL,  -- Supported: kg, g, l, ml (all convert to kg for cost calculation)
    unit_price DECIMAL(10,2),  -- Price per unit (passed from invoice or manual input)
    cost_per_portion DECIMAL(10,2),  -- Calculated: unit_price / num_portions (based on default portion grams)
    purchased_at TIMESTAMP NOT NULL,
    is_out BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Stock Count indexes
CREATE INDEX idx_stock_count_variant ON stock_count(stock_variant_id);
CREATE INDEX idx_stock_count_invoice ON stock_count(invoice_id);
CREATE INDEX idx_stock_count_is_out ON stock_count(stock_variant_id, is_out);

-- 16. Invoice Items
CREATE TABLE invoice_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_id UUID NOT NULL,
    stock_variant_id UUID REFERENCES stock_variants(id) ON DELETE SET NULL,
    detail TEXT,
    count DECIMAL(10,2) NOT NULL CHECK (count > 0),
    unit_type VARCHAR(50) NOT NULL,
    price DECIMAL(10,2) NOT NULL CHECK (price > 0),
    items_per_unit INTEGER NOT NULL DEFAULT 1 CHECK (items_per_unit > 0),
    total DECIMAL(12,2) GENERATED ALWAYS AS (count * price) STORED,
    expiration_date DATE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Invoice Items index
CREATE INDEX idx_invoice_items_stock_variant ON invoice_items(stock_variant_id);


-- 17. Customer Favorites
CREATE TABLE customer_favorites (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id UUID NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
    menu_variant_id UUID NOT NULL REFERENCES menu_variants(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(customer_id, menu_variant_id)
);

-- =============================================================================
-- BUSINESS LOGIC ENTITIES
-- =============================================================================

-- 18. Staff
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

-- 18. Promotions
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

-- 19. Orders
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

-- 20. Order Items
CREATE TABLE order_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    menu_variant_id UUID NOT NULL REFERENCES menu_variants(id) ON DELETE RESTRICT,
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

-- 21. Payments
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

-- 22. Income Invoices (Customer Sales Invoices)
CREATE TABLE income_invoices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID REFERENCES orders(id) ON DELETE SET NULL,
    payment_id UUID REFERENCES payments(id) ON DELETE SET NULL,
    customer_id VARCHAR(50),
    invoice_number VARCHAR(50) UNIQUE NOT NULL,
    invoice_type VARCHAR(20) NOT NULL DEFAULT 'sales' CHECK (invoice_type IN ('sales', 'credit_note')),
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

-- 23. Table Sessions
CREATE TABLE table_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    table_id UUID NOT NULL REFERENCES tables(id) ON DELETE CASCADE,
    started_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    ended_at TIMESTAMP,
    cleared_by_staff_id UUID REFERENCES staff(id) ON DELETE SET NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'closed'))
);

-- 24. Request Notifications
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

-- 25. Karaoke Song Requests
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

-- 26. Karaoke Song Library
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

-- 27. Loyalty Points Transactions
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

-- 28. Reservations
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

-- 29. Reviews
CREATE TABLE reviews (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id UUID REFERENCES customers(id) ON DELETE SET NULL,
    order_id UUID REFERENCES orders(id) ON DELETE SET NULL,
    menu_variant_id UUID REFERENCES menu_variants(id) ON DELETE SET NULL,
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
('Gateway', 'INVOICE_SERVICE_URL', 'http://barrest_invoice_service:8092', 'URL for invoice service'),
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

-- Insertar usuario administrador por defecto (contraseña: admin)
INSERT INTO staff (username, email, password_hash, first_name, last_name, role) VALUES
('admin', 'admin@barrest.com', '$2a$10$o4Pv9FXpT5HNIaPRS7U.xuWj2b8EyfuGp6EhGKByB8d3vdGNkgYYq', 'Sistema', 'Administrador', 'admin');

-- Insertar categorías del menú por defecto
INSERT INTO menu_categories (name, display_order, description) VALUES
('Entradas', 1, 'Aperitivos y entradas'),
('Platos Fuertes', 2, 'Platos principales'),
('Postres', 3, 'Dulces y postres'),
('Bebidas', 4, 'Cervezas, cocteles, vinos y refrescos'),
('Bocadillos', 5, 'Snacks de bar');
('Ensaladas', 6, 'Ensaladas y ensaladas'),

-- Insertar categorías de inventario por defecto
INSERT INTO stock_categories (name, description, display_order) VALUES
('Carnes', 'Res, pollo, cerdo, pescado', 1),
('Frutas y Verduras', 'Productos frescos', 2),
('Bebidas', 'Café, té, refrescos, lácteos y licores', 3),
('Salsas', 'Salsas y aderezos', 4),
('Congelados', 'Productos congelados', 5),
('Alacena', 'Productos de alacena', 6),
('Repostería', 'Pasteles, postres, galletas', 7),
('Limpieza', 'Productos de limpieza', 8),
('Envases', 'Platos, vasos, servilletas', 9);

-- Insertar sub-categorías de inventario por defecto
-- 1. Carnes
INSERT INTO stock_sub_categories (name, description, stock_category_id, display_order)
SELECT 'Res', 'Carne de res', id, 1 FROM stock_categories WHERE name = 'Carnes';
INSERT INTO stock_sub_categories (name, description, stock_category_id, display_order)
SELECT 'Pollo', 'Carne de pollo', id, 2 FROM stock_categories WHERE name = 'Carnes';
INSERT INTO stock_sub_categories (name, description, stock_category_id, display_order)
SELECT 'Cerdo', 'Carne de cerdo', id, 3 FROM stock_categories WHERE name = 'Carnes';
INSERT INTO stock_sub_categories (name, description, stock_category_id, display_order)
SELECT 'Pescado', 'Pescados frescos', id, 4 FROM stock_categories WHERE name = 'Carnes';
INSERT INTO stock_sub_categories (name, description, stock_category_id, display_order)
SELECT 'Mariscos', 'Camarones, pulpo, calamar', id, 5 FROM stock_categories WHERE name = 'Carnes';
INSERT INTO stock_sub_categories (name, description, stock_category_id, display_order)
SELECT 'Embutidos', 'Jamón, tocino, salchicha', id, 6 FROM stock_categories WHERE name = 'Carnes';

-- 2. Frutas y Verduras
INSERT INTO stock_sub_categories (name, description, stock_category_id, display_order)
SELECT 'Frutas', 'Frutas frescas', id, 1 FROM stock_categories WHERE name = 'Frutas y Verduras';
INSERT INTO stock_sub_categories (name, description, stock_category_id, display_order)
SELECT 'Verduras', 'Verduras frescas', id, 2 FROM stock_categories WHERE name = 'Frutas y Verduras';
INSERT INTO stock_sub_categories (name, description, stock_category_id, display_order)
SELECT 'Hierbas', 'Hierbas aromáticas', id, 3 FROM stock_categories WHERE name = 'Frutas y Verduras';
INSERT INTO stock_sub_categories (name, description, stock_category_id, display_order)
SELECT 'Hongos', 'Champiñones y setas', id, 4 FROM stock_categories WHERE name = 'Frutas y Verduras';

-- 3. Bebidas
INSERT INTO stock_sub_categories (name, description, stock_category_id, display_order)
SELECT 'Refrescos', 'Bebidas gaseosas', id, 1 FROM stock_categories WHERE name = 'Bebidas';
INSERT INTO stock_sub_categories (name, description, stock_category_id, display_order)
SELECT 'Jugos', 'Jugos naturales y envasados', id, 2 FROM stock_categories WHERE name = 'Bebidas';
INSERT INTO stock_sub_categories (name, description, stock_category_id, display_order)
SELECT 'Cervezas', 'Cervezas nacionales e importadas', id, 3 FROM stock_categories WHERE name = 'Bebidas';
INSERT INTO stock_sub_categories (name, description, stock_category_id, display_order)
SELECT 'Vinos', 'Vinos tintos, blancos y rosados', id, 4 FROM stock_categories WHERE name = 'Bebidas';
INSERT INTO stock_sub_categories (name, description, stock_category_id, display_order)
SELECT 'Licores', 'Destilados y licores', id, 5 FROM stock_categories WHERE name = 'Bebidas';
INSERT INTO stock_sub_categories (name, description, stock_category_id, display_order)
SELECT 'Café y Té', 'Café, té e infusiones', id, 6 FROM stock_categories WHERE name = 'Bebidas';
INSERT INTO stock_sub_categories (name, description, stock_category_id, display_order)
SELECT 'Lácteos', 'Leche y crema', id, 7 FROM stock_categories WHERE name = 'Bebidas';
INSERT INTO stock_sub_categories (name, description, stock_category_id, display_order)
SELECT 'Agua', 'Agua natural y mineral', id, 8 FROM stock_categories WHERE name = 'Bebidas';

-- 4. Salsas
INSERT INTO stock_sub_categories (name, description, stock_category_id, display_order)
SELECT 'Salsas Picantes', 'Salsas con chile', id, 1 FROM stock_categories WHERE name = 'Salsas';
INSERT INTO stock_sub_categories (name, description, stock_category_id, display_order)
SELECT 'Aderezos', 'Mayonesa, mostaza, etc.', id, 2 FROM stock_categories WHERE name = 'Salsas';
INSERT INTO stock_sub_categories (name, description, stock_category_id, display_order)
SELECT 'Salsas Base', 'Salsas para cocinar', id, 3 FROM stock_categories WHERE name = 'Salsas';
INSERT INTO stock_sub_categories (name, description, stock_category_id, display_order)
SELECT 'Aceites y Vinagres', 'Aceites y vinagres', id, 4 FROM stock_categories WHERE name = 'Salsas';

-- 5. Congelados
INSERT INTO stock_sub_categories (name, description, stock_category_id, display_order)
SELECT 'Postres Congelados', 'Helados y postres congelados', id, 1 FROM stock_categories WHERE name = 'Congelados';
INSERT INTO stock_sub_categories (name, description, stock_category_id, display_order)
SELECT 'Precosidos', 'Alimentos precocidos listos para freír', id, 2 FROM stock_categories WHERE name = 'Congelados';
INSERT INTO stock_sub_categories (name, description, stock_category_id, display_order)
SELECT 'Hielo', 'Hielo para bebidas', id, 3 FROM stock_categories WHERE name = 'Congelados';

-- 6. Alacena
INSERT INTO stock_sub_categories (name, description, stock_category_id, display_order)
SELECT 'Granos', 'Arroz, frijoles, lentejas', id, 1 FROM stock_categories WHERE name = 'Alacena';
INSERT INTO stock_sub_categories (name, description, stock_category_id, display_order)
SELECT 'Pastas', 'Pastas secas', id, 2 FROM stock_categories WHERE name = 'Alacena';
INSERT INTO stock_sub_categories (name, description, stock_category_id, display_order)
SELECT 'Harinas', 'Harinas y almidones', id, 3 FROM stock_categories WHERE name = 'Alacena';
INSERT INTO stock_sub_categories (name, description, stock_category_id, display_order)
SELECT 'Enlatados', 'Productos en lata', id, 4 FROM stock_categories WHERE name = 'Alacena';
INSERT INTO stock_sub_categories (name, description, stock_category_id, display_order)
SELECT 'Especias', 'Especias y condimentos', id, 5 FROM stock_categories WHERE name = 'Alacena';
INSERT INTO stock_sub_categories (name, description, stock_category_id, display_order)
SELECT 'Azúcares', 'Azúcar y endulzantes', id, 6 FROM stock_categories WHERE name = 'Alacena';

-- 7. Repostería
INSERT INTO stock_sub_categories (name, description, stock_category_id, display_order)
SELECT 'Postres', 'Pasteles y postres preparados', id, 1 FROM stock_categories WHERE name = 'Repostería';
INSERT INTO stock_sub_categories (name, description, stock_category_id, display_order)
SELECT 'Helados', 'Helados y toppings', id, 2 FROM stock_categories WHERE name = 'Repostería';
INSERT INTO stock_sub_categories (name, description, stock_category_id, display_order)
SELECT 'Chocolates', 'Chocolates y coberturas', id, 3 FROM stock_categories WHERE name = 'Repostería';

-- 8. Limpieza
INSERT INTO stock_sub_categories (name, description, stock_category_id, display_order)
SELECT 'Detergentes', 'Jabones y detergentes', id, 1 FROM stock_categories WHERE name = 'Limpieza';
INSERT INTO stock_sub_categories (name, description, stock_category_id, display_order)
SELECT 'Desinfectantes', 'Cloro y desinfectantes', id, 2 FROM stock_categories WHERE name = 'Limpieza';
INSERT INTO stock_sub_categories (name, description, stock_category_id, display_order)
SELECT 'Utensilios Limpieza', 'Esponjas, trapos, etc.', id, 3 FROM stock_categories WHERE name = 'Limpieza';

-- 9. Envases
INSERT INTO stock_sub_categories (name, description, stock_category_id, display_order)
SELECT 'Desechables', 'Platos y vasos desechables', id, 1 FROM stock_categories WHERE name = 'Envases';
INSERT INTO stock_sub_categories (name, description, stock_category_id, display_order)
SELECT 'Papel', 'Servilletas, papel, etc.', id, 2 FROM stock_categories WHERE name = 'Envases';
INSERT INTO stock_sub_categories (name, description, stock_category_id, display_order)
SELECT 'Bolsas', 'Bolsas para llevar', id, 3 FROM stock_categories WHERE name = 'Envases';

-- Insertar variantes de inventario por defecto
-- 1. Carnes - Res
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Filete', 'Filete de res', id FROM stock_sub_categories WHERE name = 'Res';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Bistec', 'Bistec de res', id FROM stock_sub_categories WHERE name = 'Res';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Molida', 'Carne molida de res', id FROM stock_sub_categories WHERE name = 'Res';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Costilla', 'Costilla de res', id FROM stock_sub_categories WHERE name = 'Res';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Arrachera', 'Arrachera de res', id FROM stock_sub_categories WHERE name = 'Res';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Pulpa', 'Pulpa de res', id FROM stock_sub_categories WHERE name = 'Res';

-- 1. Carnes - Pollo
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Pechuga', 'Pechuga de pollo', id FROM stock_sub_categories WHERE name = 'Pollo';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Muslo', 'Muslo de pollo', id FROM stock_sub_categories WHERE name = 'Pollo';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Ala', 'Ala de pollo', id FROM stock_sub_categories WHERE name = 'Pollo';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Pierna', 'Pierna de pollo', id FROM stock_sub_categories WHERE name = 'Pollo';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Cuarto', 'Cuarto de pollo', id FROM stock_sub_categories WHERE name = 'Pollo';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Entero', 'Pollo entero', id FROM stock_sub_categories WHERE name = 'Pollo';

-- 1. Carnes - Cerdo
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Costilla Cerdo', 'Costilla de cerdo', id FROM stock_sub_categories WHERE name = 'Cerdo';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Chuleta', 'Chuleta de cerdo', id FROM stock_sub_categories WHERE name = 'Cerdo';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Pierna Cerdo', 'Pierna de cerdo', id FROM stock_sub_categories WHERE name = 'Cerdo';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Lomo', 'Lomo de cerdo', id FROM stock_sub_categories WHERE name = 'Cerdo';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Posta', 'Posta de cerdo', id FROM stock_sub_categories WHERE name = 'Cerdo';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Chicharrón', 'Chicharrón de cerdo', id FROM stock_sub_categories WHERE name = 'Cerdo';

-- 1. Carnes - Pescado
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Tilapia', 'Tilapia fresca', id FROM stock_sub_categories WHERE name = 'Pescado';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Salmón', 'Salmón fresco', id FROM stock_sub_categories WHERE name = 'Pescado';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Mero', 'Mero fresco', id FROM stock_sub_categories WHERE name = 'Pescado';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Trucha', 'Trucha fresca', id FROM stock_sub_categories WHERE name = 'Pescado';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Corvina', 'Corvina fresca', id FROM stock_sub_categories WHERE name = 'Pescado';

-- 1. Carnes - Mariscos
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Camarón', 'Camarón fresco', id FROM stock_sub_categories WHERE name = 'Mariscos';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Pulpo', 'Pulpo fresco', id FROM stock_sub_categories WHERE name = 'Mariscos';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Calamar', 'Calamar fresco', id FROM stock_sub_categories WHERE name = 'Mariscos';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Mejillón', 'Mejillón fresco', id FROM stock_sub_categories WHERE name = 'Mariscos';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Almeja', 'Almeja fresca', id FROM stock_sub_categories WHERE name = 'Mariscos';

-- 1. Carnes - Embutidos
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Jamón', 'Jamón de cerdo', id FROM stock_sub_categories WHERE name = 'Embutidos';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Tocino', 'Tocino de cerdo', id FROM stock_sub_categories WHERE name = 'Embutidos';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Salchicha', 'Salchicha', id FROM stock_sub_categories WHERE name = 'Embutidos';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Chorizo', 'Chorizo', id FROM stock_sub_categories WHERE name = 'Embutidos';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Longaniza', 'Longaniza', id FROM stock_sub_categories WHERE name = 'Embutidos';

-- 2. Frutas y Verduras - Frutas
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Limón', 'Limón fresco', id FROM stock_sub_categories WHERE name = 'Frutas';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Naranja', 'Naranja fresca', id FROM stock_sub_categories WHERE name = 'Frutas';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Mango', 'Mango fresco', id FROM stock_sub_categories WHERE name = 'Frutas';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Piña', 'Piña fresca', id FROM stock_sub_categories WHERE name = 'Frutas';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Fresa', 'Fresa fresca', id FROM stock_sub_categories WHERE name = 'Frutas';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Aguacate', 'Aguacate fresco', id FROM stock_sub_categories WHERE name = 'Frutas';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Tomate', 'Tomate fresco', id FROM stock_sub_categories WHERE name = 'Frutas';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Plátano Maduro', 'Platano maduro fresco', id FROM stock_sub_categories WHERE name = 'Frutas';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Plátano Verde', 'Platano verde fresco', id FROM stock_sub_categories WHERE name = 'Frutas';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Banano', 'Banano fresco', id FROM stock_sub_categories WHERE name = 'Frutas';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Manzana', 'Manzana fresca', id FROM stock_sub_categories WHERE name = 'Frutas';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Pera', 'Pera fresca', id FROM stock_sub_categories WHERE name = 'Frutas';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Naranja', 'Naranja fresca', id FROM stock_sub_categories WHERE name = 'Frutas';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Uva', 'Uva fresca', id FROM stock_sub_categories WHERE name = 'Frutas';

-- 2. Frutas y Verduras - Verduras
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Cebolla', 'Cebolla fresca', id FROM stock_sub_categories WHERE name = 'Verduras';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Ajo', 'Ajo fresco', id FROM stock_sub_categories WHERE name = 'Verduras';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Chile', 'Chile fresco', id FROM stock_sub_categories WHERE name = 'Verduras';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Lechuga', 'Lechuga fresca', id FROM stock_sub_categories WHERE name = 'Verduras';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Zanahoria', 'Zanahoria fresca', id FROM stock_sub_categories WHERE name = 'Verduras';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Papa', 'Papa fresca', id FROM stock_sub_categories WHERE name = 'Verduras';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Pepino', 'Pepino fresco', id FROM stock_sub_categories WHERE name = 'Verduras';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Brócoli', 'Brócoli fresco', id FROM stock_sub_categories WHERE name = 'Verduras';

-- 2. Frutas y Verduras - Hierbas
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Cilantro', 'Cilantro fresco', id FROM stock_sub_categories WHERE name = 'Hierbas';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Perejil', 'Perejil fresco', id FROM stock_sub_categories WHERE name = 'Hierbas';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Orégano', 'Orégano fresco', id FROM stock_sub_categories WHERE name = 'Hierbas';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Romero', 'Romero fresco', id FROM stock_sub_categories WHERE name = 'Hierbas';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Albahaca', 'Albahaca fresca', id FROM stock_sub_categories WHERE name = 'Hierbas';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Menta', 'Menta fresca', id FROM stock_sub_categories WHERE name = 'Hierbas';

-- 2. Frutas y Verduras - Hongos
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Champiñón', 'Champiñón fresco', id FROM stock_sub_categories WHERE name = 'Hongos';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Portobello', 'Portobello fresco', id FROM stock_sub_categories WHERE name = 'Hongos';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Shiitake', 'Shiitake fresco', id FROM stock_sub_categories WHERE name = 'Hongos';

-- 3. Bebidas - Refrescos
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Coca-Cola', 'Coca-Cola', id FROM stock_sub_categories WHERE name = 'Refrescos';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Sprite', 'Sprite', id FROM stock_sub_categories WHERE name = 'Refrescos';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Fanta', 'Fanta naranja', id FROM stock_sub_categories WHERE name = 'Refrescos';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Fresca', 'Fresca toronja', id FROM stock_sub_categories WHERE name = 'Refrescos';

-- 3. Bebidas - Jugos
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Jugo Naranja', 'Jugo de naranja', id FROM stock_sub_categories WHERE name = 'Jugos';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Jugo Manzana', 'Jugo de manzana', id FROM stock_sub_categories WHERE name = 'Jugos';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Jugo Uva', 'Jugo de uva', id FROM stock_sub_categories WHERE name = 'Jugos';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Jugo Toronja', 'Jugo de toronja', id FROM stock_sub_categories WHERE name = 'Jugos';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Jugo Arándano', 'Jugo de arándano', id FROM stock_sub_categories WHERE name = 'Jugos';

-- 3. Bebidas - Cervezas (Costa Rican beers)
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Imperial', 'Cerveza Imperial', id FROM stock_sub_categories WHERE name = 'Cervezas';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Pilsen', 'Cerveza Pilsen', id FROM stock_sub_categories WHERE name = 'Cervezas';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Bavaria', 'Cerveza Bavaria', id FROM stock_sub_categories WHERE name = 'Cervezas';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Rock Ice', 'Cerveza Rock Ice', id FROM stock_sub_categories WHERE name = 'Cervezas';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Libertas', 'Cerveza Libertas', id FROM stock_sub_categories WHERE name = 'Cervezas';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Corona', 'Cerveza Corona', id FROM stock_sub_categories WHERE name = 'Cervezas';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Heineken', 'Cerveza Heineken', id FROM stock_sub_categories WHERE name = 'Cervezas';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Stella Artois', 'Cerveza Stella Artois', id FROM stock_sub_categories WHERE name = 'Cervezas';

-- 3. Bebidas - Vinos
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Vino Tinto Casa', 'Vino tinto de la casa', id FROM stock_sub_categories WHERE name = 'Vinos';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Vino Blanco Casa', 'Vino blanco de la casa', id FROM stock_sub_categories WHERE name = 'Vinos';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Vino Rosado', 'Vino rosado', id FROM stock_sub_categories WHERE name = 'Vinos';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Vino Espumoso', 'Vino espumoso', id FROM stock_sub_categories WHERE name = 'Vinos';

-- 3. Bebidas - Licores
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Tequila', 'Tequila', id FROM stock_sub_categories WHERE name = 'Licores';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Ron', 'Ron', id FROM stock_sub_categories WHERE name = 'Licores';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Vodka', 'Vodka', id FROM stock_sub_categories WHERE name = 'Licores';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Whisky', 'Whisky', id FROM stock_sub_categories WHERE name = 'Licores';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Ginebra', 'Ginebra', id FROM stock_sub_categories WHERE name = 'Licores';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Guaro', 'Guaro Cacique', id FROM stock_sub_categories WHERE name = 'Licores';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Brandy', 'Brandy', id FROM stock_sub_categories WHERE name = 'Licores';

-- 3. Bebidas - Café y Té
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Café Molido', 'Café molido', id FROM stock_sub_categories WHERE name = 'Café y Té';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Café Grano', 'Café en grano', id FROM stock_sub_categories WHERE name = 'Café y Té';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Té Negro', 'Té negro', id FROM stock_sub_categories WHERE name = 'Café y Té';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Té Verde', 'Té verde', id FROM stock_sub_categories WHERE name = 'Café y Té';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Manzanilla', 'Té de manzanilla', id FROM stock_sub_categories WHERE name = 'Café y Té';

-- 3. Bebidas - Lácteos
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Leche Entera', 'Leche entera', id FROM stock_sub_categories WHERE name = 'Lácteos';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Leche Deslactosada', 'Leche deslactosada', id FROM stock_sub_categories WHERE name = 'Lácteos';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Leche Semi Descremada', 'Leche semi descremada', id FROM stock_sub_categories WHERE name = 'Lácteos';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Natilla', 'Natilla', id FROM stock_sub_categories WHERE name = 'Lácteos';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Queso Crema', 'Queso crema', id FROM stock_sub_categories WHERE name = 'Lácteos';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Queso', 'Queso', id FROM stock_sub_categories WHERE name = 'Lácteos';

-- 3. Bebidas - Agua
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Agua Natural', 'Agua natural', id FROM stock_sub_categories WHERE name = 'Agua';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Agua Mineral', 'Agua mineral con gas', id FROM stock_sub_categories WHERE name = 'Agua';

-- 4. Salsas - Salsas Picantes
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Salsa Lizano', 'Salsa Lizano', id FROM stock_sub_categories WHERE name = 'Salsas Picantes';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Tabasco', 'Salsa Tabasco', id FROM stock_sub_categories WHERE name = 'Salsas Picantes';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Sriracha', 'Salsa Sriracha', id FROM stock_sub_categories WHERE name = 'Salsas Picantes';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Chilero', 'Chilero', id FROM stock_sub_categories WHERE name = 'Salsas Picantes';

-- 4. Salsas - Aderezos
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Mayonesa', 'Mayonesa', id FROM stock_sub_categories WHERE name = 'Aderezos';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Mostaza', 'Mostaza', id FROM stock_sub_categories WHERE name = 'Aderezos';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Ketchup', 'Ketchup', id FROM stock_sub_categories WHERE name = 'Aderezos';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Ranch', 'Aderezo Ranch', id FROM stock_sub_categories WHERE name = 'Aderezos';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Vinagreta', 'Vinagreta', id FROM stock_sub_categories WHERE name = 'Aderezos';

-- 4. Salsas - Salsas Base
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Salsa Soya', 'Salsa de soya', id FROM stock_sub_categories WHERE name = 'Salsas Base';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Salsa Inglesa', 'Salsa inglesa', id FROM stock_sub_categories WHERE name = 'Salsas Base';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Salsa Teriyaki', 'Salsa teriyaki', id FROM stock_sub_categories WHERE name = 'Salsas Base';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Salsa BBQ', 'Salsa BBQ', id FROM stock_sub_categories WHERE name = 'Salsas Base';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Salsa Búfalo', 'Salsa búfalo', id FROM stock_sub_categories WHERE name = 'Salsas Base';

-- 4. Salsas - Aceites y Vinagres
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Aceite Oliva', 'Aceite de oliva', id FROM stock_sub_categories WHERE name = 'Aceites y Vinagres';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Aceite Vegetal', 'Aceite vegetal', id FROM stock_sub_categories WHERE name = 'Aceites y Vinagres';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Vinagre Blanco', 'Vinagre blanco', id FROM stock_sub_categories WHERE name = 'Aceites y Vinagres';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Vinagre Balsámico', 'Vinagre balsámico', id FROM stock_sub_categories WHERE name = 'Aceites y Vinagres';

-- 5. Congelados - Postres Congelados
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Helado Vainilla', 'Helado de vainilla', id FROM stock_sub_categories WHERE name = 'Postres Congelados';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Helado Chocolate', 'Helado de chocolate', id FROM stock_sub_categories WHERE name = 'Postres Congelados';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Helado Fresa', 'Helado de fresa', id FROM stock_sub_categories WHERE name = 'Postres Congelados';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Paletas', 'Paletas de hielo', id FROM stock_sub_categories WHERE name = 'Postres Congelados';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Churros Congelados', 'Churros congelados', id FROM stock_sub_categories WHERE name = 'Postres Congelados';

-- 5. Congelados - Precosidos
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Papas Fritas', 'Papas fritas congeladas', id FROM stock_sub_categories WHERE name = 'Precosidos';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Nuggets', 'Nuggets de pollo', id FROM stock_sub_categories WHERE name = 'Precosidos';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Dedos de Pollo', 'Dedos de pollo empanizados', id FROM stock_sub_categories WHERE name = 'Precosidos';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Dedos de Queso', 'Dedos de queso', id FROM stock_sub_categories WHERE name = 'Precosidos';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Aros de Cebolla', 'Aros de cebolla congelados', id FROM stock_sub_categories WHERE name = 'Precosidos';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Jalapeño Poppers', 'Jalapeño poppers', id FROM stock_sub_categories WHERE name = 'Precosidos';

-- 5. Congelados - Hielo
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Hielo Cubo', 'Hielo en cubos', id FROM stock_sub_categories WHERE name = 'Hielo';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Hielo Frappe', 'Hielo para frappe', id FROM stock_sub_categories WHERE name = 'Hielo';

-- 6. Alacena - Granos
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Arroz Blanco', 'Arroz blanco', id FROM stock_sub_categories WHERE name = 'Granos';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Arroz Integral', 'Arroz integral', id FROM stock_sub_categories WHERE name = 'Granos';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Frijol Negro', 'Frijol negro', id FROM stock_sub_categories WHERE name = 'Granos';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Frijol Rojo', 'Frijol rojo', id FROM stock_sub_categories WHERE name = 'Granos';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Lenteja', 'Lenteja', id FROM stock_sub_categories WHERE name = 'Granos';

-- 6. Alacena - Pastas
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Spaghetti', 'Pasta spaghetti', id FROM stock_sub_categories WHERE name = 'Pastas';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Penne', 'Pasta penne', id FROM stock_sub_categories WHERE name = 'Pastas';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Macarrón', 'Pasta macarrón', id FROM stock_sub_categories WHERE name = 'Pastas';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Lasaña', 'Pasta para lasaña', id FROM stock_sub_categories WHERE name = 'Pastas';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Fideo', 'Fideo', id FROM stock_sub_categories WHERE name = 'Pastas';

-- 6. Alacena - Harinas
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Harina Trigo', 'Harina de trigo', id FROM stock_sub_categories WHERE name = 'Harinas';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Maicena', 'Maicena', id FROM stock_sub_categories WHERE name = 'Harinas';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Pan Molido', 'Pan molido', id FROM stock_sub_categories WHERE name = 'Harinas';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Harina Hotcakes', 'Harina para hotcakes', id FROM stock_sub_categories WHERE name = 'Harinas';

-- 6. Alacena - Enlatados
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Atún', 'Atún enlatado', id FROM stock_sub_categories WHERE name = 'Enlatados';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Sardina', 'Sardina enlatada', id FROM stock_sub_categories WHERE name = 'Enlatados';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Elote Enlatado', 'Elote enlatado', id FROM stock_sub_categories WHERE name = 'Enlatados';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Champiñones Enlatados', 'Champiñones enlatados', id FROM stock_sub_categories WHERE name = 'Enlatados';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Chiles Enlatados', 'Chiles enlatados', id FROM stock_sub_categories WHERE name = 'Enlatados';

-- 6. Alacena - Especias
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Sal', 'Sal de mesa', id FROM stock_sub_categories WHERE name = 'Especias';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Pimienta', 'Pimienta negra', id FROM stock_sub_categories WHERE name = 'Especias';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Comino', 'Comino', id FROM stock_sub_categories WHERE name = 'Especias';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Paprika', 'Paprika', id FROM stock_sub_categories WHERE name = 'Especias';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Canela', 'Canela', id FROM stock_sub_categories WHERE name = 'Especias';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Consomé', 'Consomé de pollo', id FROM stock_sub_categories WHERE name = 'Especias';

-- 6. Alacena - Azúcares
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Azúcar', 'Azúcar blanca', id FROM stock_sub_categories WHERE name = 'Azúcares';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Miel', 'Miel de abeja', id FROM stock_sub_categories WHERE name = 'Azúcares';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Jarabe', 'Jarabe de maple', id FROM stock_sub_categories WHERE name = 'Azúcares';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Splenda', 'Endulzante Splenda', id FROM stock_sub_categories WHERE name = 'Azúcares';

-- 7. Repostería - Postres
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Pastel Chocolate', 'Pastel de chocolate', id FROM stock_sub_categories WHERE name = 'Postres';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Pastel Tres Leches', 'Pastel tres leches', id FROM stock_sub_categories WHERE name = 'Postres';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Cheesecake', 'Cheesecake', id FROM stock_sub_categories WHERE name = 'Postres';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Flan', 'Flan', id FROM stock_sub_categories WHERE name = 'Postres';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Churros', 'Churros', id FROM stock_sub_categories WHERE name = 'Postres';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Brownie', 'Brownie', id FROM stock_sub_categories WHERE name = 'Postres';

-- 7. Repostería - Helados
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Helado Vainilla Repos', 'Helado de vainilla', id FROM stock_sub_categories WHERE name = 'Helados';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Helado Chocolate Repos', 'Helado de chocolate', id FROM stock_sub_categories WHERE name = 'Helados';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Helado Fresa Repos', 'Helado de fresa', id FROM stock_sub_categories WHERE name = 'Helados';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Helado Cookies', 'Helado cookies and cream', id FROM stock_sub_categories WHERE name = 'Helados';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Sorbete', 'Sorbete de frutas', id FROM stock_sub_categories WHERE name = 'Helados';

-- 7. Repostería - Chocolates
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Chocolate Amargo', 'Chocolate amargo', id FROM stock_sub_categories WHERE name = 'Chocolates';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Chocolate Leche', 'Chocolate con leche', id FROM stock_sub_categories WHERE name = 'Chocolates';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Cocoa', 'Cocoa en polvo', id FROM stock_sub_categories WHERE name = 'Chocolates';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Nutella', 'Nutella', id FROM stock_sub_categories WHERE name = 'Chocolates';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Chispas Chocolate', 'Chispas de chocolate', id FROM stock_sub_categories WHERE name = 'Chocolates';

-- 8. Limpieza - Detergentes
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Lavatrastes', 'Jabón lavatrastes', id FROM stock_sub_categories WHERE name = 'Detergentes';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Jabón Manos', 'Jabón para manos', id FROM stock_sub_categories WHERE name = 'Detergentes';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Detergente Multiusos', 'Detergente multiusos', id FROM stock_sub_categories WHERE name = 'Detergentes';

-- 8. Limpieza - Desinfectantes
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Cloro', 'Cloro', id FROM stock_sub_categories WHERE name = 'Desinfectantes';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Pino', 'Desinfectante de pino', id FROM stock_sub_categories WHERE name = 'Desinfectantes';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Antibacterial', 'Antibacterial', id FROM stock_sub_categories WHERE name = 'Desinfectantes';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Sanitizante', 'Sanitizante', id FROM stock_sub_categories WHERE name = 'Desinfectantes';

-- 8. Limpieza - Utensilios Limpieza
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Esponja', 'Esponja para trastes', id FROM stock_sub_categories WHERE name = 'Utensilios Limpieza';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Fibra', 'Fibra para limpiar', id FROM stock_sub_categories WHERE name = 'Utensilios Limpieza';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Trapo', 'Trapo de limpieza', id FROM stock_sub_categories WHERE name = 'Utensilios Limpieza';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Guantes', 'Guantes de limpieza', id FROM stock_sub_categories WHERE name = 'Utensilios Limpieza';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Bolsa Basura', 'Bolsa para basura', id FROM stock_sub_categories WHERE name = 'Utensilios Limpieza';

-- 9. Envases - Desechables
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Plato Foam', 'Plato de foam', id FROM stock_sub_categories WHERE name = 'Desechables';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Vaso Foam', 'Vaso de foam', id FROM stock_sub_categories WHERE name = 'Desechables';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Contenedor', 'Contenedor para llevar', id FROM stock_sub_categories WHERE name = 'Desechables';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Cubiertos Desechables', 'Cubiertos desechables', id FROM stock_sub_categories WHERE name = 'Desechables';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Popote', 'Popote', id FROM stock_sub_categories WHERE name = 'Desechables';

-- 9. Envases - Papel
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Servilleta', 'Servilletas', id FROM stock_sub_categories WHERE name = 'Papel';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Papel Aluminio', 'Papel aluminio', id FROM stock_sub_categories WHERE name = 'Papel';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Papel Encerado', 'Papel encerado', id FROM stock_sub_categories WHERE name = 'Papel';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Film Plástico', 'Film plástico', id FROM stock_sub_categories WHERE name = 'Papel';

-- 9. Envases - Bolsas
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Bolsa Llevar', 'Bolsa para llevar', id FROM stock_sub_categories WHERE name = 'Bolsas';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Bolsa Basura Grande', 'Bolsa de basura grande', id FROM stock_sub_categories WHERE name = 'Bolsas';
INSERT INTO stock_variants (name, description, stock_sub_category_id)
SELECT 'Bolsa Basura Chica', 'Bolsa de basura chica', id FROM stock_sub_categories WHERE name = 'Bolsas';

-- Insertar mesas de ejemplo
INSERT INTO tables (table_number, capacity, status) VALUES
('M1', 2, 'available'),
('M2', 2, 'available'),
('M3', 4, 'available'),
('M4', 4, 'available'),
('M5', 4, 'available'),
('M6', 6, 'available'),
('M7', 6, 'available'),
('M8', 8, 'available'),
('BARRA1', 1, 'available'),
('BARRA2', 1, 'available'),
('BARRA3', 1, 'available'),
('BARRA4', 1, 'available');

-- Tables indexes
CREATE INDEX idx_tables_status ON tables(status);

-- Customers indexes
CREATE INDEX idx_customers_customer_id ON customers(customer_id);
CREATE INDEX idx_customers_phone ON customers(phone);

-- Sub Menus indexes
CREATE INDEX idx_menu_sub_categories_category ON menu_sub_categories(category_id);
CREATE INDEX idx_menu_sub_categories_active ON menu_sub_categories(is_active);

-- Menu Variants indexes
CREATE INDEX idx_menu_variants_sub_category ON menu_variants(sub_category_id);
CREATE INDEX idx_menu_variants_available ON menu_variants(is_available);

-- Orders indexes
CREATE INDEX idx_orders_table ON orders(table_id);
CREATE INDEX idx_orders_customer ON orders(customer_id);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_payment_status ON orders(payment_status);
CREATE INDEX idx_orders_created_at ON orders(created_at);

-- Order Items indexes
CREATE INDEX idx_order_items_order ON order_items(order_id);
CREATE INDEX idx_order_items_menu_variant ON order_items(menu_variant_id);
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
CREATE INDEX idx_karaoke_requests_created_at ON karaoke_song_requests(created_at);

-- Sessions indexes
CREATE INDEX idx_sessions_token ON sessions(token);

-- Promotions indexes
CREATE INDEX idx_promotions_type ON promotions(discount_type);
CREATE INDEX idx_promotions_active ON promotions(is_active);
CREATE INDEX idx_promotions_dates ON promotions(valid_from, valid_until);

-- Settings indexes
CREATE INDEX idx_settings_service ON settings(service);
CREATE INDEX idx_settings_key ON settings(key);

-- =============================================================================
-- TRIGGERS FOR AUTOMATIC UPDATED_AT
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

CREATE TRIGGER update_staff_updated_at BEFORE UPDATE ON staff
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_menu_categories_updated_at BEFORE UPDATE ON menu_categories
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_menu_sub_categories_updated_at BEFORE UPDATE ON menu_sub_categories
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_menu_variants_updated_at BEFORE UPDATE ON menu_variants
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_menu_ingredients_updated_at BEFORE UPDATE ON menu_ingredients
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_stock_categories_updated_at BEFORE UPDATE ON stock_categories
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_stock_sub_categories_updated_at BEFORE UPDATE ON stock_sub_categories
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_stock_variants_updated_at BEFORE UPDATE ON stock_variants
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_stock_count_updated_at BEFORE UPDATE ON stock_count
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_suppliers_updated_at BEFORE UPDATE ON suppliers
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_orders_updated_at BEFORE UPDATE ON orders
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_order_items_updated_at BEFORE UPDATE ON order_items
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_payments_updated_at BEFORE UPDATE ON payments
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_outcome_invoices_updated_at BEFORE UPDATE ON outcome_invoices
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_income_invoices_updated_at BEFORE UPDATE ON income_invoices
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_invoice_items_updated_at BEFORE UPDATE ON invoice_items
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_customer_favorites_updated_at BEFORE UPDATE ON customer_favorites
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_request_notifications_updated_at BEFORE UPDATE ON request_notifications
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_karaoke_song_requests_updated_at BEFORE UPDATE ON karaoke_song_requests
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_karaoke_song_library_updated_at BEFORE UPDATE ON karaoke_song_library
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_loyalty_points_transactions_updated_at BEFORE UPDATE ON loyalty_points_transactions
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_table_sessions_updated_at BEFORE UPDATE ON table_sessions
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_reservations_updated_at BEFORE UPDATE ON reservations
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_reviews_updated_at BEFORE UPDATE ON reviews
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_promotions_updated_at BEFORE UPDATE ON promotions
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Note: invoice_items.invoice_id can reference either income_invoices or outcome_invoices
-- based on the invoice_type field. Foreign key constraints are handled in application code.

-- =============================================================================
-- END OF SCHEMA
-- =============================================================================
