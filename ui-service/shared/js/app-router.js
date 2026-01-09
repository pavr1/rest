// === BAR RESTAURANT - APP ROUTER ===
// Single Page Application router for dynamic content loading

class AppRouter {
    constructor() {
        this.currentPage = null;
        this.pageContent = document.getElementById('pageContent');
        this.pageTitle = document.getElementById('pageTitle');
        this.user = null;
        
        // Page configuration - organized by section
        this.pages = {
            // Dashboard
            'dashboard': { title: 'Dashboard', icon: 'fa-home', partial: 'pages/dashboard.html' },
            
            // Orders
            'orders-active': { title: 'Active Orders', icon: 'fa-clipboard-list', partial: 'pages/orders/active.html' },
            'orders-history': { title: 'Order History', icon: 'fa-history', partial: 'pages/orders/history.html' },
            
            // Tables
            'tables-map': { title: 'Table Map', icon: 'fa-map', partial: 'pages/tables/map.html' },
            'reservations': { title: 'Reservations', icon: 'fa-calendar-alt', partial: 'pages/tables/reservations.html' },
            
            // Menu
            'menu-items': { title: 'Menu Items', icon: 'fa-hamburger', partial: 'pages/menu/items.html' },
            'menu-categories': { title: 'Menu Categories', icon: 'fa-tags', partial: 'pages/menu/categories.html' },
            
            // Inventory
            'stock-items': { title: 'Stock Items', icon: 'fa-boxes', partial: 'pages/inventory/items.html' },
            'stock-categories': { title: 'Stock Categories', icon: 'fa-layer-group', partial: 'pages/inventory/categories.html' },
            
            // Admin
            'staff': { title: 'Staff', icon: 'fa-users', partial: 'pages/admin/staff.html' },
            'reports': { title: 'Reports', icon: 'fa-chart-bar', partial: 'pages/admin/reports.html' },
            'settings': { title: 'Settings', icon: 'fa-sliders-h', partial: 'pages/admin/settings.html' }
        };
        
        this.init();
    }
    
    async init() {
        console.log('🚀 Initializing App Router...');
        
        // Check authentication
        await this.checkAuth();
        
        // Setup user info
        this.setupUserInfo();
        
        // Setup event listeners
        this.setupEventListeners();
        
        // Load system status
        await this.loadSystemStatus();
        
        // Load initial page from URL hash or default to dashboard
        const initialPage = this.getPageFromHash() || 'dashboard';
        await this.navigateTo(initialPage);
        
        console.log('✅ App Router initialized');
    }
    
    async checkAuth() {
        if (!window.authService) {
            window.authService = new AuthService();
        }
        
        if (!window.authService.isAuthenticated()) {
            console.log('❌ Not authenticated, redirecting to login...');
            window.location.href = 'login.html';
            return;
        }
        
        const isValid = await window.authService.validateSession();
        if (!isValid) {
            console.log('❌ Session invalid, redirecting to login...');
            window.authService.clearAuthData();
            window.location.href = 'login.html';
            return;
        }
        
        this.user = window.authService.getUserData();
        console.log('✅ Authenticated as:', this.user?.user?.username);
    }
    
    setupUserInfo() {
        if (!this.user) return;
        
        const { user, role } = this.user;
        const firstName = user?.first_name || user?.username || 'User';
        const lastName = user?.last_name || '';
        const fullName = `${firstName} ${lastName}`.trim();
        const roleName = role?.name || role || 'Staff';
        
        document.getElementById('userName').textContent = fullName;
        document.getElementById('userRole').textContent = roleName;
        document.getElementById('userAvatar').textContent = firstName.charAt(0).toUpperCase();
    }
    
    setupEventListeners() {
        // Sidebar toggle (mobile)
        document.getElementById('sidebarToggle')?.addEventListener('click', () => {
            document.getElementById('sidebar').classList.toggle('show');
            document.getElementById('sidebarOverlay').classList.toggle('show');
        });
        
        // Close sidebar when clicking overlay
        document.getElementById('sidebarOverlay')?.addEventListener('click', () => {
            document.getElementById('sidebar').classList.remove('show');
            document.getElementById('sidebarOverlay').classList.remove('show');
        });
        
        // Brand click -> dashboard
        document.querySelector('.sidebar-brand')?.addEventListener('click', () => {
            this.navigateTo('dashboard');
        });
        
        // Collapsible menu toggles
        document.querySelectorAll('.nav-toggle').forEach(toggle => {
            toggle.addEventListener('click', () => {
                const targetId = toggle.dataset.target;
                const submenu = document.getElementById(targetId);
                
                toggle.classList.toggle('expanded');
                submenu.classList.toggle('expanded');
                
                this.saveMenuState();
            });
        });
        
        // Navigation links
        document.querySelectorAll('.nav-link[data-page]').forEach(link => {
            link.addEventListener('click', () => {
                const page = link.dataset.page;
                this.navigateTo(page);
                
                // Close mobile sidebar
                document.getElementById('sidebar').classList.remove('show');
                document.getElementById('sidebarOverlay').classList.remove('show');
            });
        });
        
        // Restore menu state
        this.restoreMenuState();
        
        // Logout
        document.getElementById('logoutBtn')?.addEventListener('click', async (e) => {
            e.preventDefault();
            await this.handleLogout();
        });
        
        // Theme options
        document.querySelectorAll('.theme-option').forEach(option => {
            option.addEventListener('click', (e) => {
                e.preventDefault();
                const theme = e.target.dataset.theme;
                if (window.themeManager) {
                    window.themeManager.setTheme(theme);
                }
            });
        });
        
        // Handle browser back/forward
        window.addEventListener('hashchange', () => {
            const page = this.getPageFromHash();
            if (page && page !== this.currentPage) {
                this.navigateTo(page, false);
            }
        });
    }
    
