package compute

import (
	"github.com/jinzhu/gorm"
	"github.com/nokamoto/demo20-apps/internal/mysql/core"
)

// Query defines queries for a instance table within a transaction.
type Query struct{}

// Create inserts an instance record.
func (Query) Create(tx *gorm.DB, instance *Instance) error {
	res := tx.Debug().Create(instance)
	if res.Error != nil {
		return core.Translate(res.Error)
	}
	return nil
}

// Delete deletes an instance record by an instance id.
func (Query) Delete(tx *gorm.DB, id string) error {
	res := tx.Debug().Where("instance_id = ?", id).Delete(&Instance{})
	if res.Error != nil {
		return core.Translate(res.Error)
	}
	return nil
}

// Get returns an instance record by an instance id.
func (Query) Get(tx *gorm.DB, id string) (*Instance, error) {
	var instance Instance
	res := tx.Debug().Where("instance_id = ?", id).Take(&instance)
	if res.RecordNotFound() {
		return nil, core.ErrNotFound
	}
	if res.Error != nil {
		return nil, core.Translate(res.Error)
	}
	return &instance, nil
}

// List returns instance records by a parent id.
func (Query) List(tx *gorm.DB, parentKey int64, offset int, limit int) ([]*Instance, error) {
	var instances []*Instance
	res := tx.Debug().Where("parent_key = ?", parentKey).Offset(offset).Limit(limit).Find(&instances)
	if res.Error != nil {
		return nil, core.Translate(res.Error)
	}
	return instances, nil
}
