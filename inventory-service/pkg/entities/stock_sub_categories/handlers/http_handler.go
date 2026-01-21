package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"inventory-service/pkg/entities/stock_sub_categories/models"
	sharedHttp "shared/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// HTTPHandler handles HTTP requests for stock sub-categories
type HTTPHandler struct {
	dbHandler *DBHandler
	logger    *logrus.Logger
}

// NewHTTPHandler creates a new HTTP handler
func NewHTTPHandler(dbHandler *DBHandler, logger *logrus.Logger) *HTTPHandler {
	return &HTTPHandler{dbHandler: dbHandler, logger: logger}
}

// List handles GET /api/v1/stock/sub-categories
func (h *HTTPHandler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 100
	}

	// Check if filtering by category
	categoryID := r.URL.Query().Get("category_id")

	var response *models.StockSubCategoryListResponse
	var err error

	if categoryID != "" {
		response, err = h.dbHandler.ListByCategory(categoryID, page, limit)
	} else {
		response, err = h.dbHandler.List(page, limit)
	}

	if err != nil {
		h.logger.WithError(err).Error("Failed to list stock sub-categories")
		sharedHttp.SendErrorResponse(w, http.StatusInternalServerError, "Failed to list stock sub-categories")
		return
	}

	sharedHttp.SendSuccessResponse(w, http.StatusOK, "Stock sub-categories retrieved", response)
}

// GetByID handles GET /api/v1/stock/sub-categories/:id
func (h *HTTPHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	subCategory, err := h.dbHandler.GetByID(id)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get stock sub-category")
		sharedHttp.SendErrorResponse(w, http.StatusInternalServerError, "Failed to get stock sub-category")
		return
	}

	if subCategory == nil {
		sharedHttp.SendErrorResponse(w, http.StatusNotFound, "Stock sub-category not found")
		return
	}

	sharedHttp.SendSuccessResponse(w, http.StatusOK, "Stock sub-category retrieved", subCategory)
}

// Create handles POST /api/v1/stock/sub-categories
func (h *HTTPHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.StockSubCategoryCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.WithError(err).Error("Failed to decode request body")
		sharedHttp.SendErrorResponse(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	if req.Name == "" {
		sharedHttp.SendErrorResponse(w, http.StatusBadRequest, "Name is required")
		return
	}

	if req.StockCategoryID == "" {
		sharedHttp.SendErrorResponse(w, http.StatusBadRequest, "Stock category ID is required")
		return
	}

	subCategory, err := h.dbHandler.Create(&req)
	if err != nil {
		h.logger.WithError(err).Error("Failed to create stock sub-category")
		sharedHttp.SendErrorResponse(w, http.StatusInternalServerError, "Failed to create stock sub-category")
		return
	}

	sharedHttp.SendSuccessResponse(w, http.StatusCreated, "Stock sub-category created", subCategory)
}

// Update handles PUT /api/v1/stock/sub-categories/:id
func (h *HTTPHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req models.StockSubCategoryUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.WithError(err).Error("Failed to decode request body")
		sharedHttp.SendErrorResponse(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	subCategory, err := h.dbHandler.Update(id, &req)
	if err != nil {
		h.logger.WithError(err).Error("Failed to update stock sub-category")
		sharedHttp.SendErrorResponse(w, http.StatusInternalServerError, "Failed to update stock sub-category")
		return
	}

	if subCategory == nil {
		sharedHttp.SendErrorResponse(w, http.StatusNotFound, "Stock sub-category not found")
		return
	}

	sharedHttp.SendSuccessResponse(w, http.StatusOK, "Stock sub-category updated", subCategory)
}

// Delete handles DELETE /api/v1/stock/sub-categories/:id
func (h *HTTPHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := h.dbHandler.Delete(id)
	if err != nil {
		h.logger.WithError(err).Error("Failed to delete stock sub-category")
		if err.Error() == "stock sub-category not found" {
			sharedHttp.SendErrorResponse(w, http.StatusNotFound, "Stock sub-category not found")
			return
		}
		sharedHttp.SendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	sharedHttp.SendSuccessResponse(w, http.StatusOK, "Stock sub-category deleted", nil)
}
