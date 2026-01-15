package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"invoice-service/pkg/entities/income_invoices/models"
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
	var req models.IncomeInvoiceCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sharedHttp.SendErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	invoice, err := h.dbHandler.Create(&req)
	if err != nil {
		h.logger.WithError(err).Error("Failed to create income invoice")
		sharedHttp.SendErrorResponse(w, http.StatusInternalServerError, "Failed to create income invoice")
		return
	}

	sharedHttp.SendSuccessResponse(w, http.StatusCreated, "Income invoice created successfully", invoice)
}

func (h *HTTPHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	invoice, err := h.dbHandler.GetByID(id)
	if err != nil {
		if err.Error() == "income invoice not found" {
			sharedHttp.SendErrorResponse(w, http.StatusNotFound, "Income invoice not found")
			return
		}
		h.logger.WithError(err).Error("Failed to get income invoice")
		sharedHttp.SendErrorResponse(w, http.StatusInternalServerError, "Failed to get income invoice")
		return
	}

	sharedHttp.SendSuccessResponse(w, http.StatusOK, "Income invoice retrieved successfully", invoice)
}

func (h *HTTPHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req models.IncomeInvoiceUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sharedHttp.SendErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	invoice, err := h.dbHandler.Update(id, &req)
	if err != nil {
		if err.Error() == "income invoice not found" {
			sharedHttp.SendErrorResponse(w, http.StatusNotFound, "Income invoice not found")
			return
		}
		h.logger.WithError(err).Error("Failed to update income invoice")
		sharedHttp.SendErrorResponse(w, http.StatusInternalServerError, "Failed to update income invoice")
		return
	}

	sharedHttp.SendSuccessResponse(w, http.StatusOK, "Income invoice updated successfully", invoice)
}

func (h *HTTPHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := h.dbHandler.Delete(id)
	if err != nil {
		if err.Error() == "income invoice not found" {
			sharedHttp.SendErrorResponse(w, http.StatusNotFound, "Income invoice not found")
			return
		}
		h.logger.WithError(err).Error("Failed to delete income invoice")
		sharedHttp.SendErrorResponse(w, http.StatusInternalServerError, "Failed to delete income invoice")
		return
	}

	sharedHttp.SendSuccessResponse(w, http.StatusOK, "Income invoice deleted successfully", nil)
}

func (h *HTTPHandler) List(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	page := 1
	limit := 10
	customerID := r.URL.Query().Get("customer_id")
	invoiceType := r.URL.Query().Get("invoice_type")
	status := r.URL.Query().Get("status")
	orderID := r.URL.Query().Get("order_id")

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

	req := &models.IncomeInvoiceListRequest{
		Page:  page,
		Limit: limit,
	}

	if customerID != "" {
		req.CustomerID = &customerID
	}
	if invoiceType != "" {
		req.InvoiceType = &invoiceType
	}
	if status != "" {
		req.Status = &status
	}
	if orderID != "" {
		req.OrderID = &orderID
	}

	response, err := h.dbHandler.List(req)
	if err != nil {
		h.logger.WithError(err).Error("Failed to list income invoices")
		sharedHttp.SendErrorResponse(w, http.StatusInternalServerError, "Failed to list income invoices")
		return
	}

	sharedHttp.SendSuccessResponse(w, http.StatusOK, "Income invoices retrieved successfully", response)
}
