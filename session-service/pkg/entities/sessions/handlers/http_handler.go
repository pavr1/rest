package handlers

import (
	"encoding/json"
	"net/http"
	"session-service/pkg/entities/sessions/models"
	sharedHttp "shared/http"

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
		h.logger.WithError(err).Error("Failed to decode request body")
		sharedHttp.SendErrorResponse(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	if req.Username == "" || req.Password == "" {
		h.logger.Error("Username and password are required")
		sharedHttp.SendErrorResponse(w, http.StatusBadRequest, "Username and password are required")
		return
	}

	response, err := h.dbHandler.CreateSession(&req)
	if err != nil {
		h.logger.WithError(err).Error("Login failed")
		sharedHttp.SendErrorResponse(w, http.StatusUnauthorized, "Invalid username or password")
		return
	}

	h.logger.WithField("username", req.Username).Info("Login successful")
	sharedHttp.SendSuccessResponse(w, http.StatusCreated, "Login successful", response)
}

func (h *HTTPHandler) ValidateSession(w http.ResponseWriter, r *http.Request) {
	var req models.SessionValidationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sharedHttp.SendErrorResponse(w, http.StatusBadRequest, "Invalid request")
		return
	}

	if req.Token == "" {
		sharedHttp.SendErrorResponse(w, http.StatusBadRequest, "Token is required")
		return
	}

	response, err := h.dbHandler.ValidateSession(req.Token)
	if err != nil {
		sharedHttp.SendErrorResponse(w, http.StatusInternalServerError, "Validation failed")
		return
	}

	sharedHttp.SendSuccessResponse(w, http.StatusOK, "Session validated", response)
}

func (h *HTTPHandler) LogoutSession(w http.ResponseWriter, r *http.Request) {
	var req models.SessionLogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sharedHttp.SendErrorResponse(w, http.StatusBadRequest, "Invalid request")
		return
	}

	if req.Token == "" {
		sharedHttp.SendErrorResponse(w, http.StatusBadRequest, "Token required")
		return
	}

	response, err := h.dbHandler.DeleteSession(req.Token)
	if err != nil {
		sharedHttp.SendErrorResponse(w, http.StatusInternalServerError, "Logout failed")
		return
	}

	h.logger.WithField("session_id", response.SessionID).Info("Logout successful")
	sharedHttp.SendSuccessResponse(w, http.StatusOK, "Logged out", response)
}
