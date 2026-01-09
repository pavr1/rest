package handlers

import (
	"database/sql"
	"fmt"
	"menu-service/pkg/entities/sub_menus/models"
	subMenuSQL "menu-service/pkg/entities/sub_menus/sql"
	sharedDb "shared/db"

	"github.com/sirupsen/logrus"
)

// DBHandler handles database operations for sub menus
type DBHandler struct {
	db      *sharedDb.DbHandler
	queries *subMenuSQL.Queries
	logger  *logrus.Logger
}

// NewDBHandler creates a new database handler
func NewDBHandler(db *sharedDb.DbHandler, logger *logrus.Logger) (*DBHandler, error) {
	queries, err := subMenuSQL.LoadQueries()
	if err != nil {
		return nil, fmt.Errorf("failed to load SQL queries: %w", err)
	}

	return &DBHandler{
		db:      db,
		queries: queries,
		logger:  logger,
	}, nil
}

// List returns a paginated list of sub menus
func (h *DBHandler) List(req *models.SubMenuListRequest) (*models.SubMenuListResponse, error) {
	offset := (req.Page - 1) * req.Limit

	countQuery, err := h.queries.Get(subMenuSQL.CountSubMenusQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get count query: %w", err)
	}

	var total int
	if err := h.db.QueryRow(countQuery, req.CategoryID, req.ItemType, req.IsActive).Scan(&total); err != nil {
		return nil, fmt.Errorf("failed to count sub menus: %w", err)
	}

	listQuery, err := h.queries.Get(subMenuSQL.ListSubMenusQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get list query: %w", err)
	}

	rows, err := h.db.Query(listQuery, req.CategoryID, req.ItemType, req.IsActive, req.Limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list sub menus: %w", err)
	}
	defer rows.Close()

	var subMenus []models.SubMenu
	for rows.Next() {
		subMenu, err := h.scanSubMenu(rows)
		if err != nil {
			return nil, err
		}
		subMenus = append(subMenus, *subMenu)
	}

	return &models.SubMenuListResponse{
		SubMenus: subMenus,
		Total:    total,
		Page:     req.Page,
		Limit:    req.Limit,
	}, nil
}

// GetByID returns a sub menu by ID
func (h *DBHandler) GetByID(id string) (*models.SubMenu, error) {
	query, err := h.queries.Get(subMenuSQL.GetSubMenuByIDQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get query: %w", err)
	}

	row := h.db.QueryRow(query, id)
	return h.scanSubMenuRow(row)
}

// Create creates a new sub menu
func (h *DBHandler) Create(req *models.SubMenuCreateRequest) (*models.SubMenu, error) {
	query, err := h.queries.Get(subMenuSQL.CreateSubMenuQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get query: %w", err)
	}

	row := h.db.QueryRow(query,
		req.Name, req.Description, req.CategoryID, req.ImageURL,
		req.ItemType, req.DisplayOrder, req.IsActive,
	)

	return h.scanSubMenuRowWithoutCategory(row)
}

// Update updates an existing sub menu
func (h *DBHandler) Update(id string, req *models.SubMenuUpdateRequest) (*models.SubMenu, error) {
	query, err := h.queries.Get(subMenuSQL.UpdateSubMenuQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get query: %w", err)
	}

	row := h.db.QueryRow(query, id,
		req.Name, req.Description, req.CategoryID, req.ImageURL,
		req.ItemType, req.DisplayOrder, req.IsActive,
	)

	return h.scanSubMenuRowWithoutCategory(row)
}

// Delete deletes a sub menu
func (h *DBHandler) Delete(id string) error {
	// Check for dependencies first
	checkQuery, err := h.queries.Get(subMenuSQL.CheckSubMenuDependenciesQuery)
	if err != nil {
		return fmt.Errorf("failed to get check query: %w", err)
	}

	var count int
	if err := h.db.QueryRow(checkQuery, id).Scan(&count); err != nil {
		return fmt.Errorf("failed to check dependencies: %w", err)
	}

	if count > 0 {
		return fmt.Errorf("cannot delete sub menu: it has %d menu items", count)
	}

	deleteQuery, err := h.queries.Get(subMenuSQL.DeleteSubMenuQuery)
	if err != nil {
		return fmt.Errorf("failed to get delete query: %w", err)
	}

	result, err := h.db.Exec(deleteQuery, id)
	if err != nil {
		return fmt.Errorf("failed to delete sub menu: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("sub menu not found")
	}

	h.logger.WithField("id", id).Info("Sub menu deleted")
	return nil
}

// Helper functions for scanning
func (h *DBHandler) scanSubMenu(rows *sql.Rows) (*models.SubMenu, error) {
	var subMenu models.SubMenu
	var description, categoryName, imageURL sql.NullString

	err := rows.Scan(
		&subMenu.ID, &subMenu.Name, &description, &subMenu.CategoryID, &categoryName,
		&imageURL, &subMenu.ItemType, &subMenu.DisplayOrder, &subMenu.IsActive,
		&subMenu.CreatedAt, &subMenu.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to scan sub menu: %w", err)
	}

	if description.Valid {
		subMenu.Description = &description.String
	}
	if categoryName.Valid {
		subMenu.CategoryName = categoryName.String
	}
	if imageURL.Valid {
		subMenu.ImageURL = &imageURL.String
	}

	return &subMenu, nil
}

func (h *DBHandler) scanSubMenuRow(row *sql.Row) (*models.SubMenu, error) {
	var subMenu models.SubMenu
	var description, categoryName, imageURL sql.NullString

	err := row.Scan(
		&subMenu.ID, &subMenu.Name, &description, &subMenu.CategoryID, &categoryName,
		&imageURL, &subMenu.ItemType, &subMenu.DisplayOrder, &subMenu.IsActive,
		&subMenu.CreatedAt, &subMenu.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to scan sub menu: %w", err)
	}

	if description.Valid {
		subMenu.Description = &description.String
	}
	if categoryName.Valid {
		subMenu.CategoryName = categoryName.String
	}
	if imageURL.Valid {
		subMenu.ImageURL = &imageURL.String
	}

	return &subMenu, nil
}

func (h *DBHandler) scanSubMenuRowWithoutCategory(row *sql.Row) (*models.SubMenu, error) {
	var subMenu models.SubMenu
	var description, imageURL sql.NullString

	err := row.Scan(
		&subMenu.ID, &subMenu.Name, &description, &subMenu.CategoryID,
		&imageURL, &subMenu.ItemType, &subMenu.DisplayOrder, &subMenu.IsActive,
		&subMenu.CreatedAt, &subMenu.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to scan sub menu: %w", err)
	}

	if description.Valid {
		subMenu.Description = &description.String
	}
	if imageURL.Valid {
		subMenu.ImageURL = &imageURL.String
	}

	return &subMenu, nil
}
