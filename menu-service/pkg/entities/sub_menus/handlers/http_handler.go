package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"menu-service/pkg/entities/sub_menus/models"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// HTTPHandler handles HTTP requests for sub menus
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

	req := &models.SubMenuListRequest{
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
		http.Error(w, "Failed to list sub menus", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetByID handles GET /api/v1/menu/submenus/{id}
func (h *HTTPHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	subMenu, err := h.db.GetByID(id)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get sub menu")
		http.Error(w, "Failed to get sub menu", http.StatusInternalServerError)
		return
	}

	if subMenu == nil {
		http.Error(w, "Sub menu not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subMenu)
}

// Create handles POST /api/v1/menu/submenus
func (h *HTTPHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.SubMenuCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.WithError(err).Error("Failed to decode request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}
	if req.CategoryID == "" {
		http.Error(w, "Category ID is required", http.StatusBadRequest)
		return
	}
	if req.ItemType == "" {
		http.Error(w, "Item type is required", http.StatusBadRequest)
		return
	}
	if req.ItemType != "kitchen" && req.ItemType != "bar" {
		http.Error(w, "Item type must be 'kitchen' or 'bar'", http.StatusBadRequest)
		return
	}

	subMenu, err := h.db.Create(&req)
	if err != nil {
		h.logger.WithError(err).Error("Failed to create sub menu")
		http.Error(w, "Failed to create sub menu", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(subMenu)
}

// Update handles PUT /api/v1/menu/submenus/{id}
func (h *HTTPHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req models.SubMenuUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.WithError(err).Error("Failed to decode request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate item_type if provided
	if req.ItemType != nil && *req.ItemType != "kitchen" && *req.ItemType != "bar" {
		http.Error(w, "Item type must be 'kitchen' or 'bar'", http.StatusBadRequest)
		return
	}

	subMenu, err := h.db.Update(id, &req)
	if err != nil {
		h.logger.WithError(err).Error("Failed to update sub menu")
		http.Error(w, "Failed to update sub menu", http.StatusInternalServerError)
		return
	}

	if subMenu == nil {
		http.Error(w, "Sub menu not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subMenu)
}

// Delete handles DELETE /api/v1/menu/submenus/{id}
func (h *HTTPHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := h.db.Delete(id)
	if err != nil {
		h.logger.WithError(err).Error("Failed to delete sub menu")
		if err.Error() == "sub menu not found" {
			http.Error(w, "Sub menu not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
