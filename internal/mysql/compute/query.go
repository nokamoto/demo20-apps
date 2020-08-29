package compute

import (
	"github.com/jinzhu/gorm"
	"github.com/nokamoto/demo20-apps/internal/mysql"
)

// Query defines queries for a instance table within a transaction.
type Query struct{}

// Create inserts the instance record.
func (Query) Create(tx *gorm.DB, instance *Instance) error {
	return mysql.Create(tx, instance)
}

// Delete deletes an instance record by the instance id.
func (Query) Delete(tx *gorm.DB, id string) error {
	return mysql.Delete(tx, &Instance{}, "instance_id = ?", id)
}

// Get returns an instance record by the instance id.
func (Query) Get(tx *gorm.DB, id string) (*Instance, error) {
	var instance Instance
	return &instance, mysql.Get(tx, &instance, "instance_id = ?", id)
}

// List returns instance records by the parent id.
func (Query) List(tx *gorm.DB, parentID string, offset int, limit int) ([]*Instance, error) {
	var instances []*Instance
	return instances, mysql.List(tx, &instances, offset, limit, "parent_id = ?", parentID)
}
