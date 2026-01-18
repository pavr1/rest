package sql

import (
	"embed"
	"fmt"
)

//go:embed scripts/*.sql
var sqlScripts embed.FS

// Query names
const (
	ListMenuSubCategoriesQuery      = "list_menu_sub_categories"
	CountMenuSubCategoriesQuery     = "count_menu_sub_categories"
	GetMenuSubCategoryByIDQuery     = "get_menu_sub_category_by_id"
	CreateMenuSubCategoryQuery      = "create_menu_sub_category"
	UpdateMenuSubCategoryQuery      = "update_menu_sub_category"
	DeleteMenuSubCategoryQuery      = "delete_menu_sub_category"
	CheckMenuSubCategoryDependenciesQuery = "check_menu_sub_category_dependencies"
)

// Queries holds all loaded SQL queries
type Queries struct {
	queries map[string]string
}

// LoadQueries loads all SQL queries from embedded files
func LoadQueries() (*Queries, error) {
	q := &Queries{
		queries: make(map[string]string),
	}

	files := []string{
		"list_menu_sub_categories.sql",
		"count_menu_sub_categories.sql",
		"get_menu_sub_category_by_id.sql",
		"create_menu_sub_category.sql",
		"update_menu_sub_category.sql",
		"delete_menu_sub_category.sql",
		"check_menu_sub_category_dependencies.sql",
	}

	for _, file := range files {
		content, err := sqlScripts.ReadFile("scripts/" + file)
		if err != nil {
			return nil, fmt.Errorf("failed to read SQL file %s: %w", file, err)
		}
		// Remove .sql extension for query name
		queryName := file[:len(file)-4]
		q.queries[queryName] = string(content)
	}

	return q, nil
}

// Get returns a query by name
func (q *Queries) Get(name string) (string, error) {
	query, ok := q.queries[name]
	if !ok {
		return "", fmt.Errorf("query not found: %s", name)
	}
	return query, nil
}
