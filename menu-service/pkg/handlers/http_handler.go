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
	menuIngredientHandlers "menu-service/pkg/entities/menu_ingredients/handlers"
	menuSubCategoryHandlers "menu-service/pkg/entities/menu_sub_categories/handlers"
	menuVariantHandlers "menu-service/pkg/entities/menu_variants/handlers"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type MainHTTPHandler struct {
	db                     *sharedDb.DbHandler
	httpHealthMonitor      *sharedHttp.HTTPHealthMonitor
	cancelHealthMonitor    context.CancelFunc
	menuCategoryHandler    *menuCategoryHandlers.HTTPHandler
	menuSubCategoryHandler *menuSubCategoryHandlers.HTTPHandler
	menuVariantHandler     *menuVariantHandlers.HTTPHandler
	menuIngredientHandler  *menuIngredientHandlers.HTTPHandler
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

	// Create menu category handlers
	menuCategoryDBHandler, err := menuCategoryHandlers.NewDBHandler(db, logger)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create menu category handler: %w", err)
	}
	menuCategoryHTTPHandler := menuCategoryHandlers.NewHTTPHandler(menuCategoryDBHandler, logger)

	// Create menu sub-category handlers
	menuSubCategoryDBHandler, err := menuSubCategoryHandlers.NewDBHandler(db, logger)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create menu sub-category handler: %w", err)
	}
	menuSubCategoryHTTPHandler := menuSubCategoryHandlers.NewHTTPHandler(menuSubCategoryDBHandler, logger)

	// Create menu variant handlers
	menuVariantDBHandler, err := menuVariantHandlers.NewDBHandler(db, logger)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create menu variant handler: %w", err)
	}
	menuVariantHTTPHandler := menuVariantHandlers.NewHTTPHandler(menuVariantDBHandler, logger)

	// Create menu ingredient handlers
	menuIngredientDBHandler, err := menuIngredientHandlers.NewDBHandler(db, logger)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create menu ingredient handler: %w", err)
	}
	menuIngredientHTTPHandler := menuIngredientHandlers.NewHTTPHandler(menuIngredientDBHandler, logger)

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
		db:                     db,
		httpHealthMonitor:      httpHealthMonitor,
		cancelHealthMonitor:    cancel,
		menuCategoryHandler:    menuCategoryHTTPHandler,
		menuSubCategoryHandler: menuSubCategoryHTTPHandler,
		menuVariantHandler:     menuVariantHTTPHandler,
		menuIngredientHandler:  menuIngredientHTTPHandler,
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
	router.HandleFunc("/api/v1/menu/p/health", h.HealthCheck).Methods("GET")

	// Menu Categories
	router.HandleFunc("/api/v1/menu/categories", h.menuCategoryHandler.List).Methods("GET")
	router.HandleFunc("/api/v1/menu/categories/{id}", h.menuCategoryHandler.GetByID).Methods("GET")
	router.HandleFunc("/api/v1/menu/categories", h.menuCategoryHandler.Create).Methods("POST")
	router.HandleFunc("/api/v1/menu/categories/{id}", h.menuCategoryHandler.Update).Methods("PUT")
	router.HandleFunc("/api/v1/menu/categories/{id}", h.menuCategoryHandler.Delete).Methods("DELETE")

	// Menu Sub-Categories
	router.HandleFunc("/api/v1/menu/sub-categories", h.menuSubCategoryHandler.List).Methods("GET")
	router.HandleFunc("/api/v1/menu/sub-categories/{id}", h.menuSubCategoryHandler.GetByID).Methods("GET")
	router.HandleFunc("/api/v1/menu/sub-categories", h.menuSubCategoryHandler.Create).Methods("POST")
	router.HandleFunc("/api/v1/menu/sub-categories/{id}", h.menuSubCategoryHandler.Update).Methods("PUT")
	router.HandleFunc("/api/v1/menu/sub-categories/{id}", h.menuSubCategoryHandler.Delete).Methods("DELETE")

	// Menu Variants
	router.HandleFunc("/api/v1/menu/variants", h.menuVariantHandler.List).Methods("GET")
	router.HandleFunc("/api/v1/menu/variants/{id}", h.menuVariantHandler.GetByID).Methods("GET")
	router.HandleFunc("/api/v1/menu/variants", h.menuVariantHandler.Create).Methods("POST")
	router.HandleFunc("/api/v1/menu/variants/{id}", h.menuVariantHandler.Update).Methods("PUT")
	router.HandleFunc("/api/v1/menu/variants/{id}", h.menuVariantHandler.Delete).Methods("DELETE")
	router.HandleFunc("/api/v1/menu/variants/{id}/availability", h.menuVariantHandler.UpdateAvailability).Methods("PATCH")

	// Menu Ingredients
	router.HandleFunc("/api/v1/menu/ingredients", h.menuIngredientHandler.List).Methods("GET")
	router.HandleFunc("/api/v1/menu/ingredients/{id}", h.menuIngredientHandler.GetByID).Methods("GET")
	router.HandleFunc("/api/v1/menu/ingredients", h.menuIngredientHandler.Create).Methods("POST")
	router.HandleFunc("/api/v1/menu/ingredients/{id}", h.menuIngredientHandler.Update).Methods("PUT")
	router.HandleFunc("/api/v1/menu/ingredients/{id}", h.menuIngredientHandler.Delete).Methods("DELETE")
	router.HandleFunc("/api/v1/menu/variants/{variantId}/ingredients", h.menuIngredientHandler.GetByMenuVariant).Methods("GET")
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
