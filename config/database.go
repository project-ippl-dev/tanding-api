package config

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// DatabaseConnection is function to open connection to database
func DatabaseConnection() (*sql.DB, error) {
	dbConf := Configuration().Database
	connection := fmt.Sprintf(`postgres://%s:%s@%s:%s/%s?sslmode=disable`, dbConf.Username, dbConf.Password, dbConf.Host, dbConf.Port, dbConf.Name)
	db, err := sql.Open("postgres", connection)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
