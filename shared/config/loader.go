package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	sharedModels "shared/models"

	"github.com/sirupsen/logrus"
)

// Configuration constants for bar-restaurant application
const (
	DATA_SERVICE_URL                = "http://barrest_data_service:8086"
	DATA_SERVICE_HOST               = "0.0.0.0"
	DATA_SERVICE_PORT               = 8086
	DATABASE_NAME                   = "postgres"
	DATABASE_PORT                   = 5432
	DATA_SERVICE_USER               = "postgres"
	DATA_SERVICE_PASSWORD           = "postgres123"
	DATA_SERVICE_DB_NAME            = "barrest_db"
	DATA_SERVICE_SSL_MODE           = "disable"
	DATA_SERVICE_MAX_OPEN_CONNS     = 25
	DATA_SERVICE_MAX_IDLE_CONNS     = 5
	DATA_SERVICE_CONN_MAX_LIFETIME  = 5 * time.Minute
	DATA_SERVICE_CONN_MAX_IDLE_TIME = 5 * time.Minute
	DATA_SERVICE_CONNECT_TIMEOUT    = 1 * time.Second
	DATA_SERVICE_QUERY_TIMEOUT      = 30 * time.Second
	DATA_SERVICE_MAX_RETRIES        = 3
	DATA_SERVICE_RETRY_INTERVAL     = 1 * time.Second
)

// ConfigLoader provides functionality to load configuration from the data service
type ConfigLoader struct {
	dataServiceURL string
}

// NewConfigLoader creates a new configuration loader
func NewConfigLoader(dataServiceURL string) *ConfigLoader {
	return &ConfigLoader{
		dataServiceURL: dataServiceURL,
	}
}

// Config is a generic configuration structure that can be used by all services
type Config struct {
	Values map[string]string
	Logger *logrus.Logger
}

// newConfig creates a new config with default values
func newConfig(logger *logrus.Logger) *Config {
	return &Config{
		Values: make(map[string]string),
		Logger: logger,
	}
}

// GetString returns a string value from config
func (c *Config) GetString(key string) string {
	if value, exists := c.Values[key]; exists {
		return value
	}

	c.Logger.WithFields(logrus.Fields{
		"key": key,
	}).Fatal("Key not found in config")
	return ""
}

// GetInt returns an int value from config
func (c *Config) GetInt(key string) int {
	if value, exists := c.Values[key]; exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}

	c.Logger.WithFields(logrus.Fields{
		"key": key,
	}).Fatal("Key not found in config")
	return 0
}

// GetFloat returns a float64 value from config
func (c *Config) GetFloat(key string) float64 {
	if value, exists := c.Values[key]; exists {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
	}

	c.Logger.WithFields(logrus.Fields{
		"key": key,
	}).Fatal("Key not found in config")
	return 0
}

// GetDuration returns a time.Duration value from config
func (c *Config) GetDuration(key string) time.Duration {
	if value, exists := c.Values[key]; exists {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}

	c.Logger.WithFields(logrus.Fields{
		"key": key,
	}).Fatal("Key not found in config")
	return 0
}

// Set sets a key-value pair in the config
func (c *Config) Set(key, value string) {
	c.Values[key] = value
}

