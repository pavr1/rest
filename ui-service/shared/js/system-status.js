// === BAR RESTAURANT - SYSTEM STATUS COMPONENT ===
// Reusable component for displaying system health status across pages

class SystemStatusMonitor {
    constructor(options = {}) {
        // Configuration
        this.checkInterval = options.checkInterval || 5000; // 5 seconds default
        this.autoStart = options.autoStart !== false; // Auto-start by default
        this.containerId = options.containerId || 'systemStatusContainer';
        
        // Elements
        this.container = document.getElementById(this.containerId);
        this.refreshBtn = document.getElementById('refreshStatusBtn');
        
        // State
        this.healthCheckInterval = null;
        this.isRunning = false;
        
        // Ensure CONFIG is available
        if (typeof CONFIG === 'undefined') {
            console.error('❌ CONFIG not available for SystemStatusMonitor');
            return;
        }
        
        console.log('🏥 SystemStatusMonitor initialized');
        
        // Bind methods
        this.checkGatewayHealth = this.checkGatewayHealth.bind(this);
        this.refreshStatus = this.refreshStatus.bind(this);
        
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
            this.refreshBtn.addEventListener('click', this.refreshStatus);
        }
        
        // Auto-start monitoring if enabled
        if (this.autoStart) {
            this.start();
        }
        
        console.log('✅ SystemStatusMonitor ready');
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
        // Gateway status is determined by overall system health
        const gatewayStatus = isHealthy ? 'online' : 'warning';
        this.updateStatusIndicator('gateway-status', gatewayStatus);
        
        // Update session service status
        const sessionStatus = services['session-service'] === true ? 'online' : 'offline';
        this.updateStatusIndicator('session-status', sessionStatus);
        
        // Update data service status from actual gateway response
        const dataStatus = services['data-service'] === true ? 'online' : 'offline';
        this.updateStatusIndicator('data-status', dataStatus);
    }

    updateStatusIndicator(elementId, status) {
        const element = document.getElementById(elementId);
        if (!element) return;
        
        element.classList.remove('online', 'offline', 'warning', 'loading');
        element.classList.add(status);
    }

    markAllServicesOffline() {
        this.updateStatusIndicator('gateway-status', 'offline');
        this.updateStatusIndicator('session-status', 'offline');
        this.updateStatusIndicator('data-status', 'offline');
    }

    triggerHeartbeat() {
        // Trigger heartbeat animation on all status indicators
        const indicators = [
            'gateway-status',
            'session-status',
            'data-status'
        ];
        
        console.log('💓 Triggering heartbeat animation');
        
        indicators.forEach(id => {
            const element = document.getElementById(id);
            if (!element) {
                console.warn(`⚠️ Element ${id} not found for heartbeat`);
                return;
            }
            
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
            data: this.getIndicatorStatus('data-status')
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
        
        // Load the HTML partial
        const response = await fetch('shared/partials/system-status.html');
        if (!response.ok) {
            throw new Error(`Failed to load system-status.html: ${response.status}`);
        }
        
        const html = await response.text();
        container.innerHTML = html;
        
        // Create and return the monitor
        const monitor = new SystemStatusMonitor({
            containerId: 'systemStatusContainer',
            ...options
        });
        
        console.log('✅ System status component loaded');
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
