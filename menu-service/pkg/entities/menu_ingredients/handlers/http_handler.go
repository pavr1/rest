package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"menu-service/pkg/entities/menu_ingredients/models"
	sharedHttp "shared/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// HTTPHandler handles HTTP requests for menu ingredients
type HTTPHandler struct {
	db     *DBHandler
	logger *logrus.Logger
}

// NewHTTPHandler creates a new HTTP handler
func NewHTTPHandler(db *DBHandler, logger *logrus.Logger) *HTTPHandler {
	return &HTTPHandler{
		db:     db,
		logger: logger,
	}
}

// List handles GET /api/v1/menu/ingredients
func (h *HTTPHandler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	ingredients, err := h.db.List(page, limit)
	if err != nil {
		h.logger.WithError(err).Error("Failed to list menu ingredients")
		sharedHttp.SendErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve menu ingredients")
		return
	}

	response := map[string]interface{}{
		"ingredients": ingredients,
		"total":       len(ingredients), // Simplified - in real implementation you'd have a count query
		"page":        page,
		"limit":       limit,
	}

	sharedHttp.SendSuccessResponse(w, http.StatusOK, "Menu ingredients retrieved", response)
}

// GetByID handles GET /api/v1/menu/ingredients/{id}
func (h *HTTPHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	ingredient, err := h.db.GetByID(id)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get menu ingredient by ID")
		sharedHttp.SendErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve menu ingredient")
		return
	}

	if ingredient == nil {
		sharedHttp.SendErrorResponse(w, http.StatusNotFound, "Menu ingredient not found")
		return
	}

	sharedHttp.SendSuccessResponse(w, http.StatusOK, "Menu ingredient retrieved", ingredient)
}

// Create handles POST /api/v1/menu/ingredients
func (h *HTTPHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.MenuIngredientCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.WithError(err).Error("Failed to decode create menu ingredient request")
		sharedHttp.SendErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request - must have exactly one of stock_variant_id or menu_sub_category_id
	if err := req.Validate(); err != nil {
		sharedHttp.SendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Get menu variant ID from query parameter
	menuVariantID := r.URL.Query().Get("menu_variant_id")
	if menuVariantID == "" {
		sharedHttp.SendErrorResponse(w, http.StatusBadRequest, "menu_variant_id query parameter is required")
		return
	}

	ingredient, err := h.db.Create(req, menuVariantID)
	if err != nil {
		h.logger.WithError(err).Error("Failed to create menu ingredient")
		sharedHttp.SendErrorResponse(w, http.StatusInternalServerError, "Failed to create menu ingredient")
		return
	}

	sharedHttp.SendSuccessResponse(w, http.StatusCreated, "Menu ingredient created", ingredient)
}

// Update handles PUT /api/v1/menu/ingredients/{id}
func (h *HTTPHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req models.MenuIngredientUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.WithError(err).Error("Failed to decode update menu ingredient request")
		sharedHttp.SendErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	ingredient, err := h.db.Update(id, req)
	if err != nil {
		h.logger.WithError(err).Error("Failed to update menu ingredient")
		sharedHttp.SendErrorResponse(w, http.StatusInternalServerError, "Failed to update menu ingredient")
		return
	}

	sharedHttp.SendSuccessResponse(w, http.StatusOK, "Menu ingredient updated", ingredient)
}

// Delete handles DELETE /api/v1/menu/ingredients/{id}
func (h *HTTPHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.db.Delete(id); err != nil {
		h.logger.WithError(err).Error("Failed to delete menu ingredient")
		sharedHttp.SendErrorResponse(w, http.StatusInternalServerError, "Failed to delete menu ingredient")
		return
	}

	sharedHttp.SendSuccessResponse(w, http.StatusOK, "Menu ingredient deleted successfully", nil)
}

// GetByMenuVariant handles GET /api/v1/menu/variants/{variantId}/ingredients
func (h *HTTPHandler) GetByMenuVariant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	menuVariantID := vars["variantId"]

	ingredients, err := h.db.GetByMenuVariant(menuVariantID)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get ingredients by menu variant")
		sharedHttp.SendErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve menu ingredients")
		return
	}

	response := models.MenuIngredientListResponse{
		MenuVariantID: menuVariantID,
		Ingredients:   ingredients,
	}

	sharedHttp.SendSuccessResponse(w, http.StatusOK, "Menu ingredients retrieved", response)
}
