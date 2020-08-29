package resourcemanager

import (
	"fmt"

	"github.com/nokamoto/demo20-apps/internal/mysql/resourcemanager"

	"github.com/jinzhu/gorm"
	"github.com/nokamoto/demo20-apis/cloud/resourcemanager/v1alpha"
	"github.com/nokamoto/demo20-apps/internal/application"
)

// ResourceManager defines a business logic for the cloud resource manager service.
type ResourceManager struct {
	projectQuery projectQuery
	db           *gorm.DB
}

// NewResourceManager returns ResourceManager.
func NewResourceManager(db *gorm.DB) *ResourceManager {
	return &ResourceManager{
		projectQuery: resourcemanager.Query{},
		db:           db,
	}
}

// Get returns a project by the project id.
func (r *ResourceManager) Get(id string) (*v1alpha.Project, error) {
	var project v1alpha.Project
	err := r.db.Transaction(func(tx *gorm.DB) error {
		res, err := r.projectQuery.Get(tx, id)
		if err != nil {
			return application.Error(err, application.ErrorMap{
				application.NotFound: id,
			})
		}

		project = v1alpha.Project{
			Name:        fmt.Sprintf("projects/%s", res.ProjectID),
			DisplayName: res.DisplayName,
		}

		return nil
	})
	return &project, err
}

// Create creates a project.
func (r *ResourceManager) Create(id string, project *v1alpha.Project) (*v1alpha.Project, error) {
	var res v1alpha.Project
	err := r.db.Transaction(func(tx *gorm.DB) error {
		err := r.projectQuery.Create(tx, &resourcemanager.Project{
			ProjectID:   id,
			DisplayName: project.GetDisplayName(),
		})
		if err != nil {
			return application.Error(err, application.ErrorMap{
				application.AlreadyExists: id,
			})
		}

		res = v1alpha.Project{
			Name:        fmt.Sprintf("projects/%s", id),
			DisplayName: project.GetDisplayName(),
		}

		return nil
	})
	return &res, err
}
