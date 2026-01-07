package db

import (
	"context"
	"database/sql"
	"fmt"
	"shared/config"
	"time"

	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

const (
	DBHealthCheckInterval = 1 * time.Second
)

// Config holds database configuration
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string

	// Connection pool settings
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration

	// Timeout settings
	ConnectTimeout time.Duration
	QueryTimeout   time.Duration

	// Retry settings
	MaxRetries    int
	RetryInterval time.Duration
}

// DefaultConfig returns a default configuration
func DefaultConfig(logger *logrus.Logger) *Config {
	host := config.DATABASE_NAME
	port := config.DATABASE_PORT
	user := config.DATA_SERVICE_USER
	password := config.DATA_SERVICE_PASSWORD
	dbName := config.DATA_SERVICE_DB_NAME
	sslMode := config.DATA_SERVICE_SSL_MODE

	dbConfig := &Config{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		DBName:   dbName,
		SSLMode:  sslMode,

		// Connection pool settings
		MaxOpenConns:    config.DATA_SERVICE_MAX_OPEN_CONNS,
		MaxIdleConns:    config.DATA_SERVICE_MAX_IDLE_CONNS,
		ConnMaxLifetime: config.DATA_SERVICE_CONN_MAX_LIFETIME,
		ConnMaxIdleTime: config.DATA_SERVICE_CONN_MAX_IDLE_TIME,

		// Timeout settings
		ConnectTimeout: config.DATA_SERVICE_CONNECT_TIMEOUT,
		QueryTimeout:   config.DATA_SERVICE_QUERY_TIMEOUT,

		// Retry settings
		MaxRetries:    config.DATA_SERVICE_MAX_RETRIES,
		RetryInterval: config.DATA_SERVICE_RETRY_INTERVAL,
	}

	return dbConfig
}

// DbHandler implements the IDBHandler interface
type DbHandler struct {
	ctx           context.Context
	cancelCtx     context.CancelFunc
	db            *sql.DB
	config        *Config
	logger        *logrus.Logger
	healthMonitor *DBHealthMonitor
}

// NewDbHandler creates a new database handler instance
func NewDbHandler(config *Config, logger *logrus.Logger) *DbHandler {
	return &DbHandler{
		config: config,
		logger: logger,
	}
}

