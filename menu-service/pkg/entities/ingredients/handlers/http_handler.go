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
		sharedHttp.SendErrorResponse(w, http.StatusInternalServerError, "Failed to list ingredients")
		return
	}

	sharedHttp.SendSuccessResponse(w, http.StatusOK, "Ingredients retrieved", response)
}

// Create handles POST /api/v1/menu/items/:id/ingredients
func (h *HTTPHandler) Create(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	menuItemID := vars["id"]

	var req models.IngredientCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.WithError(err).Error("Failed to decode request body")
		sharedHttp.SendErrorResponse(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	if req.StockItemID == "" {
		sharedHttp.SendErrorResponse(w, http.StatusBadRequest, "Stock item ID is required")
		return
	}

	if req.Quantity <= 0 {
		sharedHttp.SendErrorResponse(w, http.StatusBadRequest, "Quantity must be greater than 0")
		return
	}

	ingredient, err := h.dbHandler.Create(menuItemID, &req)
	if err != nil {
		h.logger.WithError(err).Error("Failed to add ingredient")
		sharedHttp.SendErrorResponse(w, http.StatusInternalServerError, "Failed to add ingredient")
		return
	}

	sharedHttp.SendSuccessResponse(w, http.StatusCreated, "Ingredient added", ingredient)
}

// Update handles PUT /api/v1/menu/items/:id/ingredients/:stockItemId
func (h *HTTPHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	menuItemID := vars["id"]
	stockItemID := vars["stockItemId"]

	var req models.IngredientUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.WithError(err).Error("Failed to decode request body")
		sharedHttp.SendErrorResponse(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	if req.Quantity <= 0 {
		sharedHttp.SendErrorResponse(w, http.StatusBadRequest, "Quantity must be greater than 0")
		return
	}

	ingredient, err := h.dbHandler.Update(menuItemID, stockItemID, &req)
	if err != nil {
		h.logger.WithError(err).Error("Failed to update ingredient")
		sharedHttp.SendErrorResponse(w, http.StatusInternalServerError, "Failed to update ingredient")
		return
	}

	if ingredient == nil {
		sharedHttp.SendErrorResponse(w, http.StatusNotFound, "Ingredient not found")
		return
	}

	sharedHttp.SendSuccessResponse(w, http.StatusOK, "Ingredient updated", ingredient)
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
			sharedHttp.SendErrorResponse(w, http.StatusNotFound, "Ingredient not found")
			return
		}
		sharedHttp.SendErrorResponse(w, http.StatusInternalServerError, "Failed to delete ingredient")
		return
	}

	sharedHttp.SendSuccessResponse(w, http.StatusOK, "Ingredient removed", nil)
}

// GetCost handles GET /api/v1/menu/items/:id/cost
func (h *HTTPHandler) GetCost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	menuItemID := vars["id"]

	// Get menu item name
	menuItem, err := h.menuItemHandler.GetByID(menuItemID)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get menu item")
		sharedHttp.SendErrorResponse(w, http.StatusInternalServerError, "Failed to get menu item")
		return
	}

	if menuItem == nil {
		sharedHttp.SendErrorResponse(w, http.StatusNotFound, "Menu item not found")
		return
	}

	cost, err := h.dbHandler.CalculateCost(menuItemID, menuItem.Name)
	if err != nil {
		h.logger.WithError(err).Error("Failed to calculate cost")
		sharedHttp.SendErrorResponse(w, http.StatusInternalServerError, "Failed to calculate cost")
		return
	}

	sharedHttp.SendSuccessResponse(w, http.StatusOK, "Cost calculated", cost)
}

// RecalculateCost handles POST /api/v1/menu/items/:id/cost/recalculate
func (h *HTTPHandler) RecalculateCost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	menuItemID := vars["id"]

	// Get menu item
	menuItem, err := h.menuItemHandler.GetByID(menuItemID)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get menu item")
		sharedHttp.SendErrorResponse(w, http.StatusInternalServerError, "Failed to get menu item")
		return
	}

	if menuItem == nil {
		sharedHttp.SendErrorResponse(w, http.StatusNotFound, "Menu item not found")
		return
	}

	// Calculate cost
	cost, err := h.dbHandler.CalculateCost(menuItemID, menuItem.Name)
	if err != nil {
		h.logger.WithError(err).Error("Failed to calculate cost")
		sharedHttp.SendErrorResponse(w, http.StatusInternalServerError, "Failed to calculate cost")
		return
	}

	// Update the menu item with the new cost
	updatedItem, err := h.menuItemHandler.UpdateCost(menuItemID, cost.TotalCost)
	if err != nil {
		h.logger.WithError(err).Error("Failed to update menu item cost")
		sharedHttp.SendErrorResponse(w, http.StatusInternalServerError, "Failed to update cost")
		return
	}

	h.logger.WithFields(map[string]interface{}{
		"menu_item_id": menuItemID,
		"new_cost":     cost.TotalCost,
	}).Info("Menu item cost recalculated")

	sharedHttp.SendSuccessResponse(w, http.StatusOK, "Cost recalculated and saved", map[string]interface{}{
		"menu_item": updatedItem,
		"cost":      cost,
	})
}
