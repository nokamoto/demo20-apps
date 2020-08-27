package compute

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/nokamoto/demo20-apis/cloud/compute/v1alpha"
	"github.com/nokamoto/demo20-apps/internal/core"
	"github.com/nokamoto/demo20-apps/internal/mysql/compute"
	mysql "github.com/nokamoto/demo20-apps/internal/mysql/core"
)

// Compute defines a business logic for the cloud compute service.
type Compute struct {
	instanceQuery instanceQuery
	projectQuery  projectQuery
	db            *gorm.DB
}

// Create creates a new compute instance.
func (c *Compute) Create(id, parentID string, instance *v1alpha.Instance) (*v1alpha.Instance, error) {
	err := c.db.Transaction(func(tx *gorm.DB) error {
		project, err := c.projectQuery.Get(tx, parentID)
		if errors.Is(err, mysql.ErrNotFound) {
			return fmt.Errorf("%s: %w", parentID, core.ErrNotFound)
		}
		if err != nil {
			return err
		}

		return c.instanceQuery.Create(tx, &compute.Instance{
			InstanceID: id,
			ParentKey:  project.ProjectKey,
			Labels:     strings.Join(instance.GetLabels(), ","),
		})
	})
	if err != nil {
		return nil, err
	}

	return &v1alpha.Instance{
		Name:   fmt.Sprintf("instances/%s", id),
		Parent: fmt.Sprintf("projects/%s", parentID),
		Labels: instance.GetLabels(),
	}, nil
}
