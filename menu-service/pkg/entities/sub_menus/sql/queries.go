package sql

import (
	"embed"
	"fmt"
)

//go:embed scripts/*.sql
var sqlScripts embed.FS

// Query names
const (
	ListSubMenusQuery               = "list_sub_menus"
	CountSubMenusQuery              = "count_sub_menus"
	GetSubMenuByIDQuery             = "get_sub_menu_by_id"
	CreateSubMenuQuery              = "create_sub_menu"
	UpdateSubMenuQuery              = "update_sub_menu"
	DeleteSubMenuQuery              = "delete_sub_menu"
	CheckSubMenuDependenciesQuery   = "check_sub_menu_dependencies"
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
		"list_sub_menus.sql",
		"count_sub_menus.sql",
		"get_sub_menu_by_id.sql",
		"create_sub_menu.sql",
		"update_sub_menu.sql",
		"delete_sub_menu.sql",
		"check_sub_menu_dependencies.sql",
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
