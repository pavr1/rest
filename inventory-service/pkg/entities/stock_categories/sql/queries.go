package sql

import (
	"embed"
	"fmt"
	"strings"
)

//go:embed scripts/*.sql
var sqlFiles embed.FS

// Queries holds all SQL queries
type Queries struct {
	queries map[string]string
}

// LoadQueries loads all SQL queries from the scripts directory
func LoadQueries() (*Queries, error) {
	queries := &Queries{
		queries: make(map[string]string),
	}

	files, err := sqlFiles.ReadDir("scripts")
	if err != nil {
		return nil, fmt.Errorf("failed to read SQL scripts directory: %w", err)
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			content, err := sqlFiles.ReadFile("scripts/" + file.Name())
			if err != nil {
				return nil, fmt.Errorf("failed to read SQL file %s: %w", file.Name(), err)
			}

			queryName := strings.TrimSuffix(file.Name(), ".sql")
			queries.queries[queryName] = string(content)
		}
	}
	return queries, nil
}

// Get retrieves a query by name
func (q *Queries) Get(name string) (string, error) {
	query, exists := q.queries[name]
	if !exists {
		return "", fmt.Errorf("query '%s' not found", name)
	}
	return query, nil
}

// SQL query constants
const (
	ListStockCategoriesQuery           = "list_stock_categories"
	CountStockCategoriesQuery          = "count_stock_categories"
	GetStockCategoryByIDQuery          = "get_stock_category_by_id"
	CreateStockCategoryQuery           = "create_stock_category"
	UpdateStockCategoryQuery           = "update_stock_category"
	DeleteStockCategoryQuery           = "delete_stock_category"
	CheckStockCategoryDependenciesQuery = "check_stock_category_dependencies"
)