// New creates a new database handler instance
func NewDatabaseHandler(config *Config, logger *logrus.Logger) (*DbHandler, error) {
	var err error
	db := NewDbHandler(config, logger)

	// Connect to database
	fmt.Println("🍺 Connecting to Bar-Restaurant Data Service...")
	if err := db.connect(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Perform initial health check
	if err := db.Ping(); err != nil {
		logger.WithError(err).Fatal("Initial database health check failed")
		return nil, fmt.Errorf("initial database health check failed")
	}

	fmt.Println("✅ Database connection established successfully, starting health monitor...")

	// Create cancellable context for health monitor
	db.ctx, db.cancelCtx = context.WithCancel(context.Background())
	db.healthMonitor, err = NewHealthMonitor(logger, DBHealthCheckInterval, db)
	if err != nil {
		return nil, fmt.Errorf("failed to create health monitor: %w", err)
	}
	db.healthMonitor.Start(db.ctx)

	return db, nil
}

// Connect establishes a connection to the database
func (h *DbHandler) connect() error {
	h.logger.WithFields(logrus.Fields{
		"host":   h.config.Host,
		"port":   h.config.Port,
		"dbname": h.config.DBName,
		"user":   h.config.User,
	}).Info("Attempting to connect to database")

	connStr := h.buildConnectionString()

	var err error
	var db *sql.DB

	// Retry connection with exponential backoff
	for attempt := 1; attempt <= h.config.MaxRetries; attempt++ {
		db, err = sql.Open("postgres", connStr)
		if err != nil {
			h.logger.WithFields(logrus.Fields{
				"attempt": attempt,
				"error":   err.Error(),
			}).Warn("Failed to open database connection")

			if attempt < h.config.MaxRetries {
				time.Sleep(h.config.RetryInterval * time.Duration(attempt))
				continue
			}
			return fmt.Errorf("failed to open database after %d attempts: %w", h.config.MaxRetries, err)
		}

		// Test the connection
		ctx, cancel := context.WithTimeout(context.Background(), h.config.ConnectTimeout)
		err = db.PingContext(ctx)
		cancel()

		if err != nil {
			h.logger.WithError(err).WithFields(logrus.Fields{
				"attempt": attempt,
				"host":    h.config.Host,
				"port":    h.config.Port,
				"dbname":  h.config.DBName,
			}).Warn("Failed to ping database")

			db.Close()
			if attempt < h.config.MaxRetries {
				time.Sleep(h.config.RetryInterval * time.Duration(attempt))
				continue
			}
			return fmt.Errorf("failed to ping database after %d attempts: %w", h.config.MaxRetries, err)
		}

		break
	}

	// Configure connection pool
	h.configureConnectionPool(db)

	h.db = db

	h.logger.WithFields(logrus.Fields{
		"host":   h.config.Host,
		"port":   h.config.Port,
		"dbname": h.config.DBName,
	}).Info("Successfully connected to database")

	return nil
}

// Close closes the database connection
func (h *DbHandler) Close() error {
	if h.db == nil {
		return nil
	}

	h.logger.Info("Closing database connection")

	// Cancel context to stop health monitor
	if h.cancelCtx != nil {
		h.cancelCtx()
	}

	err := h.db.Close()
	if err != nil {
		h.logger.WithError(err).Error("Failed to close database connection")
		return err
	}

	h.logger.Info("Database connection closed successfully")
	return nil
}

// Ping tests the database connection
func (h *DbHandler) Ping() error {
	if h.db == nil {
		h.logger.Error("Database connection is nil for Ping")
		return fmt.Errorf("database connection is nil")
	}

	// Use 2 second timeout for health checks (shorter than default)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Execute a simple query instead of PingContext
	// This forces a real statement execution and respects statement_timeout
	var result int
	err := h.db.QueryRowContext(ctx, "SELECT 1").Scan(&result)
	if err != nil {
		h.logger.WithError(err).Error("Database ping failed")
		return err
	}

	return nil
}

// BeginTx starts a new transaction
func (h *DbHandler) BeginTx(ctx context.Context) (*sql.Tx, error) {
	if h.db == nil {
		h.logger.Error("Database connection is nil for BeginTx")
		return nil, fmt.Errorf("database connection is nil")
	}

	tx, err := h.db.BeginTx(ctx, nil)
	if err != nil {
		h.logger.WithError(err).Error("Failed to begin transaction for BeginTx")
		return nil, err
	}

	h.logger.Debug("Transaction started")
	return tx, nil
}

// CommitTx commits a transaction
func (h *DbHandler) CommitTx(tx *sql.Tx) error {
	if tx == nil {
		h.logger.Error("Transaction is nil for CommitTx")
		return fmt.Errorf("transaction is nil")
	}

	err := tx.Commit()
	if err != nil {
		h.logger.WithError(err).Error("Failed to commit transaction for CommitTx")
		return err
	}

	h.logger.Debug("Transaction committed")
	return nil
}

// RollbackTx rolls back a transaction
func (h *DbHandler) RollbackTx(tx *sql.Tx) error {
	if tx == nil {
		h.logger.Error("Transaction is nil for RollbackTx")
		return fmt.Errorf("transaction is nil")
	}

	err := tx.Rollback()
	if err != nil {
		h.logger.WithError(err).Error("Failed to rollback transaction for RollbackTx")
		return err
	}

	h.logger.Debug("Transaction rolled back")
	return nil
}

// Query executes a query with logging
func (h *DbHandler) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return h.QueryContext(context.Background(), query, args...)
}

// QueryContext executes a query with context and logging
func (h *DbHandler) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	if h.db == nil {
		h.logger.Error("Database connection is nil for QueryContext")
		return nil, fmt.Errorf("database connection is nil")
	}

	start := time.Now()
	rows, err := h.db.QueryContext(ctx, query, args...)
	duration := time.Since(start)

	logEntry := h.logger.WithFields(logrus.Fields{
		"query":      h.sanitizeQuery(query),
		"duration":   duration,
		"args_count": len(args),
	})

	if err != nil {
		logEntry.WithError(err).Error("Query execution failed for QueryContext")
		return nil, h.handlePostgreSQLError(err)
	}

	logEntry.Debug("Query executed successfully")
	return rows, nil
}

// QueryRow executes a query that returns a single row
func (h *DbHandler) QueryRow(query string, args ...interface{}) *sql.Row {
	return h.QueryRowContext(context.Background(), query, args...)
}

// QueryRowContext executes a query that returns a single row with context
func (h *DbHandler) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	if h.db == nil {
		h.logger.Error("Database connection is nil for QueryRow")
		return nil
	}

	start := time.Now()
	row := h.db.QueryRowContext(ctx, query, args...)
	duration := time.Since(start)

	h.logger.WithFields(logrus.Fields{
		"query":      h.sanitizeQuery(query),
		"duration":   duration,
		"args_count": len(args),
	}).Debug("QueryRow executed")

	return row
}

// Exec executes a query without returning rows
func (h *DbHandler) Exec(query string, args ...interface{}) (sql.Result, error) {
	return h.ExecContext(context.Background(), query, args...)
}

