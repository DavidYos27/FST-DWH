package models

import (
	"database/sql"
)

// // ProviderModel Interface for Calling Provider Model
// type ProviderModel interface {
// 	GetProviderID() ([]*Provider, error)
// }

// DB Structure for embedding SQL
type DB struct {
	SQL *sql.DB
}

// NewDB Create New Connection for DB
func NewDB(dataSourceName string) (*DB, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &DB{db}, nil
}
