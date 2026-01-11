package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"inventory-service/pkg/entities/stock_item_categories/models"
	sharedHttp "shared/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// HTTPHandler handles HTTP requests for stock item categories
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
		h.logger.WithError(err).Error("Failed to list stock item categories")
		sharedHttp.SendError(w, http.StatusInternalServerError, "Failed to list stock item categories", err)
		return
	}

	sharedHttp.SendSuccess(w, http.StatusOK, "Stock item categories retrieved", response)
}

// GetByID handles GET /api/v1/stock/categories/:id
func (h *HTTPHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	category, err := h.dbHandler.GetByID(id)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get stock item category")
		sharedHttp.SendError(w, http.StatusInternalServerError, "Failed to get stock item category", err)
		return
	}

	if category == nil {
		sharedHttp.SendError(w, http.StatusNotFound, "Stock item category not found", nil)
		return
	}

	sharedHttp.SendSuccess(w, http.StatusOK, "Stock item category retrieved", category)
}

// Create handles POST /api/v1/stock/categories
func (h *HTTPHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.StockItemCategoryCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.WithError(err).Error("Failed to decode request body")
		sharedHttp.SendError(w, http.StatusBadRequest, "Invalid request format", err)
		return
	}

	if req.Name == "" {
		sharedHttp.SendError(w, http.StatusBadRequest, "Name is required", nil)
		return
	}

	category, err := h.dbHandler.Create(&req)
	if err != nil {
		h.logger.WithError(err).Error("Failed to create stock item category")
		sharedHttp.SendError(w, http.StatusInternalServerError, "Failed to create stock item category", err)
		return
	}

	sharedHttp.SendSuccess(w, http.StatusCreated, "Stock item category created", category)
}

// Update handles PUT /api/v1/stock/categories/:id
func (h *HTTPHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req models.StockItemCategoryUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.WithError(err).Error("Failed to decode request body")
		sharedHttp.SendError(w, http.StatusBadRequest, "Invalid request format", err)
		return
	}

	category, err := h.dbHandler.Update(id, &req)
	if err != nil {
		h.logger.WithError(err).Error("Failed to update stock item category")
		sharedHttp.SendError(w, http.StatusInternalServerError, "Failed to update stock item category", err)
		return
	}

	if category == nil {
		sharedHttp.SendError(w, http.StatusNotFound, "Stock item category not found", nil)
		return
	}

	sharedHttp.SendSuccess(w, http.StatusOK, "Stock item category updated", category)
}

// Delete handles DELETE /api/v1/stock/categories/:id
func (h *HTTPHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := h.dbHandler.Delete(id)
	if err != nil {
		h.logger.WithError(err).Error("Failed to delete stock item category")
		if err.Error() == "stock item category not found" {
			sharedHttp.SendError(w, http.StatusNotFound, "Stock item category not found", nil)
			return
		}
		sharedHttp.SendError(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	sharedHttp.SendSuccess(w, http.StatusOK, "Stock item category deleted", nil)
}
