package handlers

// import (
// 	"encoding/json"
// 	"net/http"

// 	"invoice-service/pkg/entities/invoice_items/models"
// 	sharedHttp "shared/http"

// 	"github.com/gorilla/mux"
// 	"github.com/sirupsen/logrus"
// )

// type HTTPHandler struct {
// 	dbHandler *DBHandler
// 	logger    *logrus.Logger
// }

// func NewHTTPHandler(dbHandler *DBHandler, logger *logrus.Logger) *HTTPHandler {
// 	return &HTTPHandler{
// 		dbHandler: dbHandler,
// 		logger:    logger,
// 	}
// }

// func (h *HTTPHandler) Create(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	invoiceID := vars["invoiceId"]

// 	var req models.InvoiceItemCreateRequest
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		sharedHttp.SendErrorResponse(w, http.StatusBadRequest, "Invalid request body")
// 		return
// 	}

// 	item, err := h.dbHandler.Create(invoiceID, &req)
// 	if err != nil {
// 		h.logger.WithError(err).Error("Failed to create invoice item")
// 		sharedHttp.SendErrorResponse(w, http.StatusInternalServerError, "Failed to create invoice item")
// 		return
// 	}

// 	sharedHttp.SendSuccessResponse(w, http.StatusCreated, "Invoice item created successfully", item)
// }

// func (h *HTTPHandler) GetByID(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	id := vars["id"]

// 	item, err := h.dbHandler.GetByID(id)
// 	if err != nil {
// 		if err.Error() == "invoice item not found" {
// 			sharedHttp.SendErrorResponse(w, http.StatusNotFound, "Invoice item not found")
// 			return
// 		}
// 		h.logger.WithError(err).Error("Failed to get invoice item")
// 		sharedHttp.SendErrorResponse(w, http.StatusInternalServerError, "Failed to get invoice item")
// 		return
// 	}

// 	sharedHttp.SendSuccessResponse(w, http.StatusOK, "Invoice item retrieved successfully", item)
// }

// func (h *HTTPHandler) Update(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	id := vars["id"]

// 	var req models.InvoiceItemUpdateRequest
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		sharedHttp.SendErrorResponse(w, http.StatusBadRequest, "Invalid request body")
// 		return
// 	}

// 	item, err := h.dbHandler.Update(id, &req)
// 	if err != nil {
// 		if err.Error() == "invoice item not found" {
// 			sharedHttp.SendErrorResponse(w, http.StatusNotFound, "Invoice item not found")
// 			return
// 		}
// 		h.logger.WithError(err).Error("Failed to update invoice item")
// 		sharedHttp.SendErrorResponse(w, http.StatusInternalServerError, "Failed to update invoice item")
// 		return
// 	}

// 	sharedHttp.SendSuccessResponse(w, http.StatusOK, "Invoice item updated successfully", item)
// }

// func (h *HTTPHandler) Delete(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	id := vars["id"]

// 	err := h.dbHandler.Delete(id)
// 	if err != nil {
// 		if err.Error() == "invoice item not found" {
// 			sharedHttp.SendErrorResponse(w, http.StatusNotFound, "Invoice item not found")
// 			return
// 		}
// 		h.logger.WithError(err).Error("Failed to delete invoice item")
// 		sharedHttp.SendErrorResponse(w, http.StatusInternalServerError, "Failed to delete invoice item")
// 		return
// 	}

// 	sharedHttp.SendSuccessResponse(w, http.StatusOK, "Invoice item deleted successfully", nil)
// }

// func (h *HTTPHandler) ListByInvoice(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	invoiceID := vars["invoiceId"]

// 	response, err := h.dbHandler.ListByInvoice(invoiceID)
// 	if err != nil {
// 		h.logger.WithError(err).Error("Failed to list invoice items")
// 		sharedHttp.SendErrorResponse(w, http.StatusInternalServerError, "Failed to list invoice items")
// 		return
// 	}

// 	sharedHttp.SendSuccessResponse(w, http.StatusOK, "Invoice items retrieved successfully", response)
// }

// // ListByInvoiceOnly gets invoice items for an invoice without sending HTTP response (for internal use)
// func (h *HTTPHandler) ListByInvoiceOnly(invoiceID string) (*models.InvoiceItemListResponse, error) {
// 	return h.dbHandler.ListByInvoice(invoiceID)
// }

// // Create creates an invoice item with specified invoice type
// func (h *HTTPHandler) Create(invoiceID, invoiceType string, req *models.InvoiceItemCreateRequest) (*models.InvoiceItem, error) {
// 	return h.dbHandler.Create(invoiceID, invoiceType, req)
// }
