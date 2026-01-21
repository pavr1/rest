package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"inventory-service/pkg/entities/stock_variants/models"
	sharedHttp "shared/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// HTTPHandler handles HTTP requests for stock variants
type HTTPHandler struct {
	dbHandler *DBHandler
	logger    *logrus.Logger
}

// NewHTTPHandler creates a new HTTP handler
func NewHTTPHandler(dbHandler *DBHandler, logger *logrus.Logger) *HTTPHandler {
	return &HTTPHandler{dbHandler: dbHandler, logger: logger}
}

// List handles GET /api/v1/stock/variants
func (h *HTTPHandler) List(w http.ResponseWriter, r *http.Request) {
	// Check if filtering by sub-category or category
	subCategoryID := r.URL.Query().Get("sub_category_id")
	categoryID := r.URL.Query().Get("category_id")

	var response *models.StockVariantListResponse
	var err error

	if subCategoryID != "" {
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		if page < 1 {
			page = 1
		}
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		if limit < 1 || limit > 100 {
			limit = 100
		}
		response, err = h.dbHandler.ListBySubCategory(subCategoryID, page, limit)
	} else if categoryID != "" {
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		if page < 1 {
			page = 1
		}
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		if limit < 1 || limit > 100 {
			limit = 100
		}
		response, err = h.dbHandler.ListByCategory(categoryID, page, limit)
	} else {
		// No filters - return all active stock variants
		response, err = h.dbHandler.ListAll()
	}

	if err != nil {
		h.logger.WithError(err).Error("Failed to list stock variants")
		sharedHttp.SendErrorResponse(w, http.StatusInternalServerError, "Failed to list stock variants")
		return
	}

	sharedHttp.SendSuccessResponse(w, http.StatusOK, "Stock variants retrieved", response)
}

// GetByID handles GET /api/v1/stock/variants/:id
func (h *HTTPHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	variant, err := h.dbHandler.GetByID(id)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get stock variant")
		sharedHttp.SendErrorResponse(w, http.StatusInternalServerError, "Failed to get stock variant")
		return
	}

	if variant == nil {
		sharedHttp.SendErrorResponse(w, http.StatusNotFound, "Stock variant not found")
		return
	}

	sharedHttp.SendSuccessResponse(w, http.StatusOK, "Stock variant retrieved", variant)
}

// Create handles POST /api/v1/stock/variants
func (h *HTTPHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.StockVariantCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.WithError(err).Error("Failed to decode request body")
		sharedHttp.SendErrorResponse(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	if req.Name == "" {
		sharedHttp.SendErrorResponse(w, http.StatusBadRequest, "Name is required")
		return
	}

	if req.StockSubCategoryID == "" {
		sharedHttp.SendErrorResponse(w, http.StatusBadRequest, "Stock sub-category ID is required")
		return
	}

	variant, err := h.dbHandler.Create(&req)
	if err != nil {
		h.logger.WithError(err).Error("Failed to create stock variant")
		sharedHttp.SendErrorResponse(w, http.StatusInternalServerError, "Failed to create stock variant")
		return
	}

	sharedHttp.SendSuccessResponse(w, http.StatusCreated, "Stock variant created", variant)
}

// Update handles PUT /api/v1/stock/variants/:id
func (h *HTTPHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req models.StockVariantUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.WithError(err).Error("Failed to decode request body")
		sharedHttp.SendErrorResponse(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	variant, err := h.dbHandler.Update(id, &req)
	if err != nil {
		h.logger.WithError(err).Error("Failed to update stock variant")
		sharedHttp.SendErrorResponse(w, http.StatusInternalServerError, "Failed to update stock variant")
		return
	}

	if variant == nil {
		sharedHttp.SendErrorResponse(w, http.StatusNotFound, "Stock variant not found")
		return
	}

	sharedHttp.SendSuccessResponse(w, http.StatusOK, "Stock variant updated", variant)
}

// Delete handles DELETE /api/v1/stock/variants/:id
func (h *HTTPHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := h.dbHandler.Delete(id)
	if err != nil {
		h.logger.WithError(err).Error("Failed to delete stock variant")
		if err.Error() == "stock variant not found" {
			sharedHttp.SendErrorResponse(w, http.StatusNotFound, "Stock variant not found")
			return
		}
		sharedHttp.SendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	sharedHttp.SendSuccessResponse(w, http.StatusOK, "Stock variant deleted", nil)
}
