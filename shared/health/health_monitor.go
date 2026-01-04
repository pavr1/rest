package health

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// ServiceHealth tracks the health state of a single dependency
type ServiceHealth struct {
	Name      string
	URL       string
	Healthy   bool
	LastCheck time.Time
}

// HealthMonitor continuously monitors dependency health
type HealthMonitor struct {
	logger   *logrus.Logger
	interval time.Duration
	client   *http.Client
	mu       sync.RWMutex
	services map[string]*ServiceHealth
}

// NewHealthMonitor creates a new health monitor
func NewHealthMonitor(logger *logrus.Logger, interval time.Duration) *HealthMonitor {
	return &HealthMonitor{
		logger:   logger,
		interval: interval,
		client: &http.Client{
			//pvillalobos this should be configurable
			Timeout: 1 * time.Second,
		},
		services: make(map[string]*ServiceHealth),
	}
}

// AddService adds a service to monitor
func (hm *HealthMonitor) AddService(name, healthURL string) {
	hm.mu.Lock()
	defer hm.mu.Unlock()
	hm.services[name] = &ServiceHealth{
		Name:    name,
		URL:     healthURL,
		Healthy: false, // Start as unhealthy until first successful check
	}
	hm.logger.WithFields(logrus.Fields{
		"service": name,
		"url":     healthURL,
	}).Info("Added service to health monitor")
}

// Start begins the background health monitoring
func (hm *HealthMonitor) Start(ctx context.Context) {
	hm.logger.WithField("interval", hm.interval).Info("🏥 Health monitor starting")

	// Initial check for all services
	hm.checkAllServices()

	ticker := time.NewTicker(hm.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			hm.logger.Info("Health monitor stopped")
			return
		case <-ticker.C:
			hm.checkAllServices()
		}
	}
}

// checkAllServices pings all registered services
func (hm *HealthMonitor) checkAllServices() {
	for _, svc := range hm.services {
		hm.checkService(svc)
	}
}

// checkService pings a single service health endpoint
func (hm *HealthMonitor) checkService(svc *ServiceHealth) {
	req, err := http.NewRequest("GET", svc.URL, nil)
	if err != nil {
		hm.setServiceHealth(svc.Name, false)
		return
	}
	req.Header.Set("X-Health-Check", "true")
	req.Header.Set("X-User-ID", "system")
	req.Header.Set("X-User-Role", "admin")

	resp, err := hm.client.Do(req)
	if err != nil {
		hm.setServiceHealth(svc.Name, false)
		hm.logger.WithFields(logrus.Fields{
			"service": svc.Name,
		}).Warn("Health check failed")
		return
	}
	defer resp.Body.Close()

	healthy := resp.StatusCode == http.StatusOK
	hm.setServiceHealth(svc.Name, healthy)
}

// setServiceHealth updates the health state thread-safely
func (hm *HealthMonitor) setServiceHealth(name string, healthy bool) {
	hm.mu.Lock()
	defer hm.mu.Unlock()
	if svc, ok := hm.services[name]; ok {
		svc.Healthy = healthy
		svc.LastCheck = time.Now()
	}
}

// IsServiceHealthy returns the health state of a specific service
func (hm *HealthMonitor) IsServiceHealthy(name string) bool {
	hm.mu.RLock()
	defer hm.mu.RUnlock()
	if svc, ok := hm.services[name]; ok {
		return svc.Healthy
	}
	return false
}

// AreAllServicesHealthy returns true if all monitored services are healthy
func (hm *HealthMonitor) AreAllServicesHealthy() bool {
	hm.mu.RLock()
	defer hm.mu.RUnlock()
	for _, svc := range hm.services {
		if !svc.Healthy {
			return false
		}
	}
	return len(hm.services) > 0
}

// GetServiceStatus returns the status of a specific service
func (hm *HealthMonitor) GetServiceStatus(name string) *ServiceHealth {
	hm.mu.RLock()
	defer hm.mu.RUnlock()
	if svc, ok := hm.services[name]; ok {
		// Return a copy to avoid race conditions
		return &ServiceHealth{
			Name:      svc.Name,
			URL:       svc.URL,
			Healthy:   svc.Healthy,
			LastCheck: svc.LastCheck,
		}
	}
	return nil
}

// GetAllServicesStatus returns status of all services
func (hm *HealthMonitor) GetAllServicesStatus() map[string]*ServiceHealth {
	hm.mu.RLock()
	defer hm.mu.RUnlock()
	result := make(map[string]*ServiceHealth)
	for name, svc := range hm.services {
		result[name] = &ServiceHealth{
			Name:      svc.Name,
			URL:       svc.URL,
			Healthy:   svc.Healthy,
			LastCheck: svc.LastCheck,
		}
	}
	return result
}

// SetServiceHealthForTesting allows tests to set health state directly
// This should only be used in tests
func (hm *HealthMonitor) SetServiceHealthForTesting(name string, healthy bool) {
	hm.mu.Lock()
	defer hm.mu.Unlock()
	if svc, ok := hm.services[name]; ok {
		svc.Healthy = healthy
		svc.LastCheck = time.Now()
	}
}
