package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"inventory-service/pkg/entities/stock_categories/models"
	sharedHttp "shared/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// HTTPHandler handles HTTP requests for stock categories
type HTTPHandler struct {
	dbHandler *DBHandler
	logger    *logrus.Logger
}

// NewHTTPHandler creates a new HTTP handler
func NewHTTPHandler(dbHandler *DBHandler, logger *logrus.Logger) *HTTPHandler {
	return &HTTPHandler{dbHandler: dbHandler, logger: logger}
}

// List handles GET /api/v1/stock/categories
func (h *HTTPHandler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}

	response, err := h.dbHandler.List(page, limit)
	if err != nil {
		h.logger.WithError(err).Error("Failed to list stock categories")
		sharedHttp.SendErrorResponse(w, http.StatusInternalServerError, "Failed to list stock categories")
		return
	}

	sharedHttp.SendSuccessResponse(w, http.StatusOK, "Stock categories retrieved", response)
}

// GetByID handles GET /api/v1/stock/categories/:id
func (h *HTTPHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	category, err := h.dbHandler.GetByID(id)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get stock category")
		sharedHttp.SendErrorResponse(w, http.StatusInternalServerError, "Failed to get stock category")
		return
	}

	if category == nil {
		sharedHttp.SendErrorResponse(w, http.StatusNotFound, "Stock category not found")
		return
	}

	sharedHttp.SendSuccessResponse(w, http.StatusOK, "Stock category retrieved", category)
}

// Create handles POST /api/v1/stock/categories
func (h *HTTPHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.StockCategoryCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.WithError(err).Error("Failed to decode request body")
		sharedHttp.SendErrorResponse(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	if req.Name == "" {
		sharedHttp.SendErrorResponse(w, http.StatusBadRequest, "Name is required")
		return
	}

	category, err := h.dbHandler.Create(&req)
	if err != nil {
		h.logger.WithError(err).Error("Failed to create stock category")
		sharedHttp.SendErrorResponse(w, http.StatusInternalServerError, "Failed to create stock category")
		return
	}

	sharedHttp.SendSuccessResponse(w, http.StatusCreated, "Stock category created", category)
}

// Update handles PUT /api/v1/stock/categories/:id
func (h *HTTPHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req models.StockCategoryUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.WithError(err).Error("Failed to decode request body")
		sharedHttp.SendErrorResponse(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	category, err := h.dbHandler.Update(id, &req)
	if err != nil {
		h.logger.WithError(err).Error("Failed to update stock category")
		sharedHttp.SendErrorResponse(w, http.StatusInternalServerError, "Failed to update stock category")
		return
	}

	if category == nil {
		sharedHttp.SendErrorResponse(w, http.StatusNotFound, "Stock category not found")
		return
	}

	sharedHttp.SendSuccessResponse(w, http.StatusOK, "Stock category updated", category)
}

// Delete handles DELETE /api/v1/stock/categories/:id
func (h *HTTPHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := h.dbHandler.Delete(id)
	if err != nil {
		h.logger.WithError(err).Error("Failed to delete stock category")
		if err.Error() == "stock category not found" {
			sharedHttp.SendErrorResponse(w, http.StatusNotFound, "Stock category not found")
			return
		}
		sharedHttp.SendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	sharedHttp.SendSuccessResponse(w, http.StatusOK, "Stock category deleted", nil)
}