// loadSettingsFromDataService calls the data service API to get settings
func (cl *ConfigLoader) loadSettingsFromDataService(serviceName string, logger *logrus.Logger) ([]sharedModels.Setting, error) {
	url := fmt.Sprintf("%s/api/v1/data/settings/by-service", cl.dataServiceURL)

	reqBody := sharedModels.GetSettingsByServiceRequest{
		Service: serviceName,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Gateway-Service", "barrest-gateway")
	req.Header.Set("X-Gateway-Session-Managed", "true")
	req.Header.Set("X-User-ID", "system")
	req.Header.Set("X-Username", "system")
	req.Header.Set("X-User-Role", "admin")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("data service returned status %d", resp.StatusCode)
	}

	var settingsResponse sharedModels.SettingsResponse
	if err := json.NewDecoder(resp.Body).Decode(&settingsResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if settingsResponse.Code != http.StatusOK {
		return nil, fmt.Errorf("data service returned error: %s", settingsResponse.Message)
	}

	jsonData, err := json.Marshal(settingsResponse.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response data: %w", err)
	}

	var settings []sharedModels.Setting
	if err := json.Unmarshal(jsonData, &settings); err != nil {
		return nil, fmt.Errorf("failed to unmarshal settings: %w", err)
	}

	logger.WithFields(logrus.Fields{
		"service":        serviceName,
		"settings_count": len(settings),
	}).Info("Successfully retrieved settings from data service")

	return settings, nil
}

// LoadConfig loads configuration for any service
func (cl *ConfigLoader) LoadConfig(serviceName string, logger *logrus.Logger) (*Config, error) {
	logger.Info("Loading configuration from data service")

	//pvillalobos config service is not ready yet, set default values for now.
	// settings, err := cl.loadSettingsFromDataService(serviceName, logger)
	// if err != nil {
	// 	return nil, err
	// }

	config := newConfig(logger)
	setDefaultValues(config, serviceName)
	//populateConfigFromSettings(config, settings, logger)

	logger.WithFields(logrus.Fields{
		"service": serviceName,
		//"settings_count": len(settings),
	}).Info("Configuration loaded from data service")

	if len(config.Values) == 0 {
		return nil, errors.New("no configuration values loaded")
	}

	return config, nil
}

// setDefaultValues sets default values based on service name
func setDefaultValues(config *Config, serviceName string) {
	switch serviceName {
	case "data":
		config.Set("SERVER_PORT", "8086")
		config.Set("SERVER_HOST", "0.0.0.0")
		config.Set("DB_HOST", "barrest_postgres")
		config.Set("DB_PORT", "5432")
		config.Set("DB_USER", "postgres")
		config.Set("DB_PASSWORD", "postgres123")
		config.Set("DB_NAME", "barrest_db")
		config.Set("DB_SSL_MODE", "disable")
		config.Set("LOG_LEVEL", "info")
	case "session":
		config.Set("SERVER_PORT", "8087")
		config.Set("SERVER_HOST", "0.0.0.0")
		config.Set("DB_HOST", "barrest_postgres")
		config.Set("DB_PORT", "5432")
		config.Set("DB_USER", "postgres")
		config.Set("DB_PASSWORD", "postgres123")
		config.Set("DB_NAME", "barrest_db")
		config.Set("DB_SSL_MODE", "disable")
		config.Set("JWT_SECRET", "barrest-super-secret-jwt-key-change-in-production")
		config.Set("JWT_EXPIRATION_TIME", "30m")
		config.Set("LOG_LEVEL", "info")
	case "orders":
		config.Set("SERVER_PORT", "8083")
		config.Set("SERVER_HOST", "0.0.0.0")
		config.Set("DB_HOST", "localhost")
		config.Set("DB_PORT", "5432")
		config.Set("DB_USER", "postgres")
		config.Set("DB_PASSWORD", "postgres123")
		config.Set("DB_NAME", "barrest_db")
		config.Set("DB_SSL_MODE", "disable")
		config.Set("JWT_SECRET", "barrest-super-secret-jwt-key-change-in-production")
		config.Set("LOG_LEVEL", "info")
		config.Set("DEFAULT_TAX_RATE", "13.0")
		config.Set("DEFAULT_SERVICE_RATE", "10.0")
	case "menu":
		config.Set("SERVER_PORT", "8087")
		config.Set("SERVER_HOST", "0.0.0.0")
		config.Set("DB_HOST", "localhost")
		config.Set("DB_PORT", "5432")
		config.Set("DB_USER", "postgres")
		config.Set("DB_PASSWORD", "postgres123")
		config.Set("DB_NAME", "barrest_db")
		config.Set("DB_SSL_MODE", "disable")
		config.Set("LOG_LEVEL", "info")
	case "inventory":
		config.Set("SERVER_PORT", "8084")
		config.Set("SERVER_HOST", "0.0.0.0")
		config.Set("DB_HOST", "localhost")
		config.Set("DB_PORT", "5432")
		config.Set("DB_USER", "postgres")
		config.Set("DB_PASSWORD", "postgres123")
		config.Set("DB_NAME", "barrest_db")
		config.Set("DB_SSL_MODE", "disable")
		config.Set("LOG_LEVEL", "info")
	case "gateway":
		config.Set("SERVER_PORT", "8082")
		config.Set("SERVER_HOST", "0.0.0.0")
		config.Set("LOG_LEVEL", "INFO")
		config.Set("GATEWAY_SERVICE_URL", "http://localhost:8082")
		config.Set("SESSION_SERVICE_URL", "http://barrest_session_service:8087")
		config.Set("ORDERS_SERVICE_URL", "http://localhost:8083")
		config.Set("MENU_SERVICE_URL", "http://localhost:8087")
		config.Set("INVENTORY_SERVICE_URL", "http://localhost:8084")
		config.Set("PAYMENT_SERVICE_URL", "http://localhost:8088")
		config.Set("CUSTOMER_SERVICE_URL", "http://localhost:8089")
		config.Set("DATA_SERVICE_URL", "http://barrest_data_service:8086")
		config.Set("CORS_ALLOWED_ORIGINS", "*")
		config.Set("CORS_ALLOWED_METHODS", "GET,POST,PUT,DELETE,OPTIONS")
		config.Set("CORS_ALLOWED_HEADERS", "Content-Type,Authorization")
	}
}

// // populateConfigFromSettings populates the config from settings
// func populateConfigFromSettings(config *Config, settings []sharedModels.Setting, logger *logrus.Logger) {
// 	for _, setting := range settings {
// 		config.Set(setting.Key, setting.Value)

// 		logger.WithFields(logrus.Fields{
// 			"key":     setting.Key,
// 			"value":   setting.Value,
// 			"service": setting.Service,
// 		}).Debug("Populated config from data service setting")
// 	}

// 	logger.WithField("settings_processed", len(settings)).Info("Config populated from data service settings")
// }
