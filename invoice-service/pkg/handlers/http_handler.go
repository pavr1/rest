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

	purchaseInvoiceHandlers "invoice-service/pkg/entities/purchase_invoices/handlers"
	invoiceDetailHandlers "invoice-service/pkg/entities/invoice_details/handlers"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type MainHTTPHandler struct {
	db                     *sharedDb.DbHandler
	httpHealthMonitor      *sharedHttp.HTTPHealthMonitor
	cancelHealthMonitor    context.CancelFunc
	purchaseInvoiceHandler *purchaseInvoiceHandlers.HTTPHandler
	invoiceDetailHandler   *invoiceDetailHandlers.HTTPHandler
	logger                 *logrus.Logger
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

	// Create purchase invoice handlers
	purchaseInvoiceDBHandler, err := purchaseInvoiceHandlers.NewDBHandler(db, logger)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create purchase invoice handler: %w", err)
	}
	purchaseInvoiceHTTPHandler := purchaseInvoiceHandlers.NewHTTPHandler(purchaseInvoiceDBHandler, logger)

	// Create invoice detail handlers
	invoiceDetailDBHandler, err := invoiceDetailHandlers.NewDBHandler(db, logger)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create invoice detail handler: %w", err)
	}
	invoiceDetailHTTPHandler := invoiceDetailHandlers.NewHTTPHandler(invoiceDetailDBHandler, logger)


	// Create cancellable context for health monitor
	ctx, cancel := context.WithCancel(context.Background())

	// Create HTTP health monitor for data-service
	httpHealthMonitor, err := sharedHttp.NewHealthMonitor(logger, 1*time.Second)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to create HTTP health monitor: %w", err)
	}
	httpHealthMonitor.AddService("data-service", sharedConfig.DATA_SERVICE_URL+"/api/v1/data/p/health")
	httpHealthMonitor.Start(ctx)

	return &MainHTTPHandler{
		db:                     db,
		httpHealthMonitor:      httpHealthMonitor,
		cancelHealthMonitor:    cancel,
		purchaseInvoiceHandler: purchaseInvoiceHTTPHandler,
		invoiceDetailHandler:   invoiceDetailHTTPHandler,
		logger:                 logger,
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
	router.HandleFunc("/api/v1/invoices/p/health", h.HealthCheck).Methods("GET")

	// Purchase Invoices
	router.HandleFunc("/api/v1/invoices/purchase", h.purchaseInvoiceHandler.List).Methods("GET")
	router.HandleFunc("/api/v1/invoices/purchase", h.purchaseInvoiceHandler.Create).Methods("POST")
	router.HandleFunc("/api/v1/invoices/purchase/{id}", h.purchaseInvoiceHandler.GetByID).Methods("GET")
	router.HandleFunc("/api/v1/invoices/purchase/{id}", h.purchaseInvoiceHandler.Update).Methods("PUT")
	router.HandleFunc("/api/v1/invoices/purchase/{id}", h.purchaseInvoiceHandler.Delete).Methods("DELETE")

	// Invoice Details
	router.HandleFunc("/api/v1/invoices/purchase/{invoiceId}/details", h.invoiceDetailHandler.ListByInvoice).Methods("GET")
	router.HandleFunc("/api/v1/invoices/purchase/{invoiceId}/details", h.invoiceDetailHandler.Create).Methods("POST")
	router.HandleFunc("/api/v1/invoices/purchase/{invoiceId}/details/{id}", h.invoiceDetailHandler.GetByID).Methods("GET")
	router.HandleFunc("/api/v1/invoices/purchase/{invoiceId}/details/{id}", h.invoiceDetailHandler.Update).Methods("PUT")
	router.HandleFunc("/api/v1/invoices/purchase/{invoiceId}/details/{id}", h.invoiceDetailHandler.Delete).Methods("DELETE")
}

func (h *MainHTTPHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"service":   "invoice-service",
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
	response["message"] = "Invoice service is healthy"
	response["services"] = healthStatus.Services

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
