package handlers

import (
	"database/sql"
	"fmt"
	"menu-service/pkg/entities/menu_categories/models"
	menuCategorySQL "menu-service/pkg/entities/menu_categories/sql"
	sharedDb "shared/db"

	"github.com/sirupsen/logrus"
)

// DBHandler handles database operations for menu categories
type DBHandler struct {
	db      *sharedDb.DbHandler
	queries *menuCategorySQL.Queries
	logger  *logrus.Logger
}

// NewDBHandler creates a new database handler
func NewDBHandler(db *sharedDb.DbHandler, logger *logrus.Logger) (*DBHandler, error) {
	queries, err := menuCategorySQL.LoadQueries()
	if err != nil {
		return nil, fmt.Errorf("failed to load SQL queries: %w", err)
	}

	return &DBHandler{
		db:      db,
		queries: queries,
		logger:  logger,
	}, nil
}

// List returns a paginated list of menu categories
func (h *DBHandler) List(page, limit int) (*models.MenuCategoryListResponse, error) {
	offset := (page - 1) * limit

	// Get total count
	countQuery, err := h.queries.Get(menuCategorySQL.CountMenuCategoriesQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get count query: %w", err)
	}

	var total int
	if err := h.db.QueryRow(countQuery).Scan(&total); err != nil {
		return nil, fmt.Errorf("failed to count menu categories: %w", err)
	}

	// Get categories
	listQuery, err := h.queries.Get(menuCategorySQL.ListMenuCategoriesQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get list query: %w", err)
	}

	rows, err := h.db.Query(listQuery, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list menu categories: %w", err)
	}
	defer rows.Close()

	var categories []models.MenuCategory
	for rows.Next() {
		var cat models.MenuCategory
		var description sql.NullString

		if err := rows.Scan(&cat.ID, &cat.Name, &cat.DisplayOrder, &description, &cat.CreatedAt, &cat.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan menu category: %w", err)
		}

		if description.Valid {
			cat.Description = &description.String
		}
		categories = append(categories, cat)
	}

	return &models.MenuCategoryListResponse{
		Categories: categories,
		Total:      total,
		Page:       page,
		Limit:      limit,
	}, nil
}

// GetByID returns a menu category by ID
func (h *DBHandler) GetByID(id string) (*models.MenuCategory, error) {
	query, err := h.queries.Get(menuCategorySQL.GetMenuCategoryByIDQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get query: %w", err)
	}

	var cat models.MenuCategory
	var description sql.NullString

	err = h.db.QueryRow(query, id).Scan(&cat.ID, &cat.Name, &cat.DisplayOrder, &description, &cat.CreatedAt, &cat.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get menu category: %w", err)
	}

	if description.Valid {
		cat.Description = &description.String
	}

	return &cat, nil
}

// Create creates a new menu category
func (h *DBHandler) Create(req *models.MenuCategoryCreateRequest) (*models.MenuCategory, error) {
	query, err := h.queries.Get(menuCategorySQL.CreateMenuCategoryQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get query: %w", err)
	}

	var cat models.MenuCategory
	var description sql.NullString

	err = h.db.QueryRow(query, req.Name, req.DisplayOrder, req.Description).Scan(
		&cat.ID, &cat.Name, &cat.DisplayOrder, &description, &cat.CreatedAt, &cat.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create menu category: %w", err)
	}

	if description.Valid {
		cat.Description = &description.String
	}

	h.logger.WithField("id", cat.ID).Info("Menu category created")
	return &cat, nil
}

// Update updates an existing menu category
func (h *DBHandler) Update(id string, req *models.MenuCategoryUpdateRequest) (*models.MenuCategory, error) {
	query, err := h.queries.Get(menuCategorySQL.UpdateMenuCategoryQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get query: %w", err)
	}

	var cat models.MenuCategory
	var description sql.NullString

	err = h.db.QueryRow(query, id, req.Name, req.DisplayOrder, req.Description).Scan(
		&cat.ID, &cat.Name, &cat.DisplayOrder, &description, &cat.CreatedAt, &cat.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to update menu category: %w", err)
	}

	if description.Valid {
		cat.Description = &description.String
	}

	h.logger.WithField("id", cat.ID).Info("Menu category updated")
	return &cat, nil
}

// Delete deletes a menu category
func (h *DBHandler) Delete(id string) error {
	// Check for dependencies first
	checkQuery, err := h.queries.Get(menuCategorySQL.CheckMenuCategoryDependenciesQuery)
	if err != nil {
		return fmt.Errorf("failed to get check query: %w", err)
	}

	var count int
	if err := h.db.QueryRow(checkQuery, id).Scan(&count); err != nil {
		return fmt.Errorf("failed to check dependencies: %w", err)
	}

	if count > 0 {
		return fmt.Errorf("cannot delete category: %d menu items depend on it", count)
	}

	// Delete the category
	deleteQuery, err := h.queries.Get(menuCategorySQL.DeleteMenuCategoryQuery)
	if err != nil {
		return fmt.Errorf("failed to get delete query: %w", err)
	}

	result, err := h.db.Exec(deleteQuery, id)
	if err != nil {
		return fmt.Errorf("failed to delete menu category: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("menu category not found")
	}

	h.logger.WithField("id", id).Info("Menu category deleted")
	return nil
}
