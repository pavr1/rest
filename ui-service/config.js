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
        data: SERVICE_URLS.gateway + '/api/v1/data/p/health'
    },
    AUTH: {
        login: SERVICE_URLS.gateway + '/api/v1/sessions/p/login',
        logout: SERVICE_URLS.gateway + '/api/v1/sessions/logout',
        validate: SERVICE_URLS.gateway + '/api/v1/sessions/p/validate',
        SESSION_ID_KEY: 'barrest_session_id',
        USER_KEY: 'barrest_user_data',
        REMEMBER_KEY: 'barrest_remember_me'
    },
    MENU: {
        // Menu Categories (top level: Drinks, Desserts, etc.)
        categories: SERVICE_URLS.gateway + '/api/v1/menu/categories',
        // Sub Menus (second level: Smoothies, Sodas, etc. - grouped by category)
        subMenus: SERVICE_URLS.gateway + '/api/v1/menu/submenus',
        // Menu Items (third level: Banana Smoothie, Pineapple Smoothie, etc. - with pricing)
        items: SERVICE_URLS.gateway + '/api/v1/menu/items',
        // Menu Item Ingredients (use with item ID: `${CONFIG.MENU.items}/${itemId}/ingredients`)
        ingredients: (itemId) => SERVICE_URLS.gateway + `/api/v1/menu/items/${itemId}/ingredients`,
        // Menu Item Cost
        itemCost: (itemId) => SERVICE_URLS.gateway + `/api/v1/menu/items/${itemId}/cost`,
        // Stock Item Categories
        stockCategories: SERVICE_URLS.gateway + '/api/v1/stock/categories',
        // Stock Items
        stockItems: SERVICE_URLS.gateway + '/api/v1/stock/items'
    }
};

console.log('🔧 Configuration loaded:', {
    gateway: SERVICE_URLS.gateway,
    menu: CONFIG.MENU
});

// Export for global access
window.CONFIG = CONFIG;
