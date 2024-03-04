package app

import (
	"database/sql"
)

// ConnectToPostgres connect to postgres
func ConnectToPostgres(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dataSourceName)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
