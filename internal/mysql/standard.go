package mysql

import (
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
)

// Create inserts the record.
func Create(tx *gorm.DB, record interface{}) error {
	res := tx.Debug().Create(record)
	if res.Error != nil {
		return Translate(res.Error)
	}
	return nil
}

// Delete deletes records by the condition(s).
func Delete(tx *gorm.DB, record interface{}, where string, args ...interface{}) error {
	res := tx.Debug().Where(where, args).Delete(record)
	if res.Error != nil {
		return Translate(res.Error)
	}
	return nil
}

// Get finds a single record by the condition(s).
func Get(tx *gorm.DB, record interface{}, where string, args ...interface{}) error {
	res := tx.Debug().Where(where, args).Take(record)
	if res.RecordNotFound() {
		return ErrNotFound
	}
	if res.Error != nil {
		return Translate(res.Error)
	}
	return nil
}

// List finds records by the condition(s).
func List(tx *gorm.DB, records interface{}, offset, limit int, where string, args ...interface{}) error {
	res := tx.Debug().Where(where, args).Offset(offset).Limit(limit).Find(records)
	if res.Error != nil {
		return Translate(res.Error)
	}
	return nil
}

// ListAll finds all records by the condition(s).
func ListAll(tx *gorm.DB, records interface{}, where string, args ...interface{}) error {
	res := tx.Debug().Where(where, args).Find(records)
	if res.Error != nil {
		return Translate(res.Error)
	}
	return nil
}

type bulkable interface {
	Args() [][]interface{}
	Fields() []string
}

// BulkInsert inserts multiple recoreds.
func BulkInsert(tx *gorm.DB, tableName string, records bulkable) error {
	fields := records.Fields()
	var escapedFields []string
	for _, f := range fields {
		escapedFields = append(escapedFields, fmt.Sprintf("`%s`", f))
	}

	var q []string
	for i := 0; i < len(fields); i++ {
		q = append(q, "?")
	}

	placeholder := fmt.Sprintf("(%s)", strings.Join(q, ","))

	var placeholders []string
	var args []interface{}
	for _, r := range records.Args() {
		placeholders = append(placeholders, placeholder)
		args = append(args, r...)
	}

	query := fmt.Sprintf(
		"INSERT INTO `%s` (%s) VALUES %s",
		tableName, strings.Join(escapedFields, ","), strings.Join(placeholders, ","),
	)

	res := tx.Debug().Exec(query, args...)
	if res.Error != nil {
		return Translate(res.Error)
	}

	return nil
}

// Update updates the record.
func Update(tx *gorm.DB, record interface{}) error {
	res := tx.Debug().Save(record)
	if res.Error != nil {
		return Translate(res.Error)
	}
	return nil
}
