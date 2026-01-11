package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"invoice-service/pkg/entities/existences/models"
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
	var req models.ExistenceCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sharedHttp.SendErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	existence, err := h.dbHandler.Create(&req)
	if err != nil {
		h.logger.WithError(err).Error("Failed to create existence")
		sharedHttp.SendErrorResponse(w, http.StatusInternalServerError, "Failed to create existence")
		return
	}

	sharedHttp.SendSuccessResponse(w, http.StatusCreated, "Existence created successfully", existence)
}

func (h *HTTPHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	existence, err := h.dbHandler.GetByID(id)
	if err != nil {
		if err.Error() == "existence not found" {
			sharedHttp.SendErrorResponse(w, http.StatusNotFound, "Existence not found")
			return
		}
		h.logger.WithError(err).Error("Failed to get existence")
		sharedHttp.SendErrorResponse(w, http.StatusInternalServerError, "Failed to get existence")
		return
	}

	sharedHttp.SendSuccessResponse(w, http.StatusOK, "Existence retrieved successfully", existence)
}

func (h *HTTPHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req models.ExistenceUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sharedHttp.SendErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	existence, err := h.dbHandler.Update(id, &req)
	if err != nil {
		if err.Error() == "existence not found" {
			sharedHttp.SendErrorResponse(w, http.StatusNotFound, "Existence not found")
			return
		}
		h.logger.WithError(err).Error("Failed to update existence")
		sharedHttp.SendErrorResponse(w, http.StatusInternalServerError, "Failed to update existence")
		return
	}

	sharedHttp.SendSuccessResponse(w, http.StatusOK, "Existence updated successfully", existence)
}

func (h *HTTPHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := h.dbHandler.Delete(id)
	if err != nil {
		if err.Error() == "existence not found" {
			sharedHttp.SendErrorResponse(w, http.StatusNotFound, "Existence not found")
			return
		}
		h.logger.WithError(err).Error("Failed to delete existence")
		sharedHttp.SendErrorResponse(w, http.StatusInternalServerError, "Failed to delete existence")
		return
	}

	sharedHttp.SendSuccessResponse(w, http.StatusOK, "Existence deleted successfully", nil)
}

func (h *HTTPHandler) List(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	page := 1
	limit := 10
	stockItemID := r.URL.Query().Get("stock_item_id")

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

	req := &models.ExistenceListRequest{
		Page:  page,
		Limit: limit,
	}

	if stockItemID != "" {
		req.StockItemID = &stockItemID
	}

	response, err := h.dbHandler.List(req)
	if err != nil {
		h.logger.WithError(err).Error("Failed to list existences")
		sharedHttp.SendErrorResponse(w, http.StatusInternalServerError, "Failed to list existences")
		return
	}

	sharedHttp.SendSuccessResponse(w, http.StatusOK, "Existences retrieved successfully", response)
}