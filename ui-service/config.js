// Bar Restaurant UI Configuration
// This file contains service URLs and API configuration

console.log('🔧 CONFIG.JS: Starting to load...');

// Environment detection and service URL configuration
function getServiceUrls() {
    const isLocalDevelopment = window.location.hostname === 'localhost' ||
                              window.location.hostname === '127.0.0.1' ||
                              window.location.hostname.includes('localhost');

    if (isLocalDevelopment) {
        console.log('🔧 Detected local development environment - using localhost URLs');
        return {
            gateway: 'http://localhost:8082'
        };
    } else {
        console.log('🔧 Detected production environment - using Docker service names');
        return {
            gateway: 'http://barrest_gateway_service:8082'
        };
    }
}

const SERVICE_URLS = getServiceUrls();

// Configuration object with all service URLs and health endpoints
const CONFIG = {
    GATEWAY_URL: SERVICE_URLS.gateway,
    API: {
        gateway: SERVICE_URLS.gateway + '/api/v1',
        LOGIN: '/api/v1/sessions/p/login',
        LOGOUT: '/api/v1/sessions/logout',
        VALIDATE: '/api/v1/sessions/p/validate'
    },
    SERVICES: {
        gateway: SERVICE_URLS.gateway + '/api/v1/gateway/p/health',
        session: SERVICE_URLS.gateway + '/api/v1/sessions/p/health',
        menu: SERVICE_URLS.gateway + '/api/v1/menu/p/health',
        data: SERVICE_URLS.gateway + '/api/v1/data/p/health',
        inventory: SERVICE_URLS.gateway + '/api/v1/inventory/p/health',
        invoice: SERVICE_URLS.gateway + '/api/v1/invoices/p/health'
    },
    AUTH: {
        login: SERVICE_URLS.gateway + '/api/v1/sessions/p/login',
        logout: SERVICE_URLS.gateway + '/api/v1/sessions/logout',
        validate: SERVICE_URLS.gateway + '/api/v1/sessions/p/validate',
        TOKEN_KEY: 'barrest_token',
        SESSION_ID_KEY: 'barrest_session_id', // Deprecated - kept for backward compatibility
        USER_KEY: 'barrest_user_data',
        REMEMBER_KEY: 'barrest_remember_me'
    },
    MENU: {
        // Menu Categories (top level: Drinks, Desserts, etc.)
        categories: SERVICE_URLS.gateway + '/api/v1/menu/categories',
        // Menu Sub-Categories (second level: Smoothies, Sodas, etc. - grouped by category)
        subCategories: SERVICE_URLS.gateway + '/api/v1/menu/sub-categories',
        // Menu Variants (third level: Banana Smoothie, Pineapple Smoothie, etc. - with pricing)
        variants: SERVICE_URLS.gateway + '/api/v1/menu/variants',
        // Menu Ingredients
        ingredients: SERVICE_URLS.gateway + '/api/v1/menu/ingredients',
        ingredientsByVariant: (variantId) => SERVICE_URLS.gateway + `/api/v1/menu/variants/${variantId}/ingredients`,
        // Stock Categories
        stockCategories: SERVICE_URLS.gateway + '/api/v1/inventory/categories',
        // Stock Sub-Categories
        stockSubCategories: SERVICE_URLS.gateway + '/api/v1/inventory/sub-categories',
        // Stock Variants
        stockVariants: SERVICE_URLS.gateway + '/api/v1/inventory/variants'
    },
    INVENTORY: {
        // Suppliers
        suppliers: SERVICE_URLS.gateway + '/api/v1/inventory/suppliers',
        supplier: (id) => SERVICE_URLS.gateway + `/api/v1/inventory/suppliers/${id}`
    },
    INVOICE: {
        // Outcome Invoices (supplier purchases)
        outcomeInvoices: SERVICE_URLS.gateway + '/api/v1/invoices/outcome',
        outcomeInvoice: (id) => SERVICE_URLS.gateway + `/api/v1/invoices/outcome/${id}`,
        // Income Invoices (customer billing)
        incomeInvoices: SERVICE_URLS.gateway + '/api/v1/invoices/income',
        incomeInvoice: (id) => SERVICE_URLS.gateway + `/api/v1/invoices/income/${id}`,
        // Invoice Items (line items)
        invoiceItems: (invoiceId) => SERVICE_URLS.gateway + `/api/v1/invoices/outcome/${invoiceId}/items`,
        invoiceItem: (invoiceId, itemId) => SERVICE_URLS.gateway + `/api/v1/invoices/outcome/${invoiceId}/items/${itemId}`
    }
};

console.log('🔧 Configuration loaded (v1.2):', {
    gateway: SERVICE_URLS.gateway,
    services: CONFIG.SERVICES,
    menu: CONFIG.MENU
});

// Export for global access
window.CONFIG = CONFIG;
