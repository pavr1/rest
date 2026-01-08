package db

import (
	"context"
	"errors"
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

// DBHealthMonitor can monitor database health
type DBHealthMonitor struct {
	logger   *logrus.Logger
	interval time.Duration

	// For DB monitoring
	dbHandler *DbHandler
	dbHealthy atomic.Bool
}

// NewHealthMonitor creates a new health monitor for database
func NewHealthMonitor(logger *logrus.Logger, interval time.Duration, dbHandler *DbHandler) (*DBHealthMonitor, error) {
	if dbHandler == nil {
		logger.Error("Database handler is not initialized")
		return nil, errors.New("database handler is not initialized")
	}

	hm := &DBHealthMonitor{
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

// Concurrent start begins the background health monitoring
func (hm *DBHealthMonitor) Start(ctx context.Context) {
	go hm.startMonitoring(ctx)
}

// startDBMonitoring monitors database health
func (hm *DBHealthMonitor) startMonitoring(ctx context.Context) {
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

// checkDatabase pings the database
func (hm *DBHealthMonitor) checkDatabase() {
	err := hm.dbHandler.Ping()

	if err != nil {
		// Only log if transitioning from healthy to unhealthy
		if hm.dbHealthy.Load() {
			hm.logger.WithError(err).Error("❌ Database health check failed")
		}
		hm.dbHealthy.Store(false)
	} else {
		// Only log if transitioning from unhealthy to healthy
		if !hm.dbHealthy.Load() {
			hm.logger.Info("✅ Database health connected successfully")
		}
		hm.dbHealthy.Store(true)
	}
}

// IsHealthy returns the current database health state (for DB monitoring)
func (hm *DBHealthMonitor) IsHealthy() bool {
	return hm.dbHealthy.Load()
}
