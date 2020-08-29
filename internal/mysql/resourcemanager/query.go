package resourcemanager

import (
	"github.com/jinzhu/gorm"
	"github.com/nokamoto/demo20-apps/internal/mysql"
)

// Query defines queries for a project table within a transaction.
type Query struct{}

// Create inserts the project record.
func (Query) Create(tx *gorm.DB, project *Project) error {
	return mysql.Create(tx, project)
}

// Delete deletes a project record by the project id.
func (Query) Delete(tx *gorm.DB, id string) error {
	return mysql.Delete(tx, &Project{}, "project_id = ?", id)
}

// Update updates the project record.
func (Query) Update(tx *gorm.DB, project *Project) error {
	return mysql.Update(tx, project)
}

// Get returns a project record by the project id.
func (Query) Get(tx *gorm.DB, id string) (*Project, error) {
	var project Project
	return &project, mysql.Get(tx, &project, "project_id = ?", id)
}
