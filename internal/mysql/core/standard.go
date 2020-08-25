package core

import (
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
