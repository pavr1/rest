package db

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
)

// DBPinger interface for database health checks
type DBPinger interface {
	Ping() error
}

// ServiceHealth tracks the health state of an HTTP service
type ServiceHealth struct {
	Name      string
	URL       string
	Healthy   bool
	LastCheck time.Time
}

// HealthMonitor can monitor either database health or HTTP services
type HealthMonitor struct {
	logger   *logrus.Logger
	interval time.Duration

	// For DB monitoring
	dbHandler *DbHandler
	dbHealthy atomic.Bool

	// For HTTP monitoring
	client   *http.Client
	mu       sync.RWMutex
	services map[string]*ServiceHealth
}

// NewHealthMonitor creates a new health monitor for database
func NewHealthMonitor(logger *logrus.Logger, interval time.Duration, dbHandler *DbHandler) (*HealthMonitor, error) {
	if dbHandler == nil {
		logger.Error("Database health monitor is not initialized")
		return nil, errors.New("database health monitor is not initialized")
	}

	hm := &HealthMonitor{
		logger:    logger,
		interval:  interval,
		dbHandler: dbHandler,
	}
	hm.logger.WithFields(logrus.Fields{
		"interval":  interval,
		"dbHandler": dbHandler,
	}).Info("Creating new health monitor")

	// Start as unhealthy until first successful check
	hm.dbHealthy.Store(false)
	return hm, nil
}

// AddService adds an HTTP service to monitor
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
	// Determine monitoring mode
	if hm.dbHandler != nil {
		hm.startDBMonitoring(ctx)
	} else if hm.services != nil {
		hm.startHTTPMonitoring(ctx)
	}
}

// startDBMonitoring monitors database health
func (hm *HealthMonitor) startDBMonitoring(ctx context.Context) {
	hm.logger.WithField("interval", hm.interval).Info("🏥 Database health monitor starting")

	// Initial check
	hm.checkDatabase()

	ticker := time.NewTicker(hm.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			hm.logger.Info("Database health monitor stopped")
			return
		case <-ticker.C:
			hm.checkDatabase()
		}
	}
}

// startHTTPMonitoring monitors HTTP services
func (hm *HealthMonitor) startHTTPMonitoring(ctx context.Context) {
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

// checkDatabase pings the database
func (hm *HealthMonitor) checkDatabase() {
	err := hm.dbHandler.Ping()

	if err != nil {
		// Only log if transitioning from healthy to unhealthy
		//if hm.dbHealthy.Load() {
		hm.logger.WithError(err).Error("Database health check failed")
		//}
		hm.dbHealthy.Store(false)
	} else {
		// Only log if transitioning from unhealthy to healthy
		if !hm.dbHealthy.Load() {
			hm.logger.Info("🔥 Database health connected successfully")
		}
		hm.dbHealthy.Store(true)
	}
}

// checkAllServices checks all HTTP services
func (hm *HealthMonitor) checkAllServices() {
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
func (hm *HealthMonitor) checkService(svc *ServiceHealth) {
	req, err := http.NewRequest("GET", svc.URL, nil)
	if err != nil {
		hm.setServiceHealth(svc.Name, false)
		return
	}
	req.Header.Set("X-Health-Check", "true")

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

// setServiceHealth updates HTTP service health state
func (hm *HealthMonitor) setServiceHealth(name string, healthy bool) {
	hm.mu.Lock()
	defer hm.mu.Unlock()
	if svc, ok := hm.services[name]; ok {
		svc.Healthy = healthy
		svc.LastCheck = time.Now()
	}
}

// IsHealthy returns the current database health state (for DB monitoring)
func (hm *HealthMonitor) IsHealthy() bool {
	return hm.dbHealthy.Load()
}

// IsServiceHealthy returns the health state of a specific HTTP service
func (hm *HealthMonitor) IsServiceHealthy(name string) bool {
	hm.mu.RLock()
	defer hm.mu.RUnlock()
	if svc, ok := hm.services[name]; ok {
		return svc.Healthy
	}
	return false
}

// AreAllServicesHealthy returns true if all monitored HTTP services are healthy
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

// GetAllServicesStatus returns status of all HTTP services
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

// SetHealthyForTesting allows tests to set health state directly
// This should only be used in tests
func (hm *HealthMonitor) SetHealthyForTesting(healthy bool) {
	hm.dbHealthy.Store(healthy)
}
