package models

import (
	"database/sql"
)
// Provider Structure for Model Provider
type Provider struct {
	ProviderID   int
	ProviderName string
	CreatedAt    string
	UpdatedAt    sql.NullString
}

// GetProviderID Get Provider ID from Provider Table
func (db *DB) GetProviderID() ([]*Provider, error) {
	rows, err := db.Query("SELECT * FROM DIMENSION.PROVIDER")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	providers := make([]*Provider, 0)

	for rows.Next() {
		p := new(Provider)
		err := rows.Scan(&p.ProviderID, &p.ProviderName, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}
		providers = append(providers, p)
	}
	return providers,err
}
