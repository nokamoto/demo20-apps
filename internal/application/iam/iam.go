package iam

import (
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/nokamoto/demo20-apis/cloud/iam/v1alpha"
	"github.com/nokamoto/demo20-apps/internal/application"
	"github.com/nokamoto/demo20-apps/internal/mysql"
	"github.com/nokamoto/demo20-apps/internal/mysql/iam"
)

// Iam defines a business logic for the cloud iam service.
type Iam struct {
	permissionQuery permissionQuery
	db              *gorm.DB
}

// NewIam returns Iam.
func NewIam(db *gorm.DB) *Iam {
	return &Iam{
		permissionQuery: &iam.PermissionQuery{},
		db:              db,
	}
}

// Create creates a permission.
func (i *Iam) Create(id string) (*v1alpha.Permission, error) {
	err := i.db.Transaction(func(tx *gorm.DB) error {
		err := i.permissionQuery.Create(tx, &iam.Permission{
			PermissionID: id,
		})
		if errors.Is(err, mysql.ErrAlreadyExists) {
			return fmt.Errorf("%s: %w", id, application.ErrAlreadyExists)
		}
		return err
	})
	if err != nil {
		return nil, nil
	}

	return &v1alpha.Permission{
		Name: fmt.Sprintf("permissions/%s", id),
	}, nil
}
