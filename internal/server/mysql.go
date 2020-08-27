package server

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
)

const (
	mysqlUser     = "MYSQL_USER"
	mysqlPassword = "MYSQL_PASSWORD"
	mysqlHost     = "MYSQL_HOST"
	mysqlPort     = "MYSQL_PORT"
	mysqlDatabase = "MYSQL_DATABASE"
)

// MySQL returns a database connection.
func MySQL() (*gorm.DB, error) {
	uri := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		os.Getenv(mysqlUser),
		os.Getenv(mysqlPassword),
		os.Getenv(mysqlHost),
		os.Getenv(mysqlPort),
		os.Getenv(mysqlDatabase),
	)
	return gorm.Open("mysql", uri)
}
