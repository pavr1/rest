package handlers

import (
	"encoding/json"
	"net/http"
	"session-service/entities/sessions/models"
	httpresponse "shared/http-response"

	"github.com/sirupsen/logrus"
)

type HTTPHandler struct {
	dbHandler *DBHandler
	logger    *logrus.Logger
}

func NewHTTPHandler(dbHandler *DBHandler, logger *logrus.Logger) *HTTPHandler {
	return &HTTPHandler{dbHandler: dbHandler, logger: logger}
}

func (h *HTTPHandler) CreateSession(w http.ResponseWriter, r *http.Request) {
	var req models.SessionCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpresponse.SendError(w, http.StatusBadRequest, "Invalid request format", err)
		return
	}

	if req.Username == "" || req.Password == "" {
		httpresponse.SendError(w, http.StatusBadRequest, "Username and password are required", nil)
		return
	}

	response, err := h.dbHandler.CreateSession(&req)
	if err != nil {
		h.logger.WithError(err).Error("Login failed")
		httpresponse.SendError(w, http.StatusUnauthorized, "Invalid username or password", nil)
		return
	}

	h.logger.WithField("username", req.Username).Info("Login successful")
	httpresponse.SendSuccess(w, http.StatusCreated, "Login successful", response)
}

func (h *HTTPHandler) ValidateSession(w http.ResponseWriter, r *http.Request) {
	var req models.SessionValidationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpresponse.SendError(w, http.StatusBadRequest, "Invalid request", err)
		return
	}

	response, err := h.dbHandler.ValidateSession(req.SessionID)
	if err != nil {
		httpresponse.SendError(w, http.StatusInternalServerError, "Validation failed", err)
		return
	}

	httpresponse.SendSuccess(w, http.StatusOK, "Session validated", response)
}

func (h *HTTPHandler) LogoutSession(w http.ResponseWriter, r *http.Request) {
	var req models.SessionLogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpresponse.SendError(w, http.StatusBadRequest, "Invalid request", err)
		return
	}

	if req.SessionID == "" {
		httpresponse.SendError(w, http.StatusBadRequest, "Session ID required", nil)
		return
	}

	response, err := h.dbHandler.DeleteSession(req.SessionID)
	if err != nil {
		httpresponse.SendError(w, http.StatusInternalServerError, "Logout failed", err)
		return
	}

	h.logger.WithField("session_id", req.SessionID).Info("Logout successful")
	httpresponse.SendSuccess(w, http.StatusOK, "Logged out", response)
}
