package handlers

import (
	"database/sql"
	"fmt"
	"session-service/pkg/entities/sessions/models"
	sessionSQL "session-service/pkg/entities/sessions/sql"
	sharedConfig "shared/config"
	"time"

	dataHandlers "data-service/pkg/handlers"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// DBHandler handles database operations for sessions
type DBHandler struct {
	db         dataHandlers.IDBHandler
	queries    sessionSQL.Queries
	jwtHandler *JWTHandler
	logger     *logrus.Logger
}

// NewDBHandler creates a new database handler with internal database connection
func NewDBHandler(cfg *sharedConfig.Config, jwtHandler *JWTHandler, logger *logrus.Logger) (*DBHandler, error) {
	// Create database configuration
	dbConfig := &dataHandlers.Config{
		Host:            cfg.GetString("DB_HOST"),
		Port:            cfg.GetInt("DB_PORT"),
		User:            cfg.GetString("DB_USER"),
		Password:        cfg.GetString("DB_PASSWORD"),
		DBName:          cfg.GetString("DB_NAME"),
		SSLMode:         cfg.GetString("DB_SSL_MODE"),
		MaxOpenConns:    25,
		MaxIdleConns:    5,
		ConnMaxLifetime: 5 * time.Minute,
		ConnMaxIdleTime: 1 * time.Minute,
		ConnectTimeout:  10 * time.Second,
		QueryTimeout:    30 * time.Second,
		MaxRetries:      3,
		RetryInterval:   2 * time.Second,
	}

	// Create database handler using data-service's handler
	db := dataHandlers.New(dbConfig, logger)
	if err := db.Connect(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	queries, err := sessionSQL.LoadQueries()
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to load SQL queries: %w", err)
	}

	return &DBHandler{
		db:         db,
		queries:    *queries,
		jwtHandler: jwtHandler,
		logger:     logger,
	}, nil
}

func (h *DBHandler) Close() error {
	if h.db != nil {
		return h.db.Close()
	}
	return nil
}

// CreateSession creates a new session for a staff member
func (h *DBHandler) CreateSession(req *models.SessionCreateRequest) (*models.SessionCreateResponse, error) {
	staff, err := h.authenticateStaff(req.Username, req.Password)
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	sessionID, err := h.jwtHandler.GenerateSessionID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate session ID: %w", err)
	}

	tokenString, expiresAt, err := h.jwtHandler.GenerateToken(sessionID, staff)
	if err != nil {
		return nil, fmt.Errorf("failed to generate JWT token: %w", err)
	}

	err = h.storeSession(sessionID, tokenString, staff.ID, expiresAt)
	if err != nil {
		return nil, fmt.Errorf("failed to store session: %w", err)
	}

	err = h.updateLastLogin(staff.ID)
	if err != nil {
		h.logger.WithError(err).Warn("Failed to update last login")
	}

	h.logger.WithFields(logrus.Fields{
		"session_id": sessionID,
		"username":   staff.Username,
		"staff_id":   staff.ID,
		"role":       staff.Role,
	}).Info("Session created successfully")

	return &models.SessionCreateResponse{
		SessionID: sessionID,
		Token:     tokenString,
		Message:   "Login successful",
		Staff:     staff,
	}, nil
}

func (h *DBHandler) authenticateStaff(username, password string) (*models.Staff, error) {
	query, err := h.queries.Get("get_staff_by_username")
	if err != nil {
		h.logger.WithError(err).Error("Failed to get staff by username query")
		return nil, fmt.Errorf("failed to get query: %w", err)
	}

	var staff models.Staff
	var passwordHash string
	var email sql.NullString
	var lastLoginAt sql.NullTime

	err = h.db.QueryRow(query, username).Scan(
		&staff.ID, &staff.Username, &email, &passwordHash,
		&staff.FirstName, &staff.LastName, &staff.Role,
		&staff.IsActive, &lastLoginAt, &staff.CreatedAt, &staff.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		h.logger.WithError(err).Error("Failed to get user")
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
		h.logger.WithError(err).Error("Failed to compare password hash and password")
		return nil, fmt.Errorf("invalid password")
	}

	if email.Valid {
		staff.Email = &email.String
	}
	if lastLoginAt.Valid {
		staff.LastLoginAt = &lastLoginAt.Time
	}

	return &staff, nil
}

