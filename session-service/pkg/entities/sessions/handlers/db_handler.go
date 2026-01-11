package handlers

import (
	"database/sql"
	"fmt"
	"session-service/pkg/entities/sessions/models"
	sessionSQL "session-service/pkg/entities/sessions/sql"
	sharedConfig "shared/config"
	"time"

	sharedDb "shared/db"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// DBHandler handles database operations for sessions
type DBHandler struct {
	db         *sharedDb.DbHandler
	queries    sessionSQL.Queries
	jwtHandler *JWTHandler
	logger     *logrus.Logger
}

// NewDBHandler creates a new database handler with internal database connection
func NewDBHandler(cfg *sharedConfig.Config, jwtHandler *JWTHandler, logger *logrus.Logger) (*DBHandler, error) {
	// Create database configuration
	dbConfig := &sharedDb.Config{
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
	db, err := sharedDb.NewDatabaseHandler(dbConfig, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create database handler: %w", err)
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

// GetDB returns the underlying database handler for health checks
func (h *DBHandler) GetDB() *sharedDb.DbHandler {
	return h.db
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

	tokenString, err := h.jwtHandler.GenerateToken(staff)
	if err != nil {
		return nil, fmt.Errorf("failed to generate JWT token: %w", err)
	}

	err = h.storeSession(sessionID, tokenString)
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

func (h *DBHandler) storeSession(sessionID, token string) error {
	query, err := h.queries.Get("create_session")
	if err != nil {
		h.logger.WithError(err).Error("Failed to get create session query")
		return fmt.Errorf("failed to get query: %w", err)
	}

	_, err = h.db.Exec(query, sessionID, token)
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

// ValidateSession validates a session token
func (h *DBHandler) ValidateSession(token string) (*models.SessionValidationResponse, error) {
	// First validate the JWT token
	claims, err := h.jwtHandler.ValidateToken(token)
	if err != nil {
		return &models.SessionValidationResponse{
			Valid:   false,
			Message: "Invalid token",
		}, nil
	}

	// Check if token is expired
	if time.Now().After(claims.ExpiresAt.Time) {
		h.deleteSessionByToken(token)
		return &models.SessionValidationResponse{
			Valid:   false,
			Message: "Session expired",
		}, nil
	}

	// Check if token exists in database
	session, err := h.getSessionByToken(token)
	if err != nil {
		if err == sql.ErrNoRows {
			return &models.SessionValidationResponse{
				Valid:   false,
				Message: "Session not found",
			}, nil
		}
		h.logger.WithError(err).Error("Failed to get session by token")
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	// Get staff information from JWT claims
	staff, err := h.getStaffByID(claims.StaffID)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get staff by ID")
		return &models.SessionValidationResponse{
			Valid:   false,
			Message: "User not found",
		}, nil
	}

	// Renew if expiring within 5 minutes
	if time.Until(claims.ExpiresAt.Time) < 5*time.Minute {
		newToken, err := h.jwtHandler.GenerateToken(staff)
		if err == nil {
			h.updateSessionToken(session.SessionID, newToken)
			// Return new token in response (optional, can be handled by client)
		}
	}

	return &models.SessionValidationResponse{
		Valid:     true,
		SessionID: session.SessionID,
		Message:   "Session valid",
		StaffID:   claims.StaffID,
		Username:  claims.Username,
		Role:      claims.Role,
		FullName:  claims.FullName,
	}, nil
}

func (h *DBHandler) getSessionByID(sessionID string) (*models.Session, error) {
	query, err := h.queries.Get(sessionSQL.GetSessionByIDQuery)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get session by ID query")
		return nil, fmt.Errorf("failed to get query: %w", err)
	}

	var session models.Session
	err = h.db.QueryRow(query, sessionID).Scan(
		&session.SessionID, &session.Token,
	)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get session by ID")
		return nil, err
	}

	return &session, nil
}

func (h *DBHandler) getSessionByToken(token string) (*models.Session, error) {
	query, err := h.queries.Get(sessionSQL.GetSessionByTokenQuery)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get session by token query")
		return nil, fmt.Errorf("failed to get query: %w", err)
	}

	var session models.Session
	err = h.db.QueryRow(query, token).Scan(
		&session.SessionID, &session.Token,
	)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get session by token")
		return nil, err
	}

	return &session, nil
}

func (h *DBHandler) getStaffByID(staffID string) (*models.Staff, error) {
	query, err := h.queries.Get("get_staff_by_id")
	if err != nil {
		h.logger.WithError(err).Error("Failed to get staff by ID query")
		return nil, fmt.Errorf("failed to get query: %w", err)
	}

	var staff models.Staff
	var email sql.NullString
	var lastLoginAt sql.NullTime

	err = h.db.QueryRow(query, staffID).Scan(
		&staff.ID, &staff.Username, &email, &staff.FirstName, &staff.LastName, &staff.Role, &staff.IsActive, &lastLoginAt, &staff.CreatedAt, &staff.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if email.Valid {
		staff.Email = &email.String
	}
	if lastLoginAt.Valid {
		staff.LastLoginAt = &lastLoginAt.Time
	}

	return &staff, nil
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

func (h *DBHandler) deleteSessionByToken(token string) error {
	query, err := h.queries.Get(sessionSQL.DeleteSessionByTokenQuery)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get delete session by token query")
		return err
	}
	_, err = h.db.Exec(query, token)
	if err != nil {
		h.logger.WithError(err).Error("Failed to delete session by token")
		return err
	}
	return nil
}

func (h *DBHandler) updateSessionToken(sessionID, token string) error {
	query, err := h.queries.Get(sessionSQL.UpdateSessionTokenQuery)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get update session token query")
		return err
	}
	_, err = h.db.Exec(query, sessionID, token)
	if err != nil {
		h.logger.WithError(err).Error("Failed to update session token")
		return err
	}
	return nil
}

// DeleteSession handles logout by token
func (h *DBHandler) DeleteSession(token string) (*models.SessionLogoutResponse, error) {
	session, err := h.getSessionByToken(token)
	if err != nil {
		if err == sql.ErrNoRows {
			return &models.SessionLogoutResponse{
				Success: false,
				Message: "Session not found",
			}, nil
		}
		h.logger.WithError(err).Error("Failed to get session")
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	if err := h.deleteSessionByToken(token); err != nil {
		return nil, fmt.Errorf("failed to delete session: %w", err)
	}

	return &models.SessionLogoutResponse{
		Success:   true,
		SessionID: session.SessionID,
		Message:   "Logged out successfully",
	}, nil
}
