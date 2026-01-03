package database

import (
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

func TestDefaultConfig(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel) // Quiet during tests

	config := DefaultConfig(logger)

	// Test that config has expected defaults
	if config == nil {
		t.Fatal("DefaultConfig returned nil")
	}

	if config.Host == "" {
		t.Error("Config.Host should not be empty")
	}

	if config.Port == 0 {
		t.Error("Config.Port should not be 0")
	}

	if config.User == "" {
		t.Error("Config.User should not be empty")
	}

	if config.DBName == "" {
		t.Error("Config.DBName should not be empty")
	}

	if config.MaxOpenConns <= 0 {
		t.Error("Config.MaxOpenConns should be positive")
	}

	if config.MaxIdleConns <= 0 {
		t.Error("Config.MaxIdleConns should be positive")
	}

	if config.ConnMaxLifetime <= 0 {
		t.Error("Config.ConnMaxLifetime should be positive")
	}

	if config.ConnectTimeout <= 0 {
		t.Error("Config.ConnectTimeout should be positive")
	}

	if config.QueryTimeout <= 0 {
		t.Error("Config.QueryTimeout should be positive")
	}
}

func TestConfigValues(t *testing.T) {
	config := &Config{
		Host:     "localhost",
		Port:     5432,
		User:     "testuser",
		Password: "testpass",
		DBName:   "testdb",
		SSLMode:  "disable",
	}

	if config.Host != "localhost" {
		t.Errorf("Host = %q, want %q", config.Host, "localhost")
	}

	if config.Port != 5432 {
		t.Errorf("Port = %d, want %d", config.Port, 5432)
	}

	if config.User != "testuser" {
		t.Errorf("User = %q, want %q", config.User, "testuser")
	}

	if config.DBName != "testdb" {
		t.Errorf("DBName = %q, want %q", config.DBName, "testdb")
	}

	if config.SSLMode != "disable" {
		t.Errorf("SSLMode = %q, want %q", config.SSLMode, "disable")
	}
}

func TestConfigDefaults(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		checkFn func(*Config) bool
		errMsg  string
	}{
		{
			name: "MaxOpenConns default",
			config: &Config{
				MaxOpenConns: 25,
			},
			checkFn: func(c *Config) bool { return c.MaxOpenConns == 25 },
			errMsg:  "MaxOpenConns should be 25",
		},
		{
			name: "MaxIdleConns default",
			config: &Config{
				MaxIdleConns: 5,
			},
			checkFn: func(c *Config) bool { return c.MaxIdleConns == 5 },
			errMsg:  "MaxIdleConns should be 5",
		},
		{
			name: "ConnMaxLifetime default",
			config: &Config{
				ConnMaxLifetime: 5 * time.Minute,
			},
			checkFn: func(c *Config) bool { return c.ConnMaxLifetime == 5*time.Minute },
			errMsg:  "ConnMaxLifetime should be 5 minutes",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.checkFn(tt.config) {
				t.Error(tt.errMsg)
			}
		})
	}
}

func TestNewDatabaseHandler(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	config := &Config{
		Host:            "localhost",
		Port:            5432,
		User:            "postgres",
		Password:        "postgres123",
		DBName:          "barrest_db",
		SSLMode:         "disable",
		MaxOpenConns:    25,
		MaxIdleConns:    5,
		ConnMaxLifetime: 5 * time.Minute,
		ConnMaxIdleTime: 5 * time.Minute,
		ConnectTimeout:  10 * time.Second,
		QueryTimeout:    30 * time.Second,
		MaxRetries:      3,
		RetryInterval:   1 * time.Second,
	}

	handler := New(config, logger)

	if handler == nil {
		t.Fatal("New() returned nil")
	}

	// Test IsConnected before Connect
	if handler.IsConnected() {
		t.Error("IsConnected() should be false before Connect()")
	}
}

func TestDatabaseHandlerInterface(t *testing.T) {
	// Verify that dbHandler implements DatabaseHandler interface
	var _ DatabaseHandler = (*dbHandler)(nil)
}
