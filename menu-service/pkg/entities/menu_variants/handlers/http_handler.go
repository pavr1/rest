package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"menu-service/pkg/entities/menu_variants/models"
	sharedHttp "shared/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// HTTPHandler handles HTTP requests for menu variants
type HTTPHandler struct {
	dbHandler *DBHandler
	logger    *logrus.Logger
}

// NewHTTPHandler creates a new HTTP handler
func NewHTTPHandler(dbHandler *DBHandler, logger *logrus.Logger) *HTTPHandler {
	return &HTTPHandler{dbHandler: dbHandler, logger: logger}
}

// List handles GET /api/v1/menu/items
func (h *HTTPHandler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}

	req := &models.MenuVariantListRequest{
		Page:  page,
		Limit: limit,
	}

	if categoryID := r.URL.Query().Get("category_id"); categoryID != "" {
		req.CategoryID = &categoryID
	}

	if subMenuID := r.URL.Query().Get("sub_category_id"); subMenuID != "" {
		req.SubCategoryID = &subMenuID
	}

	if menuType := r.URL.Query().Get("menu_type"); menuType != "" {
		req.MenuType = &menuType
	}

	if availableStr := r.URL.Query().Get("is_available"); availableStr != "" {
		isAvailable := availableStr == "true"
		req.IsAvailable = &isAvailable
	}

	response, err := h.dbHandler.List(req)
	if err != nil {
		h.logger.WithError(err).Error("Failed to list menu variants")
		sharedHttp.SendErrorResponse(w, http.StatusInternalServerError, "Failed to list menu variants")
		return
	}

	sharedHttp.SendSuccessResponse(w, http.StatusOK, "Menu items retrieved", response)
}

// GetByID handles GET /api/v1/menu/items/:id
func (h *HTTPHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	item, err := h.dbHandler.GetByID(id)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get menu item")
		sharedHttp.SendErrorResponse(w, http.StatusInternalServerError, "Failed to get menu item")
		return
	}

	if item == nil {
		sharedHttp.SendErrorResponse(w, http.StatusNotFound, "Menu item not found")
		return
	}

	sharedHttp.SendSuccessResponse(w, http.StatusOK, "Menu item retrieved", item)
}

// Create handles POST /api/v1/menu/items
func (h *HTTPHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.MenuVariantCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.WithError(err).Error("Failed to decode request body")
		sharedHttp.SendErrorResponse(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	if req.Name == "" {
		sharedHttp.SendErrorResponse(w, http.StatusBadRequest, "Name is required")
		return
	}

	if req.SubCategoryID == "" {
		sharedHttp.SendErrorResponse(w, http.StatusBadRequest, "Sub Category ID is required")
		return
	}

	item, err := h.dbHandler.Create(&req)
	if err != nil {
		h.logger.WithError(err).Error("Failed to create menu item")
		sharedHttp.SendErrorResponse(w, http.StatusInternalServerError, "Failed to create menu item")
		return
	}

	sharedHttp.SendSuccessResponse(w, http.StatusCreated, "Menu item created", item)
}

// Update handles PUT /api/v1/menu/items/:id
func (h *HTTPHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req models.MenuVariantUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.WithError(err).Error("Failed to decode request body")
		sharedHttp.SendErrorResponse(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	item, err := h.dbHandler.Update(id, &req)
	if err != nil {
		h.logger.WithError(err).Error("Failed to update menu item")
		sharedHttp.SendErrorResponse(w, http.StatusInternalServerError, "Failed to update menu item")
		return
	}

	if item == nil {
		sharedHttp.SendErrorResponse(w, http.StatusNotFound, "Menu item not found")
		return
	}

	sharedHttp.SendSuccessResponse(w, http.StatusOK, "Menu item updated", item)
}

// Delete handles DELETE /api/v1/menu/items/:id
func (h *HTTPHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := h.dbHandler.Delete(id)
	if err != nil {
		h.logger.WithError(err).Error("Failed to delete menu item")
		if err.Error() == "menu item not found" {
			sharedHttp.SendErrorResponse(w, http.StatusNotFound, "Menu item not found")
			return
		}
		sharedHttp.SendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	sharedHttp.SendSuccessResponse(w, http.StatusOK, "Menu item deleted", nil)
}

// UpdateAvailability handles PATCH /api/v1/menu/items/:id/availability
func (h *HTTPHandler) UpdateAvailability(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req models.MenuVariantAvailabilityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.WithError(err).Error("Failed to decode request body")
		sharedHttp.SendErrorResponse(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	item, err := h.dbHandler.UpdateAvailability(id, req.IsAvailable)
	if err != nil {
		h.logger.WithError(err).Error("Failed to update menu item availability")
		sharedHttp.SendErrorResponse(w, http.StatusInternalServerError, "Failed to update availability")
		return
	}

	if item == nil {
		sharedHttp.SendErrorResponse(w, http.StatusNotFound, "Menu item not found")
		return
	}

	sharedHttp.SendSuccessResponse(w, http.StatusOK, "Menu item availability updated", item)
}
