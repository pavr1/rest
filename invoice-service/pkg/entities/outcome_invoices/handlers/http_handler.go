package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"invoice-service/pkg/entities/outcome_invoices/models"
	sharedHttp "shared/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type HTTPHandler struct {
	dbHandler *DBHandler
	logger    *logrus.Logger
}

func NewHTTPHandler(dbHandler *DBHandler, logger *logrus.Logger) *HTTPHandler {
	return &HTTPHandler{
		dbHandler: dbHandler,
		logger:    logger,
	}
}

func (h *HTTPHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.OutcomeInvoiceCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sharedHttp.SendErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	invoice, err := h.dbHandler.Create(&req)
	if err != nil {
		h.logger.WithError(err).Error("Failed to create purchase invoice")
		sharedHttp.SendErrorResponse(w, http.StatusInternalServerError, "Failed to create purchase invoice")
		return
	}

	sharedHttp.SendSuccessResponse(w, http.StatusCreated, "Purchase invoice created successfully", invoice)
}

func (h *HTTPHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	invoice, err := h.dbHandler.GetByID(id)
	if err != nil {
		if err.Error() == "purchase invoice not found" {
			sharedHttp.SendErrorResponse(w, http.StatusNotFound, "Purchase invoice not found")
			return
		}
		h.logger.WithError(err).Error("Failed to get purchase invoice")
		sharedHttp.SendErrorResponse(w, http.StatusInternalServerError, "Failed to get purchase invoice")
		return
	}

	sharedHttp.SendSuccessResponse(w, http.StatusOK, "Purchase invoice retrieved successfully", invoice)
}

func (h *HTTPHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req models.OutcomeInvoiceUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sharedHttp.SendErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	invoice, err := h.dbHandler.Update(id, &req)
	if err != nil {
		if err.Error() == "purchase invoice not found" {
			sharedHttp.SendErrorResponse(w, http.StatusNotFound, "Purchase invoice not found")
			return
		}
		h.logger.WithError(err).Error("Failed to update purchase invoice")
		sharedHttp.SendErrorResponse(w, http.StatusInternalServerError, "Failed to update purchase invoice")
		return
	}

	sharedHttp.SendSuccessResponse(w, http.StatusOK, "Purchase invoice updated successfully", invoice)
}

func (h *HTTPHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := h.dbHandler.Delete(id)
	if err != nil {
		if err.Error() == "purchase invoice not found" {
			sharedHttp.SendErrorResponse(w, http.StatusNotFound, "Purchase invoice not found")
			return
		}
		h.logger.WithError(err).Error("Failed to delete purchase invoice")
		sharedHttp.SendErrorResponse(w, http.StatusInternalServerError, "Failed to delete purchase invoice")
		return
	}

	sharedHttp.SendSuccessResponse(w, http.StatusOK, "Purchase invoice deleted successfully", nil)
}

func (h *HTTPHandler) List(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	page := 1
	limit := 10
	supplierID := r.URL.Query().Get("supplier_id")

	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	req := &models.OutcomeInvoiceListRequest{
		Page:  page,
		Limit: limit,
	}

	if supplierID != "" {
		req.SupplierID = &supplierID
	}

	response, err := h.dbHandler.List(req)
	if err != nil {
		h.logger.WithError(err).Error("Failed to list purchase invoices")
		sharedHttp.SendErrorResponse(w, http.StatusInternalServerError, "Failed to list purchase invoices")
		return
	}

	sharedHttp.SendSuccessResponse(w, http.StatusOK, "Purchase invoices retrieved successfully", response)
}
