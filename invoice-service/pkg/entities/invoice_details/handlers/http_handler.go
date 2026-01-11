package handlers

import (
	"encoding/json"
	"net/http"

	"invoice-service/pkg/entities/invoice_details/models"
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
	vars := mux.Vars(r)
	invoiceID := vars["invoiceId"]

	var req models.InvoiceDetailCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sharedHttp.SendErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	detail, err := h.dbHandler.Create(invoiceID, &req)
	if err != nil {
		h.logger.WithError(err).Error("Failed to create invoice detail")
		sharedHttp.SendErrorResponse(w, http.StatusInternalServerError, "Failed to create invoice detail")
		return
	}

	sharedHttp.SendSuccessResponse(w, http.StatusCreated, "Invoice detail created successfully", detail)
}

func (h *HTTPHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	invoiceID := vars["invoiceId"]
	id := vars["id"]

	detail, err := h.dbHandler.GetByID(invoiceID, id)
	if err != nil {
		if err.Error() == "invoice detail not found" {
			sharedHttp.SendErrorResponse(w, http.StatusNotFound, "Invoice detail not found")
			return
		}
		h.logger.WithError(err).Error("Failed to get invoice detail")
		sharedHttp.SendErrorResponse(w, http.StatusInternalServerError, "Failed to get invoice detail")
		return
	}

	sharedHttp.SendSuccessResponse(w, http.StatusOK, "Invoice detail retrieved successfully", detail)
}

func (h *HTTPHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	invoiceID := vars["invoiceId"]
	id := vars["id"]

	var req models.InvoiceDetailUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sharedHttp.SendErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	detail, err := h.dbHandler.Update(invoiceID, id, &req)
	if err != nil {
		if err.Error() == "invoice detail not found" {
			sharedHttp.SendErrorResponse(w, http.StatusNotFound, "Invoice detail not found")
			return
		}
		h.logger.WithError(err).Error("Failed to update invoice detail")
		sharedHttp.SendErrorResponse(w, http.StatusInternalServerError, "Failed to update invoice detail")
		return
	}

	sharedHttp.SendSuccessResponse(w, http.StatusOK, "Invoice detail updated successfully", detail)
}

func (h *HTTPHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	invoiceID := vars["invoiceId"]
	id := vars["id"]

	err := h.dbHandler.Delete(invoiceID, id)
	if err != nil {
		if err.Error() == "invoice detail not found" {
			sharedHttp.SendErrorResponse(w, http.StatusNotFound, "Invoice detail not found")
			return
		}
		h.logger.WithError(err).Error("Failed to delete invoice detail")
		sharedHttp.SendErrorResponse(w, http.StatusInternalServerError, "Failed to delete invoice detail")
		return
	}

	sharedHttp.SendSuccessResponse(w, http.StatusOK, "Invoice detail deleted successfully", nil)
}

func (h *HTTPHandler) ListByInvoice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	invoiceID := vars["invoiceId"]

	response, err := h.dbHandler.ListByInvoice(invoiceID)
	if err != nil {
		h.logger.WithError(err).Error("Failed to list invoice details")
		sharedHttp.SendErrorResponse(w, http.StatusInternalServerError, "Failed to list invoice details")
		return
	}

	sharedHttp.SendSuccessResponse(w, http.StatusOK, "Invoice details retrieved successfully", response)
}