package http

// import (
// 	"data-service/pkg/entities/settings"
// 	"encoding/json"
// 	"net/http"
// 	sharedHttp "shared/http"
// 	sharedModels "shared/models"

// 	"github.com/sirupsen/logrus"
// )

// // HTTPHandler handles HTTP requests for settings
// type HTTPHandler struct {
// 	repository *settings.Repository
// 	logger     *logrus.Logger
// }

// // NewHTTPHandler creates a new settings HTTP handler
// func NewHTTPHandler(repository *settings.Repository, logger *logrus.Logger) *HTTPHandler {
// 	return &HTTPHandler{
// 		repository: repository,
// 		logger:     logger,
// 	}
// }

// // GetSettingsByService handles POST /api/v1/data/settings/by-service
// func (h *HTTPHandler) GetSettingsByService(w http.ResponseWriter, r *http.Request) {
// 	var req sharedModels.GetSettingsByServiceRequest
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		h.logger.WithError(err).Error("Failed to decode request body")
// 		httpresponse.SendError(w, http.StatusBadRequest, "Invalid request body", err)
// 		return
// 	}

// 	if req.Service == "" {
// 		h.logger.Error("Service name is required")
// 		httpresponse.SendError(w, http.StatusBadRequest, "Service name is required", nil)
// 		return
// 	}

// 	h.logger.WithField("service", req.Service).Info("Loading settings for service")

// 	settings, err := h.repository.GetSettingsByService(req.Service)
// 	if err != nil {
// 		h.logger.WithError(err).Error("Failed to get settings by service")
// 		httpresponse.SendError(w, http.StatusInternalServerError, "Failed to query settings", err)
// 		return
// 	}

// 	h.logger.WithFields(logrus.Fields{
// 		"service": req.Service,
// 		"count":   len(settings),
// 	}).Info("Settings loaded successfully")

// 	httpresponse.SendSuccess(w, http.StatusOK, "Settings retrieved successfully", settings)
// }

// // GetSettingByKey handles POST /api/v1/data/settings/by-key
// func (h *HTTPHandler) GetSettingByKey(w http.ResponseWriter, r *http.Request) {
// 	var req sharedModels.GetSettingsByKeyRequest
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		h.logger.WithError(err).Error("Failed to decode request body")
// 		httpresponse.SendError(w, http.StatusBadRequest, "Invalid request body", err)
// 		return
// 	}

// 	if req.Service == "" || req.Key == "" {
// 		h.logger.Error("Service and key are required")
// 		httpresponse.SendError(w, http.StatusBadRequest, "Service and key are required", nil)
// 		return
// 	}

// 	setting, err := h.repository.GetSettingByKey(req.Service, req.Key)
// 	if err != nil {
// 		h.logger.WithError(err).Error("Failed to get setting by key")
// 		httpresponse.SendError(w, http.StatusNotFound, "Setting not found", err)
// 		return
// 	}

// 	httpresponse.SendSuccess(w, http.StatusOK, "Setting retrieved successfully", setting)
// }

// // UpdateSetting handles PUT /api/v1/data/settings
// func (h *HTTPHandler) UpdateSetting(w http.ResponseWriter, r *http.Request) {
// 	var req sharedModels.UpdateSettingRequest
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		h.logger.WithError(err).Error("Failed to decode request body")
// 		httpresponse.SendError(w, http.StatusBadRequest, "Invalid request body", err)
// 		return
// 	}

// 	if req.Service == "" || req.Key == "" {
// 		h.logger.Error("Service and key are required")
// 		httpresponse.SendError(w, http.StatusBadRequest, "Service and key are required", nil)
// 		return
// 	}

// 	if err := h.repository.UpdateSetting(req.Service, req.Key, req.Value); err != nil {
// 		h.logger.WithError(err).Error("Failed to update setting")
// 		httpresponse.SendError(w, http.StatusInternalServerError, "Failed to update setting", err)
// 		return
// 	}

// 	httpresponse.SendSuccess(w, http.StatusOK, "Setting updated successfully", nil)
// }

// // CreateSetting handles POST /api/v1/data/settings
// func (h *HTTPHandler) CreateSetting(w http.ResponseWriter, r *http.Request) {
// 	var setting sharedModels.Setting
// 	if err := json.NewDecoder(r.Body).Decode(&setting); err != nil {
// 		h.logger.WithError(err).Error("Failed to decode request body")
// 		httpresponse.SendError(w, http.StatusBadRequest, "Invalid request body", err)
// 		return
// 	}

// 	if setting.Service == "" || setting.Key == "" {
// 		h.logger.Error("Service and key are required")
// 		httpresponse.SendError(w, http.StatusBadRequest, "Service and key are required", nil)
// 		return
// 	}

// 	if err := h.repository.CreateSetting(setting); err != nil {
// 		h.logger.WithError(err).Error("Failed to create setting")
// 		httpresponse.SendError(w, http.StatusInternalServerError, "Failed to create setting", err)
// 		return
// 	}

// 	httpresponse.SendSuccess(w, http.StatusCreated, "Setting created successfully", nil)
// }

// // DeleteSetting handles DELETE /api/v1/data/settings
// func (h *HTTPHandler) DeleteSetting(w http.ResponseWriter, r *http.Request) {
// 	var req sharedModels.GetSettingsByKeyRequest
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		h.logger.WithError(err).Error("Failed to decode request body")
// 		httpresponse.SendError(w, http.StatusBadRequest, "Invalid request body", err)
// 		return
// 	}

// 	if req.Service == "" || req.Key == "" {
// 		h.logger.Error("Service and key are required")
// 		httpresponse.SendError(w, http.StatusBadRequest, "Service and key are required", nil)
// 		return
// 	}

// 	if err := h.repository.DeleteSetting(req.Service, req.Key); err != nil {
// 		h.logger.WithError(err).Error("Failed to delete setting")
// 		httpresponse.SendError(w, http.StatusInternalServerError, "Failed to delete setting", err)
// 		return
// 	}

// 	httpresponse.SendSuccess(w, http.StatusOK, "Setting deleted successfully", nil)
// }
