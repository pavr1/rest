package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"menu-service/pkg/entities/menu_sub_categories/models"
	sharedHttp "shared/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// HTTPHandler handles HTTP requests for menu sub-categories
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

// List handles GET /api/v1/menu/submenus
func (h *HTTPHandler) List(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}

	req := &models.MenuSubCategoryListRequest{
		Page:  page,
		Limit: limit,
	}

	// Parse optional filters
	if categoryID := r.URL.Query().Get("category_id"); categoryID != "" {
		req.CategoryID = &categoryID
	}
	if itemType := r.URL.Query().Get("item_type"); itemType != "" {
		req.ItemType = &itemType
	}
	if isActiveStr := r.URL.Query().Get("is_active"); isActiveStr != "" {
		isActive := isActiveStr == "true"
		req.IsActive = &isActive
	}

	response, err := h.db.List(req)
	if err != nil {
		h.logger.WithError(err).Error("Failed to list sub menus")
		sharedHttp.SendErrorResponse(w, http.StatusInternalServerError, "Failed to list sub menus")
		return
	}

	sharedHttp.SendSuccessResponse(w, http.StatusOK, "Sub menus retrieved", response)
}

// GetByID handles GET /api/v1/menu/submenus/{id}
func (h *HTTPHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	subMenu, err := h.db.GetByID(id)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get sub menu")
		sharedHttp.SendErrorResponse(w, http.StatusInternalServerError, "Failed to get sub menu")
		return
	}

	if subMenu == nil {
		sharedHttp.SendErrorResponse(w, http.StatusNotFound, "Sub menu not found")
		return
	}

	sharedHttp.SendSuccessResponse(w, http.StatusOK, "Sub menu retrieved", subMenu)
}

// Create handles POST /api/v1/menu/submenus
func (h *HTTPHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.MenuSubCategoryCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.WithError(err).Error("Failed to decode request body")
		sharedHttp.SendErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate required fields
	if req.Name == "" {
		sharedHttp.SendErrorResponse(w, http.StatusBadRequest, "Name is required")
		return
	}
	if req.CategoryID == "" {
		sharedHttp.SendErrorResponse(w, http.StatusBadRequest, "Category ID is required")
		return
	}
	if req.ItemType == "" {
		sharedHttp.SendErrorResponse(w, http.StatusBadRequest, "Item type is required")
		return
	}
	if req.ItemType != "kitchen" && req.ItemType != "bar" {
		sharedHttp.SendErrorResponse(w, http.StatusBadRequest, "Item type must be 'kitchen' or 'bar'")
		return
	}

	subMenu, err := h.db.Create(&req)
	if err != nil {
		h.logger.WithError(err).Error("Failed to create sub menu")
		sharedHttp.SendErrorResponse(w, http.StatusInternalServerError, "Failed to create sub menu")
		return
	}

	sharedHttp.SendSuccessResponse(w, http.StatusCreated, "Sub menu created", subMenu)
}

// Update handles PUT /api/v1/menu/submenus/{id}
func (h *HTTPHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req models.MenuSubCategoryUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.WithError(err).Error("Failed to decode request body")
		sharedHttp.SendErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate item_type if provided
	if req.ItemType != nil && *req.ItemType != "kitchen" && *req.ItemType != "bar" {
		sharedHttp.SendErrorResponse(w, http.StatusBadRequest, "Item type must be 'kitchen' or 'bar'")
		return
	}

	subMenu, err := h.db.Update(id, &req)
	if err != nil {
		h.logger.WithError(err).Error("Failed to update sub menu")
		sharedHttp.SendErrorResponse(w, http.StatusInternalServerError, "Failed to update sub menu")
		return
	}

	if subMenu == nil {
		sharedHttp.SendErrorResponse(w, http.StatusNotFound, "Sub menu not found")
		return
	}

	sharedHttp.SendSuccessResponse(w, http.StatusOK, "Sub menu updated", subMenu)
}

// Delete handles DELETE /api/v1/menu/submenus/{id}
func (h *HTTPHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := h.db.Delete(id)
	if err != nil {
		h.logger.WithError(err).Error("Failed to delete sub menu")
		if err.Error() == "sub menu not found" {
			sharedHttp.SendErrorResponse(w, http.StatusNotFound, "Sub menu not found")
			return
		}
		sharedHttp.SendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	sharedHttp.SendSuccessResponse(w, http.StatusOK, "Sub menu deleted", nil)
}