    getPageFromHash() {
        const hash = window.location.hash.slice(1);
        return hash || null;
    }
    
    async navigateTo(pageName, updateHash = true) {
        const pageConfig = this.pages[pageName];
        
        if (!pageConfig) {
            console.warn(`Page not found: ${pageName}`);
            pageName = 'dashboard';
        }
        
        console.log(`📄 Navigating to: ${pageName}`);
        
        // Update current page
        this.currentPage = pageName;
        
        // Update URL hash
        if (updateHash) {
            window.location.hash = pageName;
        }
        
        // Update page title
        this.updatePageTitle(pageName);
        
        // Update active nav item
        this.updateActiveNav(pageName);
        
        // Show loading
        this.showLoading();
        
        // Load page content
        await this.loadPageContent(pageName);
    }
    
    updatePageTitle(pageName) {
        const config = this.pages[pageName] || this.pages['dashboard'];
        const titleEl = document.getElementById('pageTitle');
        titleEl.innerHTML = `<i class="fas ${config.icon}"></i><span>${config.title}</span>`;
        document.title = `Bar Restaurant - ${config.title}`;
    }
    
    updateActiveNav(pageName) {
        // Remove all active states
        document.querySelectorAll('.nav-link').forEach(link => {
            link.classList.remove('active');
        });
        document.querySelectorAll('.nav-toggle').forEach(toggle => {
            toggle.classList.remove('active');
        });
        
        // Add active state to current page
        const activeLink = document.querySelector(`.nav-link[data-page="${pageName}"]`);
        if (activeLink) {
            activeLink.classList.add('active');
            
            // Expand parent submenu if needed
            const submenu = activeLink.closest('.nav-submenu');
            if (submenu) {
                submenu.classList.add('expanded');
                const toggle = document.querySelector(`[data-target="${submenu.id}"]`);
                if (toggle) {
                    toggle.classList.add('expanded');
                    toggle.classList.add('active');
                }
            }
        }
    }
    
    showLoading() {
        this.pageContent.innerHTML = `
            <div class="content-loading">
                <div class="spinner-border text-primary" role="status">
                    <span class="visually-hidden">Loading...</span>
                </div>
            </div>
        `;
    }
    
    async loadPageContent(pageName) {
        const config = this.pages[pageName] || this.pages['dashboard'];
        
        try {
            const response = await fetch(config.partial);
            if (!response.ok) {
                throw new Error(`Failed to load ${config.partial}`);
            }
            
            const html = await response.text();
            this.pageContent.innerHTML = html;
            
            // Execute any inline scripts in the loaded content
            this.executePageScripts();
            
            // Initialize page-specific functionality
            this.initializePage(pageName);
            
        } catch (error) {
            console.error('Error loading page:', error);
            this.pageContent.innerHTML = `
                <div class="text-center py-5">
                    <i class="fas fa-exclamation-triangle fa-3x text-warning mb-3"></i>
                    <h4>Page Not Found</h4>
                    <p class="text-muted">The requested page could not be loaded.</p>
                    <button class="btn btn-primary" onclick="app.navigateTo('dashboard')">
                        <i class="fas fa-home me-2"></i>Go to Dashboard
                    </button>
                </div>
            `;
        }
    }
    
    executePageScripts() {
        // Find and execute scripts in loaded content
        const scripts = this.pageContent.querySelectorAll('script');
        scripts.forEach(oldScript => {
            const newScript = document.createElement('script');
            if (oldScript.src) {
                newScript.src = oldScript.src;
            } else {
                newScript.textContent = oldScript.textContent;
            }
            oldScript.parentNode.replaceChild(newScript, oldScript);
        });
    }
    
    initializePage(pageName) {
        // Page-specific initialization can be added here
        console.log(`✅ Page loaded: ${pageName}`);
    }
    
    saveMenuState() {
        const expandedMenus = [];
        document.querySelectorAll('.nav-toggle.expanded').forEach(toggle => {
            expandedMenus.push(toggle.dataset.target);
        });
        localStorage.setItem('barrest_menu_state', JSON.stringify(expandedMenus));
    }
    
    restoreMenuState() {
        const saved = localStorage.getItem('barrest_menu_state');
        if (saved) {
            const expandedMenus = JSON.parse(saved);
            expandedMenus.forEach(targetId => {
                const toggle = document.querySelector(`[data-target="${targetId}"]`);
                const submenu = document.getElementById(targetId);
                if (toggle && submenu) {
                    toggle.classList.add('expanded');
                    submenu.classList.add('expanded');
                }
            });
        }
    }
    
    async loadSystemStatus() {
        try {
            await loadSystemStatus('headerSystemStatus', {
                checkInterval: 30000,
                autoStart: true,
                variant: 'dot'
            });
        } catch (error) {
            console.error('Failed to load system status:', error);
        }
    }
    
    async handleLogout() {
        const result = await Swal.fire({
            title: 'Logout',
            text: 'Are you sure you want to logout?',
            icon: 'question',
            showCancelButton: true,
            confirmButtonColor: '#dc3545',
            cancelButtonColor: '#6c757d',
            confirmButtonText: 'Yes, logout',
            cancelButtonText: 'Cancel'
        });
        
        if (result.isConfirmed) {
            try {
                await window.authService.logout();
                window.location.href = 'login.html';
            } catch (error) {
                console.error('Logout error:', error);
                window.location.href = 'login.html';
            }
        }
    }
}

// Export for global access
window.AppRouter = AppRouter;
