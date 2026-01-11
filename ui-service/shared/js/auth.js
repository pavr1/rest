// === BAR RESTAURANT - AUTHENTICATION SERVICE ===

class AuthService {
    constructor() {
        // Ensure CONFIG is available
        if (typeof CONFIG === 'undefined') {
            console.error('❌ CONFIG not available when creating AuthService');
            throw new Error('CONFIG must be loaded before AuthService');
        }

        // Use the gateway URL for authentication
        this.baseURL = CONFIG.GATEWAY_URL;
        this.tokenKey = CONFIG.AUTH.TOKEN_KEY || CONFIG.AUTH.SESSION_ID_KEY; // Fallback for backward compatibility
        this.userKey = CONFIG.AUTH.USER_KEY;
        this.rememberKey = CONFIG.AUTH.REMEMBER_KEY;

        console.log('✅ AuthService initialized');
    }

    // === MAIN LOGIN METHOD ===
    
    async login(username, password, rememberMe = false) {
        try {
            console.log('🔑 Attempting login for:', username);
            
            const response = await fetch(`${this.baseURL}${CONFIG.API.LOGIN}`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ username, password })
            });

            const { success, result } = await isSuccessfulResponse(response);
            
            if (success) {
                // Store authentication data (token instead of session_id)
                this.setToken(result.data.token, rememberMe);
                this.setUserData(result.data.user, result.data.role, result.data.permissions || []);
                
                console.log('✅ Login successful for:', username);
                
                return {
                    success: true,
                    user: result.data.user,
                    role: result.data.role,
                    permissions: result.data.permissions || []
                };
            } else {
                const errorMessage = result.message || `Login failed (${result.code})`;
                throw new Error(errorMessage);
            }
            
        } catch (error) {
            console.error('❌ Login failed:', error);
            throw error;
        }
    }

    // === LOGOUT METHOD ===
    
    async logout() {
        try {
            console.log('🚪 Attempting logout...');
            
            const token = this.getToken();
            if (!token) {
                this.clearAuthData();
                return { success: true };
            }

            const response = await fetch(`${this.baseURL}${CONFIG.API.LOGOUT}`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
                body: JSON.stringify({ token: token })
            });

            // Always clear local data regardless of server response
            this.clearAuthData();
            
            console.log('✅ Logout successful');
            return { success: true };
            
        } catch (error) {
            console.error('❌ Logout error:', error);
            // Always clear local data even if request fails
            this.clearAuthData();
            return { success: true, warning: 'Server logout failed' };
        }
    }

    // === SESSION VALIDATION ===
    
    async validateSession() {
        try {
            const token = this.getToken();
            if (!token) {
                return false;
            }

            const response = await fetch(`${this.baseURL}${CONFIG.API.VALIDATE}`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ token: token })
            });

            if (response.ok) {
                const data = await response.json();
                return data.data?.valid || data.valid || false;
            }
            
            return false;
        } catch (error) {
            console.error('❌ Session validation error:', error);
            return false;
        }
    }

    // === SESSION MANAGEMENT ===
    
    setToken(token, rememberMe = false) {
        const storage = rememberMe ? localStorage : sessionStorage;
        storage.setItem(this.tokenKey, token);
        
        if (rememberMe) {
            localStorage.setItem(this.rememberKey, 'true');
        } else {
            localStorage.removeItem(this.rememberKey);
        }
        
        console.log('💾 Token stored');
    }

    getToken() {
        // Check sessionStorage first, then localStorage
        let token = sessionStorage.getItem(this.tokenKey);
        if (!token) {
            token = localStorage.getItem(this.tokenKey);
        }
        return token;
    }

    // Backward compatibility methods (deprecated - use getToken/setToken)
    setSessionId(sessionId, rememberMe = false) {
        console.warn('⚠️ setSessionId is deprecated, use setToken instead');
        this.setToken(sessionId, rememberMe);
    }

    getSessionId() {
        console.warn('⚠️ getSessionId is deprecated, use getToken instead');
        return this.getToken();
    }

    setUserData(user, role, permissions) {
        const userData = { user, role, permissions };
        const storage = this.isRememberMe() ? localStorage : sessionStorage;
        storage.setItem(this.userKey, JSON.stringify(userData));
        console.log('💾 User data stored:', { username: user?.username, role: role?.name });
    }

    getUserData() {
        const storage = this.isRememberMe() ? localStorage : sessionStorage;
        const data = storage.getItem(this.userKey);
        return data ? JSON.parse(data) : null;
    }

    isRememberMe() {
        return localStorage.getItem(this.rememberKey) === 'true';
    }

    clearAuthData() {
        sessionStorage.removeItem(this.tokenKey);
        sessionStorage.removeItem(this.userKey);
        localStorage.removeItem(this.tokenKey);
        localStorage.removeItem(this.userKey);
        localStorage.removeItem(this.rememberKey);
        console.log('🧹 Auth data cleared');
    }

    isAuthenticated() {
        const token = this.getToken();
        if (!token) {
            return false;
        }
        return token.length > 0 && token !== 'null' && token !== 'undefined';
    }

    getCurrentUser() {
        const userData = this.getUserData();
        return userData?.user || null;
    }

    getCurrentRole() {
        const userData = this.getUserData();
        return userData?.role || null;
    }

    getPermissions() {
        const userData = this.getUserData();
        return userData?.permissions || [];
    }

    hasPermission(permission) {
        const permissions = this.getPermissions();
        return permissions.includes(permission);
    }
}

// === GLOBAL AUTH SERVICE INITIALIZATION ===

function initializeAuthService() {
    try {
        if (window.authService) {
            return window.authService;
        }
        
        const authService = new AuthService();
        window.authService = authService;
        
        return authService;
        
    } catch (error) {
        console.error('❌ Failed to initialize AuthService:', error);
        throw error;
    }
}

// === AUTHENTICATED REQUEST HELPER ===

async function makeAuthenticatedRequest(url, options = {}) {
    const authService = window.authService;
    
    if (!authService) {
        throw new Error('Authentication service not available');
    }
    
    if (!authService.isAuthenticated()) {
        console.warn('⚠️ User not authenticated');
        window.location.href = 'login.html';
        throw new Error('User not authenticated');
    }
    
    const token = authService.getToken();
    
    const headers = {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`,
        ...options.headers
    };
    
    const response = await fetch(url, {
        ...options,
        headers
    });
    
    if (response.status === 401) {
        console.warn('⚠️ Session expired');
        authService.clearAuthData();
        window.location.href = 'login.html';
        throw new Error('Session expired');
    }
    
    return response;
}

// Export for global access
window.AuthService = AuthService;
window.initializeAuthService = initializeAuthService;
window.makeAuthenticatedRequest = makeAuthenticatedRequest;

// Auto-initialize when script loads
document.addEventListener('DOMContentLoaded', () => {
    try {
        initializeAuthService();
    } catch (error) {
        console.warn('⚠️ AuthService will be initialized later');
    }
});

