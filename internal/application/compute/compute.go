package compute

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/nokamoto/demo20-apis/cloud/compute/v1alpha"
	"github.com/nokamoto/demo20-apps/internal/mysql/compute"
)

// Compute defines a business logic for the cloud compute service.
type Compute struct {
	instanceQuery instanceQuery
	db            *gorm.DB
}

// NewCompute returns Compute.
func NewCompute(db *gorm.DB) *Compute {
	return &Compute{
		instanceQuery: compute.Query{},
		db:            db,
	}
}

// RandomName returns a genrated name randomly from the parent id.
func (c *Compute) RandomName(parentID string) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz0123456789")

	b := make([]rune, 8)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	return fmt.Sprintf("%s-%s", parentID, string(b))
}

// Create creates a new compute instance.
func (c *Compute) Create(id, parentID string, instance *v1alpha.Instance) (*v1alpha.Instance, error) {
	err := c.db.Transaction(func(tx *gorm.DB) error {
		return c.instanceQuery.Create(tx, &compute.Instance{
			InstanceID: id,
			ParentID:   parentID,
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
