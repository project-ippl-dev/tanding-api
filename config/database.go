package config

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// DatabaseConnection is function to open connection to database
func NewDatabase(dbConf DatabaseConfig) (*sql.DB, error) {
	connection := fmt.Sprintf(`postgres://%s:%s@%s:%s/%s?sslmode=disable`, dbConf.Username, dbConf.Password, dbConf.Host, dbConf.Port, dbConf.Name)
	db, err := sql.Open("postgres", connection)
	if err != nil {
		return nil, fmt.Errorf("sql connect err : %s", err.Error())
	}
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("sql test connection err : %s", err.Error())
	}
	return db, nil
}
