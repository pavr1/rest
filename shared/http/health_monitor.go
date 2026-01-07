package http

import (
	"context"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
)

// ServiceHealth tracks the health state of an HTTP service
type ServiceHealth struct {
	Name      string
	URL       string
	Healthy   bool
	LastCheck time.Time
}

// HealthMonitor can monitor either database health or HTTP services
type HTTPHealthMonitor struct {
	logger   *logrus.Logger
	interval time.Duration

	// For HTTP monitoring
	client   *http.Client
	mu       sync.RWMutex
	services map[string]*ServiceHealth
	healthy  atomic.Bool
}

// NewHealthMonitor creates a new health monitor for database
func NewHealthMonitor(logger *logrus.Logger, interval time.Duration) (*HTTPHealthMonitor, error) {
	hm := &HTTPHealthMonitor{
		logger:   logger,
		interval: interval,
		client:   &http.Client{},
		services: make(map[string]*ServiceHealth),
	}
	hm.logger.WithFields(logrus.Fields{
		"interval": interval,
	}).Info("Creating new health monitor")

	// Start as unhealthy until first successful check
	hm.healthy.Store(true)
	return hm, nil
}

func (hm *HTTPHealthMonitor) AddService(name string, url string) {
	hm.mu.Lock()
	defer hm.mu.Unlock()
	hm.services[name] = &ServiceHealth{
		Name: name,
		URL:  url,
	}
}

// Concurrent start begins the background health monitoring
func (hm *HTTPHealthMonitor) Start(ctx context.Context) {
	go hm.startHTTPMonitoring(ctx)
}

// startHTTPMonitoring monitors HTTP services
func (hm *HTTPHealthMonitor) startHTTPMonitoring(ctx context.Context) {
	hm.logger.WithField("interval", hm.interval).Info("🏥 HTTP health monitor starting")

	// Initial check
	hm.checkAllServices()

	ticker := time.NewTicker(hm.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			hm.logger.Info("HTTP health monitor stopped")
			return
		case <-ticker.C:
			hm.checkAllServices()
		}
	}
}

// checkAllServices checks all HTTP services
func (hm *HTTPHealthMonitor) checkAllServices() {
	hm.mu.RLock()
	services := make([]*ServiceHealth, 0, len(hm.services))
	for _, svc := range hm.services {
		services = append(services, svc)
	}
	hm.mu.RUnlock()

	for _, svc := range services {
		hm.checkService(svc)
	}
}

// checkService pings a single HTTP service
func (hm *HTTPHealthMonitor) checkService(svc *ServiceHealth) {
	req, err := http.NewRequest("GET", svc.URL, nil)
	if err != nil {
		hm.logger.WithFields(logrus.Fields{
			"service": svc.Name,
			"error":   err.Error(),
		}).Error("Failed to create HTTP request")

		hm.setServiceHealth(svc.Name, false)
		return
	}
	req.Header.Set("X-Health-Check", "true")

	resp, err := hm.client.Do(req)
	if err != nil {
		hm.setServiceHealth(svc.Name, false)
		hm.logger.WithFields(logrus.Fields{
			"service": svc.Name,
		}).Error("Health check failed")
		return
	}
	defer resp.Body.Close()

	healthy := resp.StatusCode == http.StatusOK
	hm.setServiceHealth(svc.Name, healthy)
}

// setServiceHealth updates HTTP service health state
func (hm *HTTPHealthMonitor) setServiceHealth(name string, healthy bool) {
	hm.mu.Lock()
	defer hm.mu.Unlock()
	if svc, ok := hm.services[name]; ok {
		svc.Healthy = healthy
		svc.LastCheck = time.Now()
	}
}

// HealthStatus represents the overall health with individual service statuses
type HealthStatus struct {
	IsHealthy bool            `json:"is_healthy"`
	Services  map[string]bool `json:"services"`
}

// GetHealthStatus returns overall health and individual service statuses
func (hm *HTTPHealthMonitor) GetHealthStatus() HealthStatus {
	hm.mu.RLock()
	defer hm.mu.RUnlock()

	services := make(map[string]bool)
	allHealthy := true

	for name, svc := range hm.services {
		services[name] = svc.Healthy
		if !svc.Healthy {
			allHealthy = false
		}
	}

	// If no services registered, consider unhealthy
	if len(services) == 0 {
		allHealthy = false
	}

	return HealthStatus{
		IsHealthy: allHealthy,
		Services:  services,
	}
}

// IsHealthy returns the current overall health state
func (hm *HTTPHealthMonitor) IsHealthy() bool {
	return hm.GetHealthStatus().IsHealthy
}
