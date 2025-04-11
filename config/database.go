package config

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

// DatabaseConnection is function to open connection to database
func DatabaseConnection() (*sql.DB, error) {
	dbConf := Configuration().Database
	connection := fmt.Sprintf(`%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Asia%%2FJakarta`, dbConf.Username, dbConf.Password, dbConf.Host, dbConf.Port, dbConf.Name)
	return sql.Open("postgres", connection)
}
