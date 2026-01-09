package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"menu-service/pkg/entities/stock_items/models"
	sharedHttp "shared/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// HTTPHandler handles HTTP requests for stock items
type HTTPHandler struct {
	dbHandler *DBHandler
	logger    *logrus.Logger
}

// NewHTTPHandler creates a new HTTP handler
func NewHTTPHandler(dbHandler *DBHandler, logger *logrus.Logger) *HTTPHandler {
	return &HTTPHandler{dbHandler: dbHandler, logger: logger}
}

// List handles GET /api/v1/stock/items
func (h *HTTPHandler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}

	req := &models.StockItemListRequest{
		Page:  page,
		Limit: limit,
	}

	if categoryID := r.URL.Query().Get("category_id"); categoryID != "" {
		req.CategoryID = &categoryID
	}

	response, err := h.dbHandler.List(req)
	if err != nil {
		h.logger.WithError(err).Error("Failed to list stock items")
		sharedHttp.SendError(w, http.StatusInternalServerError, "Failed to list stock items", err)
		return
	}

	sharedHttp.SendSuccess(w, http.StatusOK, "Stock items retrieved", response)
}

// GetByID handles GET /api/v1/stock/items/:id
func (h *HTTPHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	item, err := h.dbHandler.GetByID(id)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get stock item")
		sharedHttp.SendError(w, http.StatusInternalServerError, "Failed to get stock item", err)
		return
	}

	if item == nil {
		sharedHttp.SendError(w, http.StatusNotFound, "Stock item not found", nil)
		return
	}

	sharedHttp.SendSuccess(w, http.StatusOK, "Stock item retrieved", item)
}

// Create handles POST /api/v1/stock/items
func (h *HTTPHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.StockItemCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.WithError(err).Error("Failed to decode request body")
		sharedHttp.SendError(w, http.StatusBadRequest, "Invalid request format", err)
		return
	}

	if req.Name == "" {
		sharedHttp.SendError(w, http.StatusBadRequest, "Name is required", nil)
		return
	}

	if req.Unit == "" {
		sharedHttp.SendError(w, http.StatusBadRequest, "Unit is required", nil)
		return
	}

	item, err := h.dbHandler.Create(&req)
	if err != nil {
		h.logger.WithError(err).Error("Failed to create stock item")
		sharedHttp.SendError(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	sharedHttp.SendSuccess(w, http.StatusCreated, "Stock item created", item)
}

// Update handles PUT /api/v1/stock/items/:id
func (h *HTTPHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req models.StockItemUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.WithError(err).Error("Failed to decode request body")
		sharedHttp.SendError(w, http.StatusBadRequest, "Invalid request format", err)
		return
	}

	item, err := h.dbHandler.Update(id, &req)
	if err != nil {
		h.logger.WithError(err).Error("Failed to update stock item")
		sharedHttp.SendError(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if item == nil {
		sharedHttp.SendError(w, http.StatusNotFound, "Stock item not found", nil)
		return
	}

	sharedHttp.SendSuccess(w, http.StatusOK, "Stock item updated", item)
}

// Delete handles DELETE /api/v1/stock/items/:id
func (h *HTTPHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := h.dbHandler.Delete(id)
	if err != nil {
		h.logger.WithError(err).Error("Failed to delete stock item")
		if err.Error() == "stock item not found" {
			sharedHttp.SendError(w, http.StatusNotFound, "Stock item not found", nil)
			return
		}
		sharedHttp.SendError(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	sharedHttp.SendSuccess(w, http.StatusOK, "Stock item deleted", nil)
}
