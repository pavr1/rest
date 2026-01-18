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
	ListMenuVariantsQuery           = "list_menu_variants"
	CountMenuVariantsQuery          = "count_menu_variants"
	GetMenuVariantByIDQuery         = "get_menu_variant_by_id"
	CreateMenuVariantQuery          = "create_menu_variant"
	UpdateMenuVariantQuery          = "update_menu_variant"
	DeleteMenuVariantQuery          = "delete_menu_variant"
	UpdateMenuVariantAvailabilityQuery = "update_menu_variant_availability"
	UpdateMenuVariantImageQuery     = "update_menu_variant_image"
	UpdateMenuVariantCostQuery      = "update_menu_variant_cost"
)
