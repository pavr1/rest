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
            'menu-categories': { title: 'Categories', icon: 'fa-tags', partial: 'pages/menu/categories.html' },
            'menu-sub-categories': { title: 'Sub-Categories', icon: 'fa-layer-group', partial: 'pages/menu/sub-categories.html' },
            'menu-variants': { title: 'Variants', icon: 'fa-hamburger', partial: 'pages/menu/variants.html' },
            'menu-ingredients': { title: 'Ingredients', icon: 'fa-list-ul', partial: 'pages/menu/ingredients.html' },
            
            // Inventory
            'inventory-categories': { title: 'Inventory Categories', icon: 'fa-layer-group', partial: 'pages/inventory/categories.html' },
            'inventory-sub-categories': { title: 'Inventory Sub-Categories', icon: 'fa-list', partial: 'pages/inventory/sub-categories.html' },
            'inventory-variants': { title: 'Inventory Variants', icon: 'fa-boxes', partial: 'pages/inventory/variants.html' },
            'stock-count': { title: 'Stock Count', icon: 'fa-warehouse', partial: 'pages/inventory/stock-count.html' },
            'menu-ingredients': { title: 'Ingredients', icon: 'fa-list-ul', partial: 'pages/menu/ingredients.html' },
            'suppliers': { title: 'Suppliers', icon: 'fa-truck', partial: 'pages/inventory/suppliers.html' },

            // Invoices
            'outcome-invoices': { title: 'Purchases', icon: 'fa-shopping-cart', partial: 'pages/invoices/outcome.html' },
            'income-invoices': { title: 'Sells', icon: 'fa-cash-register', partial: 'pages/invoices/income.html' },

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
                const isExpanding = !toggle.classList.contains('expanded');
                
                // Collapse all other menus first, then expand selected after a short delay
                if (isExpanding) {
                    document.querySelectorAll('.nav-toggle.expanded').forEach(otherToggle => {
                        if (otherToggle !== toggle) {
                            otherToggle.classList.remove('expanded');
                            const otherSubmenu = document.getElementById(otherToggle.dataset.target);
                            if (otherSubmenu) otherSubmenu.classList.remove('expanded');
                        }
                    });
                    
                    // Small delay to let close animation finish before opening
                    setTimeout(() => {
                        toggle.classList.add('expanded');
                        submenu.classList.add('expanded');
                        this.saveMenuState();
                    }, 150);
                } else {
                    // Just collapse if already expanded
                    toggle.classList.remove('expanded');
                    submenu.classList.remove('expanded');
                    this.saveMenuState();
                }
            });
        });
        
        // Navigation links
        document.querySelectorAll('.nav-link[data-page]').forEach(link => {
            link.addEventListener('click', (e) => {
                // If this element is also a nav-toggle, don't navigate on click
                // The toggle functionality will handle the expand/collapse
                if (link.classList.contains('nav-toggle')) {
                    return;
                }

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
        const config = this.pages[pageName];

        if (!config) {
            console.error(`Page config not found for: ${pageName}, falling back to dashboard`);
            return this.loadPageContent('dashboard');
        }

        console.log(`🔄 Loading page content for: ${pageName} from ${config.partial}`);

        try {
            // Add cache-busting parameter to prevent browser caching
            const cacheBust = `?v=${Date.now()}`;
            const response = await fetch(config.partial + cacheBust);
            console.log(`📡 Fetch response for ${config.partial}:`, response.status, response.statusText);

            if (!response.ok) {
                throw new Error(`HTTP ${response.status}: Failed to load ${config.partial}`);
            }

            const html = await response.text();
            console.log(`✅ Loaded ${html.length} characters for ${pageName}`);

            this.pageContent.innerHTML = html;

            // Execute any inline scripts in the loaded content
            this.executePageScripts();

        // Initialize page-specific functionality
        this.initializePage(pageName);

        // Debug logging for inventory pages
        if (pageName.includes('inventory')) {
            console.log(`🎯 Inventory page loaded: ${pageName}`);
            console.log('Page content length:', this.pageContent.innerHTML.length);
            console.log('Page content preview:', this.pageContent.innerHTML.substring(0, 200) + '...');
        }

        } catch (error) {
            console.error('❌ Error loading page:', error);
            this.pageContent.innerHTML = `
                <div class="text-center py-5">
                    <i class="fas fa-exclamation-triangle fa-3x text-warning mb-3"></i>
                    <h4>Page Load Error</h4>
                    <p class="text-muted">Failed to load page: ${pageName}</p>
                    <p class="text-muted small">Error: ${error.message}</p>
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
