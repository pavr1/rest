package settings

// import (
// 	"data-service/pkg/database"
// 	"data-service/pkg/entities/settings/sql"
// 	sharedModels "shared/models"
// )

// // Repository handles database operations for settings
// type Repository struct {
// 	db      database.DatabaseHandler
// 	queries *sql.Queries
// }

// // NewRepository creates a new settings repository
// func NewRepository(db database.DatabaseHandler) (*Repository, error) {
// 	queries, err := sql.LoadQueries()
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &Repository{
// 		db:      db,
// 		queries: queries,
// 	}, nil
// }

// // GetSettingsByService retrieves all settings for a specific service
// func (r *Repository) GetSettingsByService(service string) ([]sharedModels.Setting, error) {
// 	query, err := r.queries.Get(sql.GetSettingsByServiceQuery)
// 	if err != nil {
// 		return nil, err
// 	}

// 	rows, err := r.db.Query(query, service)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var settings []sharedModels.Setting
// 	for rows.Next() {
// 		var s sharedModels.Setting
// 		if err := rows.Scan(&s.SettingID, &s.Service, &s.Key, &s.Value, &s.Description, &s.CreatedAt, &s.UpdatedAt); err != nil {
// 			return nil, err
// 		}
// 		settings = append(settings, s)
// 	}

// 	if err := rows.Err(); err != nil {
// 		return nil, err
// 	}

// 	return settings, nil
// }

// // GetSettingByKey retrieves a specific setting by service and key
// func (r *Repository) GetSettingByKey(service, key string) (*sharedModels.Setting, error) {
// 	query, err := r.queries.Get(sql.GetSettingByKeyQuery)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var s sharedModels.Setting
// 	err = r.db.QueryRow(query, service, key).Scan(
// 		&s.SettingID, &s.Service, &s.Key, &s.Value, &s.Description, &s.CreatedAt, &s.UpdatedAt,
// 	)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &s, nil
// }

// // UpdateSetting updates a setting value
// func (r *Repository) UpdateSetting(service, key, value string) error {
// 	query, err := r.queries.Get(sql.UpdateSettingQuery)
// 	if err != nil {
// 		return err
// 	}

// 	_, err = r.db.Exec(query, value, service, key)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// // CreateSetting creates a new setting
// func (r *Repository) CreateSetting(setting sharedModels.Setting) error {
// 	query, err := r.queries.Get(sql.CreateSettingQuery)
// 	if err != nil {
// 		return err
// 	}

// 	_, err = r.db.Exec(query, setting.Service, setting.Key, setting.Value, setting.Description)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// // DeleteSetting deletes a setting
// func (r *Repository) DeleteSetting(service, key string) error {
// 	query, err := r.queries.Get(sql.DeleteSettingQuery)
// 	if err != nil {
// 		return err
// 	}

// 	_, err = r.db.Exec(query, service, key)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