func (h *DBHandler) storeSession(sessionID, token, staffID string, expiresAt time.Time) error {
	query, err := h.queries.Get("create_session")
	if err != nil {
		h.logger.WithError(err).Error("Failed to get create session query")
		return fmt.Errorf("failed to get query: %w", err)
	}

	_, err = h.db.Exec(query, sessionID, token, staffID, expiresAt)
	if err != nil {
		h.logger.WithError(err).Error("Failed to create session")
		return fmt.Errorf("failed to create session: %w", err)
	}

	return nil
}

func (h *DBHandler) updateLastLogin(staffID string) error {
	query, err := h.queries.Get("update_last_login")
	if err != nil {
		h.logger.WithError(err).Error("Failed to get update last login query")
		return fmt.Errorf("failed to get query: %w", err)
	}

	_, err = h.db.Exec(query, staffID)
	return err
}

// ValidateSession validates a session
func (h *DBHandler) ValidateSession(sessionID string) (*models.SessionValidationResponse, error) {
	session, staff, err := h.getSessionByID(sessionID)
	if err != nil {
		if err == sql.ErrNoRows {
			return &models.SessionValidationResponse{
				Valid:   false,
				Message: "Session not found",
			}, nil
		}
		h.logger.WithError(err).Error("Failed to get session")
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	claims, err := h.jwtHandler.ValidateToken(session.Token)
	if err != nil {
		h.deleteSession(sessionID)
		return &models.SessionValidationResponse{
			Valid:   false,
			Message: "Invalid token",
		}, nil
	}

	if time.Now().After(claims.ExpiresAt.Time) {
		h.deleteSession(sessionID)
		return &models.SessionValidationResponse{
			Valid:   false,
			Message: "Session expired",
		}, nil
	}

	// Renew if expiring within 5 minutes
	if time.Until(claims.ExpiresAt.Time) < 5*time.Minute {
		newToken, newExpiry, err := h.jwtHandler.GenerateToken(sessionID, staff)
		if err == nil {
			h.updateSessionToken(sessionID, newToken, newExpiry)
		}
	}

	return &models.SessionValidationResponse{
		Valid:     true,
		SessionID: sessionID,
		Message:   "Session valid",
		StaffID:   staff.ID,
		Username:  staff.Username,
		Role:      staff.Role,
		FullName:  fmt.Sprintf("%s %s", staff.FirstName, staff.LastName),
	}, nil
}

func (h *DBHandler) getSessionByID(sessionID string) (*models.Session, *models.Staff, error) {
	query, err := h.queries.Get(sessionSQL.GetSessionByIDQuery)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get session by ID query")
		return nil, nil, fmt.Errorf("failed to get query: %w", err)
	}

	var session models.Session
	var staff models.Staff
	var email sql.NullString

	err = h.db.QueryRow(query, sessionID).Scan(
		&session.SessionID, &session.Token, &session.StaffID, &session.CreatedAt, &session.ExpiresAt,
		&staff.ID, &staff.Username, &email, &staff.FirstName, &staff.LastName, &staff.Role, &staff.IsActive,
	)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get session by ID")
		return nil, nil, err
	}

	if email.Valid {
		staff.Email = &email.String
	}

	return &session, &staff, nil
}

func (h *DBHandler) deleteSession(sessionID string) error {
	query, err := h.queries.Get(sessionSQL.DeleteSessionQuery)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get delete session query")
		return err
	}
	_, err = h.db.Exec(query, sessionID)
	if err != nil {
		h.logger.WithError(err).Error("Failed to delete session")
		return err
	}
	return nil
}

func (h *DBHandler) updateSessionToken(sessionID, token string, expiresAt time.Time) error {
	query, err := h.queries.Get(sessionSQL.UpdateSessionTokenQuery)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get update session token query")
		return err
	}
	_, err = h.db.Exec(query, sessionID, token, expiresAt)
	if err != nil {
		h.logger.WithError(err).Error("Failed to update session token")
		return err
	}
	return nil
}

// DeleteSession handles logout
func (h *DBHandler) DeleteSession(sessionID string) (*models.SessionLogoutResponse, error) {
	_, _, err := h.getSessionByID(sessionID)
	if err != nil {
		if err == sql.ErrNoRows {
			return &models.SessionLogoutResponse{
				Success:   false,
				SessionID: sessionID,
				Message:   "Session not found",
			}, nil
		}
		h.logger.WithError(err).Error("Failed to get session")
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	if err := h.deleteSession(sessionID); err != nil {
		return nil, fmt.Errorf("failed to delete session: %w", err)
	}

	return &models.SessionLogoutResponse{
		Success:   true,
		SessionID: sessionID,
		Message:   "Logged out successfully",
	}, nil
}
