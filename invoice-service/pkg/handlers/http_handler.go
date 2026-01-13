package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	sharedConfig "shared/config"
	sharedDb "shared/db"
	sharedHttp "shared/http"

	incomeInvoiceHandlers "invoice-service/pkg/entities/income_invoices/handlers"
	"invoice-service/pkg/entities/income_invoices/models"
	outcomeInvoiceHandlers "invoice-service/pkg/entities/outcome_invoices/handlers"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type MainHTTPHandler struct {
	db                    *sharedDb.DbHandler
	httpHealthMonitor     *sharedHttp.HTTPHealthMonitor
	cancelHealthMonitor   context.CancelFunc
	outcomeInvoiceHandler *outcomeInvoiceHandlers.HTTPHandler
	incomeInvoiceHandler  *incomeInvoiceHandlers.HTTPHandler
	logger                *logrus.Logger
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

	// Create outcome invoice handlers
	outcomeInvoiceDBHandler, err := outcomeInvoiceHandlers.NewDBHandler(db, logger)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create outcome invoice handler: %w", err)
	}
	outcomeInvoiceHTTPHandler := outcomeInvoiceHandlers.NewHTTPHandler(outcomeInvoiceDBHandler, logger)

	// Create income invoice handlers
	incomeInvoiceDBHandler, err := incomeInvoiceHandlers.NewDBHandler(db, logger)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create income invoice handler: %w", err)
	}
	incomeInvoiceHTTPHandler := incomeInvoiceHandlers.NewHTTPHandler(incomeInvoiceDBHandler, logger)

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
		db:                    db,
		httpHealthMonitor:     httpHealthMonitor,
		cancelHealthMonitor:   cancel,
		outcomeInvoiceHandler: outcomeInvoiceHTTPHandler,
		incomeInvoiceHandler:  incomeInvoiceHTTPHandler,
		logger:                logger,
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

	// Outcome Invoices (expenses from suppliers)
	router.HandleFunc("/api/v1/invoices/outcome", h.outcomeInvoiceHandler.List).Methods("GET")
	router.HandleFunc("/api/v1/invoices/outcome", h.outcomeInvoiceHandler.Create).Methods("POST")
	router.HandleFunc("/api/v1/invoices/outcome/{id}", h.GetOutcomeInvoiceByID).Methods("GET")
	router.HandleFunc("/api/v1/invoices/outcome/{id}", h.outcomeInvoiceHandler.Update).Methods("PUT")
	router.HandleFunc("/api/v1/invoices/outcome/{id}", h.outcomeInvoiceHandler.Delete).Methods("DELETE")

	// Income Invoices (revenue from customers)
	router.HandleFunc("/api/v1/invoices/income", h.incomeInvoiceHandler.List).Methods("GET")
	router.HandleFunc("/api/v1/invoices/income", h.incomeInvoiceHandler.Create).Methods("POST")
	router.HandleFunc("/api/v1/invoices/income/{id}", h.GetIncomeInvoiceByID).Methods("GET")
	router.HandleFunc("/api/v1/invoices/income/{id}", h.incomeInvoiceHandler.Update).Methods("PUT")
	router.HandleFunc("/api/v1/invoices/income/{id}", h.incomeInvoiceHandler.Delete).Methods("DELETE")

	// Invoice Items are now handled within invoice CRUD operations
	// No separate endpoints for invoice items
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

// GetOutcomeInvoiceByID gets an outcome invoice with its invoice items
func (h *MainHTTPHandler) GetOutcomeInvoiceByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Get the invoice
	invoice, err := h.outcomeInvoiceHandler.GetByIDOnly(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			sharedHttp.SendErrorResponse(w, http.StatusNotFound, "Outcome invoice not found")
			return
		}
		h.logger.WithError(err).Error("Failed to get outcome invoice")
		sharedHttp.SendErrorResponse(w, http.StatusInternalServerError, "Failed to get outcome invoice")
		return
	}

	// Get invoice items for this invoice
	itemsResponse, err := h.outcomeInvoiceHandler.GetByIDOnly(id)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get invoice items")
		// Don't fail the request, just return empty items
		invoice.InvoiceItems = []models.InvoiceItem{}
	} else {
		invoice.InvoiceItems = itemsResponse.InvoiceItems
	}

	sharedHttp.SendSuccessResponse(w, http.StatusOK, "Outcome invoice retrieved successfully", invoice)
}

// GetIncomeInvoiceByID gets an income invoice with its invoice items
func (h *MainHTTPHandler) GetIncomeInvoiceByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Get the invoice
	invoice, err := h.incomeInvoiceHandler.GetByIDOnly(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			sharedHttp.SendErrorResponse(w, http.StatusNotFound, "Income invoice not found")
			return
		}
		h.logger.WithError(err).Error("Failed to get income invoice")
		sharedHttp.SendErrorResponse(w, http.StatusInternalServerError, "Failed to get income invoice")
		return
	}

	// Get invoice items for this invoice
	itemsResponse, err := h.incomeInvoiceHandler.GetByIDOnly(id)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get invoice items")
		// Don't fail the request, just return empty items
		invoice.InvoiceItems = []models.InvoiceItem{}
	} else {
		invoice.InvoiceItems = itemsResponse.InvoiceItems
	}

	sharedHttp.SendSuccessResponse(w, http.StatusOK, "Income invoice retrieved successfully", invoice)
}
