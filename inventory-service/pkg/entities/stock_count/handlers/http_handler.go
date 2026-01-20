package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"inventory-service/pkg/entities/stock_count/models"
	sharedHttp "shared/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// HTTPHandler handles HTTP requests for stock count
type HTTPHandler struct {
	dbHandler *DBHandler
	logger    *logrus.Logger
}

// NewHTTPHandler creates a new HTTP handler
func NewHTTPHandler(dbHandler *DBHandler, logger *logrus.Logger) *HTTPHandler {
	return &HTTPHandler{dbHandler: dbHandler, logger: logger}
}

// List handles GET /api/v1/inventory/stock-count
func (h *HTTPHandler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}

	// Check if filtering by variant
	variantID := r.URL.Query().Get("stock_variant_id")
	
	var response *models.StockCountListResponse
	var err error
	
	if variantID != "" {
		response, err = h.dbHandler.ListByVariant(variantID, page, limit)
	} else {
		response, err = h.dbHandler.List(page, limit)
	}
	
	if err != nil {
		h.logger.WithError(err).Error("Failed to list stock count records")
		sharedHttp.SendErrorResponse(w, http.StatusInternalServerError, "Failed to list stock count records")
		return
	}

	sharedHttp.SendSuccessResponse(w, http.StatusOK, "Stock count records retrieved", response)
}

// GetByID handles GET /api/v1/inventory/stock-count/:id
func (h *HTTPHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	stockCount, err := h.dbHandler.GetByID(id)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get stock count record")
		sharedHttp.SendErrorResponse(w, http.StatusInternalServerError, "Failed to get stock count record")
		return
	}

	if stockCount == nil {
		sharedHttp.SendErrorResponse(w, http.StatusNotFound, "Stock count record not found")
		return
	}

	sharedHttp.SendSuccessResponse(w, http.StatusOK, "Stock count record retrieved", stockCount)
}

// Create handles POST /api/v1/inventory/stock-count
func (h *HTTPHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.StockCountCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.WithError(err).Error("Failed to decode request body")
		sharedHttp.SendErrorResponse(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	if req.StockVariantID == "" {
		sharedHttp.SendErrorResponse(w, http.StatusBadRequest, "Stock variant ID is required")
		return
	}

	if req.InvoiceID == "" {
		sharedHttp.SendErrorResponse(w, http.StatusBadRequest, "Invoice ID is required")
		return
	}

	if req.Count <= 0 {
		sharedHttp.SendErrorResponse(w, http.StatusBadRequest, "Count must be greater than 0")
		return
	}

	if req.Unit == "" {
		sharedHttp.SendErrorResponse(w, http.StatusBadRequest, "Unit is required")
		return
	}

	if req.PurchasedAt.IsZero() {
		sharedHttp.SendErrorResponse(w, http.StatusBadRequest, "Purchased at is required")
		return
	}

	stockCount, err := h.dbHandler.Create(&req)
	if err != nil {
		h.logger.WithError(err).Error("Failed to create stock count record")
		sharedHttp.SendErrorResponse(w, http.StatusInternalServerError, "Failed to create stock count record")
		return
	}

	sharedHttp.SendSuccessResponse(w, http.StatusCreated, "Stock count record created", stockCount)
}

// Update handles PUT /api/v1/inventory/stock-count/:id
func (h *HTTPHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req models.StockCountUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.WithError(err).Error("Failed to decode request body")
		sharedHttp.SendErrorResponse(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	if req.Count != nil && *req.Count <= 0 {
		sharedHttp.SendErrorResponse(w, http.StatusBadRequest, "Count must be greater than 0")
		return
	}

	stockCount, err := h.dbHandler.Update(id, &req)
	if err != nil {
		h.logger.WithError(err).Error("Failed to update stock count record")
		sharedHttp.SendErrorResponse(w, http.StatusInternalServerError, "Failed to update stock count record")
		return
	}

	if stockCount == nil {
		sharedHttp.SendErrorResponse(w, http.StatusNotFound, "Stock count record not found")
		return
	}

	sharedHttp.SendSuccessResponse(w, http.StatusOK, "Stock count record updated", stockCount)
}

// MarkOut handles PATCH /api/v1/inventory/stock-count/:id/out
func (h *HTTPHandler) MarkOut(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req models.StockCountMarkOutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.WithError(err).Error("Failed to decode request body")
		sharedHttp.SendErrorResponse(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	stockCount, err := h.dbHandler.MarkOut(id, req.IsOut)
	if err != nil {
		h.logger.WithError(err).Error("Failed to mark stock out")
		sharedHttp.SendErrorResponse(w, http.StatusInternalServerError, "Failed to mark stock out")
		return
	}

	if stockCount == nil {
		sharedHttp.SendErrorResponse(w, http.StatusNotFound, "Stock count record not found")
		return
	}

	message := "Stock marked as available"
	if req.IsOut {
		message = "Stock marked as out"
	}
	sharedHttp.SendSuccessResponse(w, http.StatusOK, message, stockCount)
}

// Delete handles DELETE /api/v1/inventory/stock-count/:id
func (h *HTTPHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := h.dbHandler.Delete(id)
	if err != nil {
		h.logger.WithError(err).Error("Failed to delete stock count record")
		if err.Error() == "stock count record not found" {
			sharedHttp.SendErrorResponse(w, http.StatusNotFound, "Stock count record not found")
			return
		}
		sharedHttp.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	sharedHttp.SendSuccessResponse(w, http.StatusOK, "Stock count record deleted", nil)
}
