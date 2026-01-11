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

	menuCategoryHandlers "menu-service/pkg/entities/menu_categories/handlers"
	menuItemHandlers "menu-service/pkg/entities/menu_items/handlers"
	subMenuHandlers "menu-service/pkg/entities/sub_menus/handlers"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type MainHTTPHandler struct {
	db                  *sharedDb.DbHandler
	httpHealthMonitor   *sharedHttp.HTTPHealthMonitor
	cancelHealthMonitor context.CancelFunc
	menuCategoryHandler *menuCategoryHandlers.HTTPHandler
	subMenuHandler      *subMenuHandlers.HTTPHandler
	menuItemHandler     *menuItemHandlers.HTTPHandler
	logger              *logrus.Logger
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

	// Create menu category handlers
	menuCategoryDBHandler, err := menuCategoryHandlers.NewDBHandler(db, logger)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create menu category handler: %w", err)
	}
	menuCategoryHTTPHandler := menuCategoryHandlers.NewHTTPHandler(menuCategoryDBHandler, logger)

	// Create sub menu handlers
	subMenuDBHandler, err := subMenuHandlers.NewDBHandler(db, logger)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create sub menu handler: %w", err)
	}
	subMenuHTTPHandler := subMenuHandlers.NewHTTPHandler(subMenuDBHandler, logger)

	// Create menu item handlers
	menuItemDBHandler, err := menuItemHandlers.NewDBHandler(db, logger)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create menu item handler: %w", err)
	}
	menuItemHTTPHandler := menuItemHandlers.NewHTTPHandler(menuItemDBHandler, logger)

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
		db:                  db,
		httpHealthMonitor:   httpHealthMonitor,
		cancelHealthMonitor: cancel,
		menuCategoryHandler: menuCategoryHTTPHandler,
		subMenuHandler:      subMenuHTTPHandler,
		menuItemHandler:     menuItemHTTPHandler,
		logger:              logger,
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
	router.HandleFunc("/api/v1/menu/p/health", h.HealthCheck).Methods("GET")

	// Menu Categories
	router.HandleFunc("/api/v1/menu/categories", h.menuCategoryHandler.List).Methods("GET")
	router.HandleFunc("/api/v1/menu/categories/{id}", h.menuCategoryHandler.GetByID).Methods("GET")
	router.HandleFunc("/api/v1/menu/categories", h.menuCategoryHandler.Create).Methods("POST")
	router.HandleFunc("/api/v1/menu/categories/{id}", h.menuCategoryHandler.Update).Methods("PUT")
	router.HandleFunc("/api/v1/menu/categories/{id}", h.menuCategoryHandler.Delete).Methods("DELETE")

	// Sub Menus
	router.HandleFunc("/api/v1/menu/submenus", h.subMenuHandler.List).Methods("GET")
	router.HandleFunc("/api/v1/menu/submenus/{id}", h.subMenuHandler.GetByID).Methods("GET")
	router.HandleFunc("/api/v1/menu/submenus", h.subMenuHandler.Create).Methods("POST")
	router.HandleFunc("/api/v1/menu/submenus/{id}", h.subMenuHandler.Update).Methods("PUT")
	router.HandleFunc("/api/v1/menu/submenus/{id}", h.subMenuHandler.Delete).Methods("DELETE")

	// Menu Items
	router.HandleFunc("/api/v1/menu/items", h.menuItemHandler.List).Methods("GET")
	router.HandleFunc("/api/v1/menu/items/{id}", h.menuItemHandler.GetByID).Methods("GET")
	router.HandleFunc("/api/v1/menu/items", h.menuItemHandler.Create).Methods("POST")
	router.HandleFunc("/api/v1/menu/items/{id}", h.menuItemHandler.Update).Methods("PUT")
	router.HandleFunc("/api/v1/menu/items/{id}", h.menuItemHandler.Delete).Methods("DELETE")
	router.HandleFunc("/api/v1/menu/items/{id}/availability", h.menuItemHandler.UpdateAvailability).Methods("PATCH")
}

func (h *MainHTTPHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"service":   "menu-service",
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
	response["message"] = "Menu service is healthy"
	response["services"] = healthStatus.Services

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
