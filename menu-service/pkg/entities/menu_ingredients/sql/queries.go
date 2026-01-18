package sql

import (
	"embed"
	"fmt"
)

//go:embed scripts/*.sql
var sqlScripts embed.FS

// Query names
const (
	ListMenuIngredientsQuery          = "list_menu_ingredients"
	GetMenuIngredientByIDQuery        = "get_menu_ingredient_by_id"
	CreateMenuIngredientQuery         = "create_menu_ingredient"
	UpdateMenuIngredientQuery         = "update_menu_ingredient"
	DeleteMenuIngredientQuery         = "delete_menu_ingredient"
	GetIngredientsByMenuVariantQuery  = "get_ingredients_by_menu_variant"
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
		"list_menu_ingredients.sql",
		"get_menu_ingredient_by_id.sql",
		"create_menu_ingredient.sql",
		"update_menu_ingredient.sql",
		"delete_menu_ingredient.sql",
		"get_ingredients_by_menu_variant.sql",
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
