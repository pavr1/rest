package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"inventory-service/pkg/entities/suppliers/models"
	sharedHttp "shared/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// HTTPHandler handles HTTP requests for suppliers
type HTTPHandler struct {
	dbHandler *DBHandler
	logger    *logrus.Logger
}

// NewHTTPHandler creates a new HTTP handler
func NewHTTPHandler(dbHandler *DBHandler, logger *logrus.Logger) *HTTPHandler {
	return &HTTPHandler{dbHandler: dbHandler, logger: logger}
}

// List handles GET /api/v1/inventory/suppliers
func (h *HTTPHandler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}

	req := &models.SupplierListRequest{
		Page:  page,
		Limit: limit,
	}

	if name := r.URL.Query().Get("name"); name != "" {
		req.Name = &name
	}
	if email := r.URL.Query().Get("email"); email != "" {
		req.Email = &email
	}
	if phone := r.URL.Query().Get("phone"); phone != "" {
		req.Phone = &phone
	}

	response, err := h.dbHandler.List(req)
	if err != nil {
		h.logger.WithError(err).Error("Failed to list suppliers")
		sharedHttp.SendErrorResponse(w, http.StatusInternalServerError, "Failed to list suppliers")
		return
	}

	sharedHttp.SendSuccessResponse(w, http.StatusOK, "Suppliers retrieved successfully", response)
}

// GetByID handles GET /api/v1/inventory/suppliers/{id}
func (h *HTTPHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		sharedHttp.SendErrorResponse(w, http.StatusBadRequest, "Supplier ID is required")
		return
	}

	supplier, err := h.dbHandler.GetByID(id)
	if err != nil {
		if err.Error() == "supplier not found" {
			sharedHttp.SendErrorResponse(w, http.StatusNotFound, "Supplier not found")
			return
		}
		h.logger.WithError(err).Error("Failed to get supplier")
		sharedHttp.SendErrorResponse(w, http.StatusInternalServerError, "Failed to get supplier")
		return
	}

	sharedHttp.SendSuccessResponse(w, http.StatusOK, "Supplier retrieved successfully", supplier)
}

// Create handles POST /api/v1/inventory/suppliers
func (h *HTTPHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.SupplierCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sharedHttp.SendErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Name == "" {
		sharedHttp.SendErrorResponse(w, http.StatusBadRequest, "Supplier name is required")
		return
	}

	supplier, err := h.dbHandler.Create(&req)
	if err != nil {
		h.logger.WithError(err).Error("Failed to create supplier")
		sharedHttp.SendErrorResponse(w, http.StatusInternalServerError, "Failed to create supplier")
		return
	}

	sharedHttp.SendSuccessResponse(w, http.StatusCreated, "Supplier created successfully", supplier)
}

// Update handles PUT /api/v1/inventory/suppliers/{id}
func (h *HTTPHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		sharedHttp.SendErrorResponse(w, http.StatusBadRequest, "Supplier ID is required")
		return
	}

	var req models.SupplierUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sharedHttp.SendErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	supplier, err := h.dbHandler.Update(id, &req)
	if err != nil {
		if err.Error() == "supplier not found" {
			sharedHttp.SendErrorResponse(w, http.StatusNotFound, "Supplier not found")
			return
		}
		h.logger.WithError(err).Error("Failed to update supplier")
		sharedHttp.SendErrorResponse(w, http.StatusInternalServerError, "Failed to update supplier")
		return
	}

	sharedHttp.SendSuccessResponse(w, http.StatusOK, "Supplier updated successfully", supplier)
}

// Delete handles DELETE /api/v1/inventory/suppliers/{id}
func (h *HTTPHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		sharedHttp.SendErrorResponse(w, http.StatusBadRequest, "Supplier ID is required")
		return
	}

	err := h.dbHandler.Delete(id)
	if err != nil {
		if err.Error() == "supplier not found" {
			sharedHttp.SendErrorResponse(w, http.StatusNotFound, "Supplier not found")
			return
		}
		h.logger.WithError(err).Error("Failed to delete supplier")
		sharedHttp.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	sharedHttp.SendSuccessResponse(w, http.StatusOK, "Supplier deleted successfully", nil)
}