// ExecContext executes a query without returning rows with context
func (h *DbHandler) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	if h.db == nil {
		h.logger.Error("Database connection is nil for ExecContext")
		return nil, fmt.Errorf("database connection is nil")
	}

	start := time.Now()
	result, err := h.db.ExecContext(ctx, query, args...)
	duration := time.Since(start)

	logEntry := h.logger.WithFields(logrus.Fields{
		"query":      h.sanitizeQuery(query),
		"duration":   duration,
		"args_count": len(args),
	})

	if err != nil {
		logEntry.WithError(err).Error("Exec execution failed for ExecContext")
		return nil, h.handlePostgreSQLError(err)
	}

	logEntry.Debug("Exec executed successfully")
	return result, nil
}

// Prepare creates a prepared statement
func (h *DbHandler) Prepare(query string) (*sql.Stmt, error) {
	return h.PrepareContext(context.Background(), query)
}

// PrepareContext creates a prepared statement with context
func (h *DbHandler) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	if h.db == nil {
		h.logger.Error("Database connection is nil for PrepareContext")
		return nil, fmt.Errorf("database connection is nil")
	}

	stmt, err := h.db.PrepareContext(ctx, query)
	if err != nil {
		h.logger.WithFields(logrus.Fields{
			"query": h.sanitizeQuery(query),
		}).WithError(err).Error("Failed to prepare statement for PrepareContext")
		return nil, h.handlePostgreSQLError(err)
	}

	h.logger.WithFields(logrus.Fields{
		"query": h.sanitizeQuery(query),
	}).Debug("Statement prepared successfully")

	return stmt, nil
}

// GetDB returns the underlying sql.DB instance
func (h *DbHandler) GetDB() *sql.DB {
	return h.db
}

// GetStats returns database connection statistics
func (h *DbHandler) GetStats() sql.DBStats {
	if h.db == nil {
		return sql.DBStats{}
	}
	return h.db.Stats()
}

// IsConnected returns the connection status
func (h *DbHandler) IsConnected() bool {
	return h.healthMonitor.IsHealthy()
}

// buildConnectionString creates the PostgreSQL connection string
func (h *DbHandler) buildConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s connect_timeout=%d statement_timeout=%d tcp_user_timeout=%d",
		h.config.Host,
		h.config.Port,
		h.config.User,
		h.config.Password,
		h.config.DBName,
		h.config.SSLMode,
		int(h.config.ConnectTimeout.Seconds()),
		1000, // 1 second statement timeout in milliseconds
		2000, // 2 second TCP user timeout in milliseconds (forces TCP to give up faster)
	)
}

// configureConnectionPool sets up the connection pool
func (h *DbHandler) configureConnectionPool(db *sql.DB) {
	db.SetMaxOpenConns(h.config.MaxOpenConns)
	db.SetMaxIdleConns(h.config.MaxIdleConns)
	db.SetConnMaxLifetime(h.config.ConnMaxLifetime)
	db.SetConnMaxIdleTime(h.config.ConnMaxIdleTime)

	h.logger.WithFields(logrus.Fields{
		"max_open_conns":     h.config.MaxOpenConns,
		"max_idle_conns":     h.config.MaxIdleConns,
		"conn_max_lifetime":  h.config.ConnMaxLifetime,
		"conn_max_idle_time": h.config.ConnMaxIdleTime,
	}).Info("Database connection pool configured")
}

// sanitizeQuery removes sensitive information from queries for logging
func (h *DbHandler) sanitizeQuery(query string) string {
	if len(query) > 100 {
		return query[:100] + "..."
	}
	return query
}

// handlePostgreSQLError handles PostgreSQL-specific errors
func (h *DbHandler) handlePostgreSQLError(err error) error {
	if pqErr, ok := err.(*pq.Error); ok {
		h.logger.WithFields(logrus.Fields{
			"code":       pqErr.Code,
			"constraint": pqErr.Constraint,
			"detail":     pqErr.Detail,
			"hint":       pqErr.Hint,
			"position":   pqErr.Position,
			"table":      pqErr.Table,
			"column":     pqErr.Column,
		}).Error("PostgreSQL error occurred")

		switch pqErr.Code {
		case "23505": // unique_violation
			return fmt.Errorf("duplicate entry: %s", pqErr.Detail)
		case "23503": // foreign_key_violation
			return fmt.Errorf("foreign key constraint violation: %s", pqErr.Detail)
		case "23502": // not_null_violation
			return fmt.Errorf("required field missing: %s", pqErr.Column)
		default:
			return fmt.Errorf("database error [%s]: %s", pqErr.Code, pqErr.Message)
		}
	}

	return err
}

// IDBHandler defines the interface for database operations
type IDBHandler interface {
	// Connection management
	connect() error
	Close() error
	Ping() error

	// Transaction management
	BeginTx(ctx context.Context) (*sql.Tx, error)
	CommitTx(tx *sql.Tx) error
	RollbackTx(tx *sql.Tx) error

	// Query operations
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row

	// Execute operations
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)

	// Prepared statements
	Prepare(query string) (*sql.Stmt, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)

	// Utility methods
	GetDB() *sql.DB
	GetStats() sql.DBStats
	IsConnected() bool
}
