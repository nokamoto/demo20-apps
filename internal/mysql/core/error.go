package core

import (
	"errors"

	"github.com/jinzhu/gorm"

	"github.com/go-sql-driver/mysql"
)

// https://dev.mysql.com/doc/mysql-errors/8.0/en/server-error-reference.html
const (
	DupEntry = 1062
)

var (
	// ErrAlreadyExists represents a duplicated record.
	ErrAlreadyExists = errors.New("already exists")
	// ErrNotFound represents a record not found.
	ErrNotFound = errors.New("not found")
)

// Translate converts a mysql server error to an error defined in this package.
func Translate(err error) error {
	if gorm.IsRecordNotFoundError(err) {
		return ErrNotFound
	}
	if e, ok := err.(*mysql.MySQLError); ok {
		switch e.Number {
		case DupEntry:
			return ErrAlreadyExists
		}
	}
	return err
}
