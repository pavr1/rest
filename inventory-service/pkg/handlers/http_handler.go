package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	sharedConfig "shared/config"
	sharedDb "shared/db"
	sharedHttp "shared/http"

	stockCategoryHandlers "inventory-service/pkg/entities/stock_item_categories/handlers"
	stockItemHandlers "inventory-service/pkg/entities/stock_items/handlers"
	supplierHandlers "inventory-service/pkg/entities/suppliers/handlers"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type MainHTTPHandler struct {
	db                   *sharedDb.DbHandler
	httpHealthMonitor    *sharedHttp.HTTPHealthMonitor
	cancelHealthMonitor  context.CancelFunc
	stockCategoryHandler *stockCategoryHandlers.HTTPHandler
	stockItemHandler     *stockItemHandlers.HTTPHandler
	supplierHandler      *supplierHandlers.HTTPHandler
	logger               *logrus.Logger
}

func NewHTTPHandler(cfg *sharedConfig.Config, logger *logrus.Logger) (*MainHTTPHandler, error) {
	// Create database configuration
	dbConfig := &sharedDb.Config{
		Host:            cfg.GetString("DB_HOST"),
		Port:            cfg.GetInt("DB_PORT"),
		User:            cfg.GetString("DB_USER"),
		Password:        cfg.GetString("DB_PASSWORD"),
		DBName:          cfg.GetString("DB_NAME"),
		SSLMode:         cfg.GetString("DB_SSL_MODE"),
		MaxOpenConns:    25,
		MaxIdleConns:    5,
		ConnMaxLifetime: 5 * time.Minute,
		ConnMaxIdleTime: 1 * time.Minute,
		ConnectTimeout:  10 * time.Second,
		QueryTimeout:    30 * time.Second,
		MaxRetries:      3,
		RetryInterval:   2 * time.Second,
	}

	db, err := sharedDb.NewDatabaseHandler(dbConfig, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create database handler: %w", err)
	}

	// Create stock category handlers
	stockCategoryDBHandler, err := stockCategoryHandlers.NewDBHandler(db, logger)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create stock category handler: %w", err)
	}
	stockCategoryHTTPHandler := stockCategoryHandlers.NewHTTPHandler(stockCategoryDBHandler, logger)

	// Create stock item handlers
	stockItemDBHandler, err := stockItemHandlers.NewDBHandler(db, logger)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create stock item handler: %w", err)
	}
	stockItemHTTPHandler := stockItemHandlers.NewHTTPHandler(stockItemDBHandler, logger)

	// Create supplier handlers
	supplierDBHandler, err := supplierHandlers.NewDBHandler(db, logger)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create supplier handler: %w", err)
	}
	supplierHTTPHandler := supplierHandlers.NewHTTPHandler(supplierDBHandler, logger)

	// Create cancellable context for health monitor
	ctx, cancel := context.WithCancel(context.Background())

	//pvillalobos this should be configurable
	// Create HTTP health monitor for data-service
	httpHealthMonitor, err := sharedHttp.NewHealthMonitor(logger, 1*time.Second)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to create HTTP health monitor: %w", err)
	}
	httpHealthMonitor.AddService("data-service", sharedConfig.DATA_SERVICE_URL+"/api/v1/data/p/health")
	httpHealthMonitor.Start(ctx)

	return &MainHTTPHandler{
		db:                   db,
		httpHealthMonitor:    httpHealthMonitor,
		cancelHealthMonitor:  cancel,
		stockCategoryHandler: stockCategoryHTTPHandler,
		stockItemHandler:     stockItemHTTPHandler,
		supplierHandler:      supplierHTTPHandler,
		logger:               logger,
	}, nil
}

func (h *MainHTTPHandler) CloseDB() error {
	// Stop health monitor
	if h.cancelHealthMonitor != nil {
		h.cancelHealthMonitor()
	}

	err := h.db.Close()
	if err != nil {
		h.logger.WithError(err).Error("Failed to close database")
		return err
	}
	return nil
}

func (h *MainHTTPHandler) SetupRoutes(router *mux.Router) {
	// Health check
	router.HandleFunc("/api/v1/inventory/p/health", h.HealthCheck).Methods("GET")

	// Stock Item Categories
	router.HandleFunc("/api/v1/stock/categories", h.stockCategoryHandler.List).Methods("GET")
	router.HandleFunc("/api/v1/stock/categories/{id}", h.stockCategoryHandler.GetByID).Methods("GET")
	router.HandleFunc("/api/v1/stock/categories", h.stockCategoryHandler.Create).Methods("POST")
	router.HandleFunc("/api/v1/stock/categories/{id}", h.stockCategoryHandler.Update).Methods("PUT")
	router.HandleFunc("/api/v1/stock/categories/{id}", h.stockCategoryHandler.Delete).Methods("DELETE")

	// Stock Items
	router.HandleFunc("/api/v1/stock/items", h.stockItemHandler.List).Methods("GET")
	router.HandleFunc("/api/v1/stock/items/{id}", h.stockItemHandler.GetByID).Methods("GET")
	router.HandleFunc("/api/v1/stock/items", h.stockItemHandler.Create).Methods("POST")
	router.HandleFunc("/api/v1/stock/items/{id}", h.stockItemHandler.Update).Methods("PUT")
	router.HandleFunc("/api/v1/stock/items/{id}", h.stockItemHandler.Delete).Methods("DELETE")

	// Suppliers
	router.HandleFunc("/api/v1/inventory/suppliers", h.supplierHandler.List).Methods("GET")
	router.HandleFunc("/api/v1/inventory/suppliers/{id}", h.supplierHandler.GetByID).Methods("GET")
	router.HandleFunc("/api/v1/inventory/suppliers", h.supplierHandler.Create).Methods("POST")
	router.HandleFunc("/api/v1/inventory/suppliers/{id}", h.supplierHandler.Update).Methods("PUT")
	router.HandleFunc("/api/v1/inventory/suppliers/{id}", h.supplierHandler.Delete).Methods("DELETE")

}

func (h *MainHTTPHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"service":   "inventory-service",
		"timestamp": time.Now(),
	}

	// Check cached health state from background monitor
	healthStatus := h.httpHealthMonitor.GetHealthStatus()
	if !healthStatus.IsHealthy {
		response["status"] = "unhealthy"
		response["message"] = "Dependent services are not healthy"
		response["services"] = healthStatus.Services

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(response)
		return
	}

	response["status"] = "healthy"
	response["message"] = "Inventory service is healthy"
	response["services"] = healthStatus.Services

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
