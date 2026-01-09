package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	ingredientHandlers "menu-service/pkg/entities/ingredients/handlers"
	menuCategoryHandlers "menu-service/pkg/entities/menu_categories/handlers"
	menuItemHandlers "menu-service/pkg/entities/menu_items/handlers"
	stockCategoryHandlers "menu-service/pkg/entities/stock_item_categories/handlers"
	stockItemHandlers "menu-service/pkg/entities/stock_items/handlers"

	sharedConfig "shared/config"
	sharedDb "shared/db"
	sharedHttp "shared/http"
	sharedLogger "shared/logger"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type MainHTTPHandler struct {
	menuCategoryHandler  *menuCategoryHandlers.HTTPHandler
	menuItemHandler      *menuItemHandlers.HTTPHandler
	stockCategoryHandler *stockCategoryHandlers.HTTPHandler
	stockItemHandler     *stockItemHandlers.HTTPHandler
	ingredientHandler    *ingredientHandlers.HTTPHandler
	logger               *logrus.Logger
}

func main() {
	logger := sharedLogger.SetupLogger(sharedLogger.SERVICE_MENU_SERVICE, "INFO")
	logger.Info("🍽️ Starting Bar-Restaurant Menu Service")

	configLoader := sharedConfig.NewConfigLoader(sharedConfig.DATA_SERVICE_URL)
	config, err := configLoader.LoadConfig("menu", logger)
	if err != nil {
		logger.WithError(err).Fatal("Failed to load configuration")
	}

	logger.WithFields(logrus.Fields{
		"port": config.GetString("SERVER_PORT"),
		"host": config.GetString("SERVER_HOST"),
	}).Info("Configuration loaded")

	mainHandler, db, err := newMainHTTPHandler(config, logger)
	if err != nil {
		logger.WithError(err).Fatal("Failed to create HTTP handler")
	}
	defer db.Close()

	router := mux.NewRouter()
	mainHandler.SetupRoutes(router)

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", config.GetString("SERVER_HOST"), config.GetString("SERVER_PORT")),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.WithField("addr", server.Addr).Info("HTTP server starting")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.WithError(err).Fatal("Server failed")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.WithError(err).Fatal("Shutdown failed")
	}
	logger.Info("Menu Service stopped")
}

func newMainHTTPHandler(cfg *sharedConfig.Config, logger *logrus.Logger) (*MainHTTPHandler, *sharedDb.DbHandler, error) {
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
		return nil, nil, fmt.Errorf("failed to create database handler: %w", err)
	}

	// Create menu category handlers
	menuCategoryDBHandler, err := menuCategoryHandlers.NewDBHandler(db, logger)
	if err != nil {
		db.Close()
		return nil, nil, fmt.Errorf("failed to create menu category handler: %w", err)
	}
	menuCategoryHTTPHandler := menuCategoryHandlers.NewHTTPHandler(menuCategoryDBHandler, logger)

	// Create menu item handlers
	menuItemDBHandler, err := menuItemHandlers.NewDBHandler(db, logger)
	if err != nil {
		db.Close()
		return nil, nil, fmt.Errorf("failed to create menu item handler: %w", err)
	}
	menuItemHTTPHandler := menuItemHandlers.NewHTTPHandler(menuItemDBHandler, logger)

	// Create stock category handlers
	stockCategoryDBHandler, err := stockCategoryHandlers.NewDBHandler(db, logger)
	if err != nil {
		db.Close()
		return nil, nil, fmt.Errorf("failed to create stock category handler: %w", err)
	}
	stockCategoryHTTPHandler := stockCategoryHandlers.NewHTTPHandler(stockCategoryDBHandler, logger)

	// Create stock item handlers
	stockItemDBHandler, err := stockItemHandlers.NewDBHandler(db, logger)
	if err != nil {
		db.Close()
		return nil, nil, fmt.Errorf("failed to create stock item handler: %w", err)
	}
	stockItemHTTPHandler := stockItemHandlers.NewHTTPHandler(stockItemDBHandler, logger)

	// Create ingredient handlers
	ingredientDBHandler, err := ingredientHandlers.NewDBHandler(db, logger)
	if err != nil {
		db.Close()
		return nil, nil, fmt.Errorf("failed to create ingredient handler: %w", err)
	}
	ingredientHTTPHandler := ingredientHandlers.NewHTTPHandler(ingredientDBHandler, menuItemDBHandler, logger)

	return &MainHTTPHandler{
		menuCategoryHandler:  menuCategoryHTTPHandler,
		menuItemHandler:      menuItemHTTPHandler,
		stockCategoryHandler: stockCategoryHTTPHandler,
		stockItemHandler:     stockItemHTTPHandler,
		ingredientHandler:    ingredientHTTPHandler,
		logger:               logger,
	}, db, nil
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

	// Menu Items
	router.HandleFunc("/api/v1/menu/items", h.menuItemHandler.List).Methods("GET")
	router.HandleFunc("/api/v1/menu/items/{id}", h.menuItemHandler.GetByID).Methods("GET")
	router.HandleFunc("/api/v1/menu/items", h.menuItemHandler.Create).Methods("POST")
	router.HandleFunc("/api/v1/menu/items/{id}", h.menuItemHandler.Update).Methods("PUT")
	router.HandleFunc("/api/v1/menu/items/{id}", h.menuItemHandler.Delete).Methods("DELETE")
	router.HandleFunc("/api/v1/menu/items/{id}/availability", h.menuItemHandler.UpdateAvailability).Methods("PATCH")

	// Menu Item Ingredients
	router.HandleFunc("/api/v1/menu/items/{id}/ingredients", h.ingredientHandler.List).Methods("GET")
	router.HandleFunc("/api/v1/menu/items/{id}/ingredients", h.ingredientHandler.Create).Methods("POST")
	router.HandleFunc("/api/v1/menu/items/{id}/ingredients/{stockItemId}", h.ingredientHandler.Update).Methods("PUT")
	router.HandleFunc("/api/v1/menu/items/{id}/ingredients/{stockItemId}", h.ingredientHandler.Delete).Methods("DELETE")

	// Cost Calculation
	router.HandleFunc("/api/v1/menu/items/{id}/cost", h.ingredientHandler.GetCost).Methods("GET")
	router.HandleFunc("/api/v1/menu/items/{id}/cost/recalculate", h.ingredientHandler.RecalculateCost).Methods("POST")

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
}

func (h *MainHTTPHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	client := &http.Client{Timeout: 1 * time.Second}
	resp, err := client.Get(sharedConfig.DATA_SERVICE_URL + "/api/v1/data/p/health")
	if err != nil {
		h.logger.WithError(err).Error("data-service is not healthy")
		sharedHttp.SendError(w, http.StatusServiceUnavailable, "data-service is not healthy", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		sharedHttp.SendError(w, http.StatusServiceUnavailable, "data-service is not healthy", nil)
		return
	}

	sharedHttp.SendSuccess(w, http.StatusOK, "Menu service healthy", map[string]interface{}{
		"status":  "healthy",
		"service": "menu-service",
	})
}
