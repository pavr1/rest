// === BAR RESTAURANT - SYSTEM STATUS COMPONENT ===
// Reusable component for displaying system health status across pages

class SystemStatusMonitor {
    constructor(options = {}) {
        // Configuration
        this.checkInterval = options.checkInterval || 5000; // 5 seconds default
        this.autoStart = options.autoStart !== false; // Auto-start by default
        this.containerId = options.containerId || 'systemStatusContainer';
        this.variant = options.variant || 'full'; // 'full', 'compact', or 'dot'
        
        // Elements
        this.container = document.getElementById(this.containerId);
        this.refreshBtn = document.getElementById('refreshStatusBtn');
        this.mainDot = document.getElementById('main-status');
        this.dotTrigger = document.getElementById('statusDotTrigger');
        this.dropdownMenu = document.getElementById('statusDropdownMenu');
        this.inventoryStatus = document.getElementById('inventory-status');
        // State
        this.healthCheckInterval = null;
        this.isRunning = false;
        this.isDropdownOpen = false;
        
        // Ensure CONFIG is available
        if (typeof CONFIG === 'undefined') {
            console.error('❌ CONFIG not available for SystemStatusMonitor');
            return;
        }
        
        console.log('🏥 SystemStatusMonitor initialized');
        
        // Bind methods
        this.checkGatewayHealth = this.checkGatewayHealth.bind(this);
        this.refreshStatus = this.refreshStatus.bind(this);
        this.toggleDropdown = this.toggleDropdown.bind(this);
        this.closeDropdown = this.closeDropdown.bind(this);
        
        // Initialize
        this.init();
    }

    init() {
        if (!this.container) {
            console.warn('⚠️ System status container not found');
            return;
        }
        
        // Bind refresh button
        if (this.refreshBtn) {
            this.refreshBtn.addEventListener('click', (e) => {
                e.stopPropagation();
                this.refreshStatus();
            });
        }
        
        // Bind dot trigger for dropdown variant
        if (this.dotTrigger && this.dropdownMenu) {
            this.dotTrigger.addEventListener('click', this.toggleDropdown);
            
            // Initialize Bootstrap tooltip if available
            if (typeof bootstrap !== 'undefined' && bootstrap.Tooltip) {
                this.tooltip = new bootstrap.Tooltip(this.dotTrigger);
            }
            
            // Close dropdown when clicking outside
            document.addEventListener('click', (e) => {
                if (!this.container.contains(e.target)) {
                    this.closeDropdown();
                }
            });
        }
        
        // Auto-start monitoring if enabled
        if (this.autoStart) {
            this.start();
        }
        
        console.log('✅ SystemStatusMonitor ready');
    }
    
    toggleDropdown(e) {
        e.stopPropagation();
        this.isDropdownOpen = !this.isDropdownOpen;
        
        if (this.isDropdownOpen) {
            this.dropdownMenu.classList.add('show');
            // Hide tooltip when dropdown is open
            if (this.tooltip) {
                this.tooltip.hide();
            }
        } else {
            this.dropdownMenu.classList.remove('show');
        }
    }
    
    closeDropdown() {
        this.isDropdownOpen = false;
        if (this.dropdownMenu) {
            this.dropdownMenu.classList.remove('show');
        }
    }

    // === CONTROL METHODS ===

    start() {
        if (this.isRunning) {
            console.warn('⚠️ SystemStatusMonitor already running');
            return;
        }
        
        console.log('🏥 Starting system health monitoring...');
        this.isRunning = true;
        
        // Initial check
        this.checkGatewayHealth();
        
        // Schedule periodic checks
        this.healthCheckInterval = setInterval(
            this.checkGatewayHealth,
            this.checkInterval
        );
    }

    stop() {
        if (!this.isRunning) {
            return;
        }
        
        console.log('🛑 Stopping system health monitoring...');
        this.isRunning = false;
        
        if (this.healthCheckInterval) {
            clearInterval(this.healthCheckInterval);
            this.healthCheckInterval = null;
        }
    }

    async refreshStatus() {
        if (!this.refreshBtn) return;
        
        const icon = this.refreshBtn.querySelector('i');
        if (icon) {
            icon.classList.add('fa-spin');
        }
        this.refreshBtn.disabled = true;
        
        await this.checkGatewayHealth();
        
        if (icon) {
            icon.classList.remove('fa-spin');
        }
        this.refreshBtn.disabled = false;
    }

    // === HEALTH CHECK METHODS ===

    async checkGatewayHealth() {
        // Trigger heartbeat animation on all indicators
        this.triggerHeartbeat();
        
        try {
            const controller = new AbortController();
            const timeoutId = setTimeout(() => controller.abort(), 5000);
            
            const response = await fetch(CONFIG.SERVICES.gateway, {
                method: 'GET',
                signal: controller.signal
            });
            
            clearTimeout(timeoutId);
            
            // Parse response even if status is 503 (degraded)
            // Gateway returns 503 when some services are unhealthy but still returns service details
            if (response.ok || response.status === 503) {
                const data = await response.json();
                
                // Gateway returns: { is_healthy: bool, services: { "service-name": bool } }
                this.updateStatusFromGateway(data.services || {}, data.is_healthy);
            } else {
                console.warn('🏥 Gateway returned unexpected status:', response.status);
                this.markAllServicesOffline();
            }
        } catch (error) {
            console.error('🔄 Health check failed:', error.message);
            this.markAllServicesOffline();
        }
    }

