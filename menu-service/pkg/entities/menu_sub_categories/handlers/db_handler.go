package handlers

import (
	"database/sql"
	"fmt"
	"menu-service/pkg/entities/menu_sub_categories/models"
	menuSubCategorySQL "menu-service/pkg/entities/menu_sub_categories/sql"
	sharedDb "shared/db"

	"github.com/sirupsen/logrus"
)

// DBHandler handles database operations for menu sub-categories
type DBHandler struct {
	db      *sharedDb.DbHandler
	queries *menuSubCategorySQL.Queries
	logger  *logrus.Logger
}

// NewDBHandler creates a new database handler
func NewDBHandler(db *sharedDb.DbHandler, logger *logrus.Logger) (*DBHandler, error) {
	queries, err := menuSubCategorySQL.LoadQueries()
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
func (h *DBHandler) List(req *models.MenuSubCategoryListRequest) (*models.MenuSubCategoryListResponse, error) {
	offset := (req.Page - 1) * req.Limit

	countQuery, err := h.queries.Get(menuSubCategorySQL.CountMenuSubCategoriesQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get count query: %w", err)
	}

	var total int
	if err := h.db.QueryRow(countQuery, req.CategoryID, req.ItemType, req.IsActive).Scan(&total); err != nil {
		return nil, fmt.Errorf("failed to count sub menus: %w", err)
	}

	listQuery, err := h.queries.Get(menuSubCategorySQL.ListMenuSubCategoriesQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get list query: %w", err)
	}

	rows, err := h.db.Query(listQuery, req.CategoryID, req.ItemType, req.IsActive, req.Limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list sub menus: %w", err)
	}
	defer rows.Close()

	var subMenus []models.MenuSubCategory
	for rows.Next() {
		subMenu, err := h.scanMenuSubCategory(rows)
		if err != nil {
			return nil, err
		}
		subMenus = append(subMenus, *subMenu)
	}

	return &models.MenuSubCategoryListResponse{
		SubCategories: subMenus,
		Total:         total,
		Page:          req.Page,
		Limit:         req.Limit,
	}, nil
}

// GetByID returns a sub menu by ID
func (h *DBHandler) GetByID(id string) (*models.MenuSubCategory, error) {
	query, err := h.queries.Get(menuSubCategorySQL.GetMenuSubCategoryByIDQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get query: %w", err)
	}

	row := h.db.QueryRow(query, id)
	return h.scanMenuSubCategoryRow(row)
}

// Create creates a new sub menu
func (h *DBHandler) Create(req *models.MenuSubCategoryCreateRequest) (*models.MenuSubCategory, error) {
	query, err := h.queries.Get(menuSubCategorySQL.CreateMenuSubCategoryQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get query: %w", err)
	}

	row := h.db.QueryRow(query,
		req.Name, req.Description, req.CategoryID,
		req.ItemType, req.DisplayOrder, req.IsActive,
	)

	return h.scanMenuSubCategoryRowWithoutCategory(row)
}

// Update updates an existing sub menu
func (h *DBHandler) Update(id string, req *models.MenuSubCategoryUpdateRequest) (*models.MenuSubCategory, error) {
	query, err := h.queries.Get(menuSubCategorySQL.UpdateMenuSubCategoryQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get query: %w", err)
	}

	row := h.db.QueryRow(query, id,
		req.Name, req.Description, req.CategoryID,
		req.ItemType, req.DisplayOrder, req.IsActive,
	)

	return h.scanMenuSubCategoryRowWithoutCategory(row)
}

// Delete deletes a sub menu
func (h *DBHandler) Delete(id string) error {
	// Check for dependencies first
	checkQuery, err := h.queries.Get(menuSubCategorySQL.CheckMenuSubCategoryDependenciesQuery)
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

	deleteQuery, err := h.queries.Get(menuSubCategorySQL.DeleteMenuSubCategoryQuery)
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
func (h *DBHandler) scanMenuSubCategory(rows *sql.Rows) (*models.MenuSubCategory, error) {
	var subMenu models.MenuSubCategory
	var description, categoryName sql.NullString

	err := rows.Scan(
		&subMenu.ID, &subMenu.Name, &description, &subMenu.CategoryID, &categoryName,
		&subMenu.ItemType, &subMenu.DisplayOrder, &subMenu.IsActive,
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

	return &subMenu, nil
}

func (h *DBHandler) scanMenuSubCategoryRow(row *sql.Row) (*models.MenuSubCategory, error) {
	var subMenu models.MenuSubCategory
	var description, categoryName sql.NullString

	err := row.Scan(
		&subMenu.ID, &subMenu.Name, &description, &subMenu.CategoryID, &categoryName,
		&subMenu.ItemType, &subMenu.DisplayOrder, &subMenu.IsActive,
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

	return &subMenu, nil
}

func (h *DBHandler) scanMenuSubCategoryRowWithoutCategory(row *sql.Row) (*models.MenuSubCategory, error) {
	var subMenu models.MenuSubCategory
	var description sql.NullString

	err := row.Scan(
		&subMenu.ID, &subMenu.Name, &description, &subMenu.CategoryID,
		&subMenu.ItemType, &subMenu.DisplayOrder, &subMenu.IsActive,
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

	return &subMenu, nil
}
