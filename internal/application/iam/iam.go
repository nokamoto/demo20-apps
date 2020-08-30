package iam

import (
	"encoding/base64"
	"fmt"
	"math/rand"

	"github.com/jinzhu/gorm"
	"github.com/nokamoto/demo20-apis/cloud/iam/v1alpha"
	"github.com/nokamoto/demo20-apps/internal/application"
	"github.com/nokamoto/demo20-apps/internal/mysql/iam"
)

// Iam defines a business logic for the cloud iam service.
type Iam struct {
	permissionQuery  permissionQuery
	machineUserQuery machineUserQuery
	db               *gorm.DB
}

// NewIam returns Iam.
func NewIam(db *gorm.DB) *Iam {
	return &Iam{
		permissionQuery:  &iam.PermissionQuery{},
		machineUserQuery: &iam.MachineUserQuery{},
		db:               db,
	}
}

// Create creates a permission.
func (i *Iam) Create(id string) (*v1alpha.Permission, error) {
	err := i.db.Transaction(func(tx *gorm.DB) error {
		err := i.permissionQuery.Create(tx, &iam.Permission{
			PermissionID: id,
		})
		return application.Error(err, application.ErrorMap{
			application.AlreadyExists: id,
		})
	})
	if err != nil {
		return nil, err
	}

	return &v1alpha.Permission{
		Name: fmt.Sprintf("permissions/%s", id),
	}, nil
}

// RandomMachineUserID returns a genrated name randomly from the machine user id.
func (*Iam) RandomMachineUserID() string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz0123456789")

	b := make([]rune, 32)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	return string(b)
}

// MachineUserAPIKey returns an api key of the machine user id.
func (*Iam) MachineUserAPIKey(id string) string {
	return base64.RawStdEncoding.EncodeToString([]byte(id))
}

// CreateMachineUser creates a machine user.
func (i *Iam) CreateMachineUser(parentID string, machineUser *v1alpha.MachineUser) (*v1alpha.MachineUser, error) {
	id := i.RandomMachineUserID()

	err := i.db.Transaction(func(tx *gorm.DB) error {
		err := i.machineUserQuery.Create(tx, &iam.MachineUser{
			MachineUserID: id,
			DisplayName:   machineUser.GetDisplayName(),
			ParentID:      parentID,
		})
		return application.Error(err, application.ErrorMap{
			application.AlreadyExists: id,
		})
	})
	if err != nil {
		return nil, err
	}

	return &v1alpha.MachineUser{
		Name:        fmt.Sprintf("machineusers/%s", id),
		DisplayName: machineUser.GetDisplayName(),
		ApiKey:      i.MachineUserAPIKey(id),
	}, nil
}