    updateStatusFromGateway(services, isHealthy) {
        // Update main dot (for dot variant) - overall system health
        const mainStatus = isHealthy ? 'online' : 'offline';
        this.updateStatusIndicator('main-status', mainStatus);
        
        // Update tooltip with current status
        this.updateTooltip(isHealthy);
        
        // Gateway status is determined by overall system health
        const gatewayStatus = isHealthy ? 'online' : 'warning';
        this.updateStatusIndicator('gateway-status', gatewayStatus);
        
        // Update session service status
        const sessionStatus = services['session-service'] === true ? 'online' : 'offline';
        this.updateStatusIndicator('session-status', sessionStatus);
        
        // Update menu service status
        const menuStatus = services['menu-service'] === true ? 'online' : 'offline';
        this.updateStatusIndicator('menu-status', menuStatus);
        
        // Update data service status from actual gateway response
        const dataStatus = services['data-service'] === true ? 'online' : 'offline';
        this.updateStatusIndicator('data-status', dataStatus);

        // Update inventory service status
        const inventoryStatus = services['inventory-service'] === true ? 'online' : 'offline';
        this.updateStatusIndicator('inventory-status', inventoryStatus);

        // Update invoice service status
        const invoiceStatus = services['invoice-service'] === true ? 'online' : 'offline';
        this.updateStatusIndicator('invoice-status', invoiceStatus);
    }
    
    updateTooltip(isHealthy) {
        if (this.tooltip && this.dotTrigger) {
            const statusText = isHealthy ? 'All Systems Operational' : 'System Issues Detected';
            this.dotTrigger.setAttribute('data-bs-title', statusText);
            this.tooltip.dispose();
            this.tooltip = new bootstrap.Tooltip(this.dotTrigger);
        }
    }

    updateStatusIndicator(elementId, status) {
        const element = document.getElementById(elementId);
        if (!element) return;
        
        element.classList.remove('online', 'offline', 'warning', 'loading');
        element.classList.add(status);
    }

    markAllServicesOffline() {
        this.updateStatusIndicator('main-status', 'offline');
        this.updateStatusIndicator('gateway-status', 'offline');
        this.updateStatusIndicator('session-status', 'offline');
        this.updateStatusIndicator('menu-status', 'offline');
        this.updateStatusIndicator('data-status', 'offline');
        this.updateStatusIndicator('inventory-status', 'offline');
    }

    triggerHeartbeat() {
        // Trigger heartbeat animation on all status indicators (including main dot)
        const indicators = [
            'main-status',
            'gateway-status',
            'session-status',
            'menu-status',
            'data-status',
            'inventory-status'
        ];
        
        indicators.forEach(id => {
            const element = document.getElementById(id);
            if (!element) return;
            
            // Remove heartbeat class first to restart animation if it's already running
            element.classList.remove('heartbeat');
            
            // Force reflow to restart animation
            void element.offsetWidth;
            
            // Add heartbeat class
            element.classList.add('heartbeat');
            
            // Remove after animation completes (800ms)
            setTimeout(() => {
                element.classList.remove('heartbeat');
            }, 800);
        });
    }

    // === UTILITY METHODS ===

    getStatus() {
        return {
            gateway: this.getIndicatorStatus('gateway-status'),
            session: this.getIndicatorStatus('session-status'),
            menu: this.getIndicatorStatus('menu-status'),
            data: this.getIndicatorStatus('data-status'),
            inventory: this.getIndicatorStatus('inventory-status')
        };
    }

    getIndicatorStatus(elementId) {
        const element = document.getElementById(elementId);
        if (!element) return null;
        
        if (element.classList.contains('online')) return 'online';
        if (element.classList.contains('offline')) return 'offline';
        if (element.classList.contains('warning')) return 'warning';
        if (element.classList.contains('loading')) return 'loading';
        
        return 'unknown';
    }

    destroy() {
        this.stop();
        
        if (this.refreshBtn) {
            this.refreshBtn.removeEventListener('click', this.refreshStatus);
        }
        
        console.log('🧹 SystemStatusMonitor destroyed');
    }
}

// === HELPER FUNCTION TO LOAD SYSTEM STATUS HTML ===

async function loadSystemStatus(containerId, options = {}) {
    try {
        const container = document.getElementById(containerId);
        if (!container) {
            console.error('❌ Container not found:', containerId);
            return null;
        }
        
        // Determine which partial to load
        let partialPath = 'shared/partials/system-status.html';
        let variant = 'full';
        
        if (options.variant === 'dot') {
            partialPath = 'shared/partials/system-status-dot.html';
            variant = 'dot';
        } else if (options.compact === true) {
            partialPath = 'shared/partials/system-status-compact.html';
            variant = 'compact';
        }
        
        // Load the HTML partial
        const response = await fetch(partialPath);
        if (!response.ok) {
            throw new Error(`Failed to load ${partialPath}: ${response.status}`);
        }
        
        const html = await response.text();
        container.innerHTML = html;
        
        // Create and return the monitor
        const monitor = new SystemStatusMonitor({
            containerId: 'systemStatusContainer',
            variant: variant,
            ...options
        });
        
        console.log(`✅ System status component loaded (${variant})`);
        return monitor;
        
    } catch (error) {
        console.error('❌ Failed to load system status component:', error);
        return null;
    }
}

// Export for global access
window.SystemStatusMonitor = SystemStatusMonitor;
window.loadSystemStatus = loadSystemStatus;

console.log('✅ system-status.js loaded');
