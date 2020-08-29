package server

import (
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

func mySQL() (*gorm.DB, error) {
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
