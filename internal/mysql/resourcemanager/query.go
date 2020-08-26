package resourcemanager

import (
	"github.com/jinzhu/gorm"
	"github.com/nokamoto/demo20-apps/internal/mysql/core"
)

// Query defines queries for a project table within a transaction.
type Query struct{}

// Create inserts the project record.
func (Query) Create(tx *gorm.DB, project *Project) error {
	return core.Create(tx, project)
}

// Delete deletes a project record by the project id.
func (Query) Delete(tx *gorm.DB, id string) error {
	return core.Delete(tx, &Project{}, "project_id = ?", id)
}

// Update updates the project record.
func (Query) Update(tx *gorm.DB, project *Project) error {
	return core.Update(tx, project)
}

// Get returns a project record by the project id.
func (Query) Get(tx *gorm.DB, id string) (*Project, error) {
	var project Project
	return &project, core.Get(tx, &project, "project_id = ?", id)
}
