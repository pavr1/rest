package handlers

import (
	"encoding/json"
	"net/http"

	"menu-service/pkg/entities/ingredients/models"
	menuItemHandlers "menu-service/pkg/entities/menu_items/handlers"
	sharedHttp "shared/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// HTTPHandler handles HTTP requests for ingredients
type HTTPHandler struct {
	dbHandler       *DBHandler
	menuItemHandler *menuItemHandlers.DBHandler
	logger          *logrus.Logger
}

// NewHTTPHandler creates a new HTTP handler
func NewHTTPHandler(dbHandler *DBHandler, menuItemHandler *menuItemHandlers.DBHandler, logger *logrus.Logger) *HTTPHandler {
	return &HTTPHandler{
		dbHandler:       dbHandler,
		menuItemHandler: menuItemHandler,
		logger:          logger,
	}
}

// List handles GET /api/v1/menu/items/:id/ingredients
func (h *HTTPHandler) List(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	menuItemID := vars["id"]

	response, err := h.dbHandler.List(menuItemID)
	if err != nil {
		h.logger.WithError(err).Error("Failed to list ingredients")
		sharedHttp.SendError(w, http.StatusInternalServerError, "Failed to list ingredients", err)
		return
	}

	sharedHttp.SendSuccess(w, http.StatusOK, "Ingredients retrieved", response)
}

// Create handles POST /api/v1/menu/items/:id/ingredients
func (h *HTTPHandler) Create(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	menuItemID := vars["id"]

	var req models.IngredientCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.WithError(err).Error("Failed to decode request body")
		sharedHttp.SendError(w, http.StatusBadRequest, "Invalid request format", err)
		return
	}

	if req.StockItemID == "" {
		sharedHttp.SendError(w, http.StatusBadRequest, "Stock item ID is required", nil)
		return
	}

	if req.Quantity <= 0 {
		sharedHttp.SendError(w, http.StatusBadRequest, "Quantity must be greater than 0", nil)
		return
	}

	ingredient, err := h.dbHandler.Create(menuItemID, &req)
	if err != nil {
		h.logger.WithError(err).Error("Failed to add ingredient")
		sharedHttp.SendError(w, http.StatusInternalServerError, "Failed to add ingredient", err)
		return
	}

	sharedHttp.SendSuccess(w, http.StatusCreated, "Ingredient added", ingredient)
}

// Update handles PUT /api/v1/menu/items/:id/ingredients/:stockItemId
func (h *HTTPHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	menuItemID := vars["id"]
	stockItemID := vars["stockItemId"]

	var req models.IngredientUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.WithError(err).Error("Failed to decode request body")
		sharedHttp.SendError(w, http.StatusBadRequest, "Invalid request format", err)
		return
	}

	if req.Quantity <= 0 {
		sharedHttp.SendError(w, http.StatusBadRequest, "Quantity must be greater than 0", nil)
		return
	}

	ingredient, err := h.dbHandler.Update(menuItemID, stockItemID, &req)
	if err != nil {
		h.logger.WithError(err).Error("Failed to update ingredient")
		sharedHttp.SendError(w, http.StatusInternalServerError, "Failed to update ingredient", err)
		return
	}

	if ingredient == nil {
		sharedHttp.SendError(w, http.StatusNotFound, "Ingredient not found", nil)
		return
	}

	sharedHttp.SendSuccess(w, http.StatusOK, "Ingredient updated", ingredient)
}

// Delete handles DELETE /api/v1/menu/items/:id/ingredients/:stockItemId
func (h *HTTPHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	menuItemID := vars["id"]
	stockItemID := vars["stockItemId"]

	err := h.dbHandler.Delete(menuItemID, stockItemID)
	if err != nil {
		h.logger.WithError(err).Error("Failed to delete ingredient")
		if err.Error() == "ingredient not found" {
			sharedHttp.SendError(w, http.StatusNotFound, "Ingredient not found", nil)
			return
		}
		sharedHttp.SendError(w, http.StatusInternalServerError, "Failed to delete ingredient", err)
		return
	}

	sharedHttp.SendSuccess(w, http.StatusOK, "Ingredient removed", nil)
}

// GetCost handles GET /api/v1/menu/items/:id/cost
func (h *HTTPHandler) GetCost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	menuItemID := vars["id"]

	// Get menu item name
	menuItem, err := h.menuItemHandler.GetByID(menuItemID)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get menu item")
		sharedHttp.SendError(w, http.StatusInternalServerError, "Failed to get menu item", err)
		return
	}

	if menuItem == nil {
		sharedHttp.SendError(w, http.StatusNotFound, "Menu item not found", nil)
		return
	}

	cost, err := h.dbHandler.CalculateCost(menuItemID, menuItem.Name)
	if err != nil {
		h.logger.WithError(err).Error("Failed to calculate cost")
		sharedHttp.SendError(w, http.StatusInternalServerError, "Failed to calculate cost", err)
		return
	}

	sharedHttp.SendSuccess(w, http.StatusOK, "Cost calculated", cost)
}

// RecalculateCost handles POST /api/v1/menu/items/:id/cost/recalculate
func (h *HTTPHandler) RecalculateCost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	menuItemID := vars["id"]

	// Get menu item
	menuItem, err := h.menuItemHandler.GetByID(menuItemID)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get menu item")
		sharedHttp.SendError(w, http.StatusInternalServerError, "Failed to get menu item", err)
		return
	}

	if menuItem == nil {
		sharedHttp.SendError(w, http.StatusNotFound, "Menu item not found", nil)
		return
	}

	// Calculate cost
	cost, err := h.dbHandler.CalculateCost(menuItemID, menuItem.Name)
	if err != nil {
		h.logger.WithError(err).Error("Failed to calculate cost")
		sharedHttp.SendError(w, http.StatusInternalServerError, "Failed to calculate cost", err)
		return
	}

	// Update the menu item with the new cost
	updatedItem, err := h.menuItemHandler.UpdateCost(menuItemID, cost.TotalCost)
	if err != nil {
		h.logger.WithError(err).Error("Failed to update menu item cost")
		sharedHttp.SendError(w, http.StatusInternalServerError, "Failed to update cost", err)
		return
	}

	h.logger.WithFields(map[string]interface{}{
		"menu_item_id": menuItemID,
		"new_cost":     cost.TotalCost,
	}).Info("Menu item cost recalculated")

	sharedHttp.SendSuccess(w, http.StatusOK, "Cost recalculated and saved", map[string]interface{}{
		"menu_item": updatedItem,
		"cost":      cost,
	})
}
