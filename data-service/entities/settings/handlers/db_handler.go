package handlers

import (
	"data-service/pkg/database"
	sharedModels "shared/models"

	"github.com/sirupsen/logrus"
)

// DBHandler handles database operations for settings
type DBHandler struct {
	db     database.DatabaseHandler
	logger *logrus.Logger
}

// NewDBHandler creates a new settings database handler
func NewDBHandler(db database.DatabaseHandler, logger *logrus.Logger) *DBHandler {
	return &DBHandler{
		db:     db,
		logger: logger,
	}
}

// GetSettingsByService retrieves all settings for a specific service
func (h *DBHandler) GetSettingsByService(service string) ([]sharedModels.Setting, error) {
	query := `SELECT setting_id, service, key, value, description, created_at, updated_at 
			  FROM settings WHERE service = $1`

	rows, err := h.db.Query(query, service)
	if err != nil {
		h.logger.WithError(err).Error("Failed to query settings")
		return nil, err
	}
	defer rows.Close()

	var settings []sharedModels.Setting
	for rows.Next() {
		var s sharedModels.Setting
		if err := rows.Scan(&s.SettingID, &s.Service, &s.Key, &s.Value, &s.Description, &s.CreatedAt, &s.UpdatedAt); err != nil {
			h.logger.WithError(err).Error("Failed to scan setting row")
			continue
		}
		settings = append(settings, s)
	}

	if err := rows.Err(); err != nil {
		h.logger.WithError(err).Error("Error iterating settings rows")
		return nil, err
	}

	return settings, nil
}

// GetSettingByKey retrieves a specific setting by service and key
func (h *DBHandler) GetSettingByKey(service, key string) (*sharedModels.Setting, error) {
	query := `SELECT setting_id, service, key, value, description, created_at, updated_at 
			  FROM settings WHERE service = $1 AND key = $2`

	var s sharedModels.Setting
	err := h.db.QueryRow(query, service, key).Scan(
		&s.SettingID, &s.Service, &s.Key, &s.Value, &s.Description, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		h.logger.WithError(err).Error("Failed to query setting by key")
		return nil, err
	}

	return &s, nil
}

// UpdateSetting updates a setting value
func (h *DBHandler) UpdateSetting(service, key, value string) error {
	query := `UPDATE settings SET value = $1, updated_at = NOW() 
			  WHERE service = $2 AND key = $3`

	_, err := h.db.Exec(query, value, service, key)
	if err != nil {
		h.logger.WithError(err).Error("Failed to update setting")
		return err
	}

	return nil
}

// CreateSetting creates a new setting
func (h *DBHandler) CreateSetting(setting sharedModels.Setting) error {
	query := `INSERT INTO settings (service, key, value, description) 
			  VALUES ($1, $2, $3, $4)`

	_, err := h.db.Exec(query, setting.Service, setting.Key, setting.Value, setting.Description)
	if err != nil {
		h.logger.WithError(err).Error("Failed to create setting")
		return err
	}

	return nil
}

// DeleteSetting deletes a setting
func (h *DBHandler) DeleteSetting(service, key string) error {
	query := `DELETE FROM settings WHERE service = $1 AND key = $2`

	_, err := h.db.Exec(query, service, key)
	if err != nil {
		h.logger.WithError(err).Error("Failed to delete setting")
		return err
	}

	return nil
}
