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
	roleQuery        roleQuery
	roleBindingQuery roleBindingQuery
	db               *gorm.DB
}

// NewIam returns Iam.
func NewIam(db *gorm.DB) *Iam {
	return &Iam{
		permissionQuery:  &iam.PermissionQuery{},
		machineUserQuery: &iam.MachineUserQuery{},
		roleQuery:        &iam.RoleQuery{},
		roleBindingQuery: &iam.RoleBindingQuery{},
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

// FromMachineUserAPIKey returns a machine user id of the api key.
func (*Iam) FromMachineUserAPIKey(apiKey string) (string, error) {
	id, err := base64.RawStdEncoding.DecodeString(apiKey)
	if err != nil {
		return "", err
	}
	return string(id), nil
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
		Parent:      fmt.Sprintf("projects/%s", parentID),
	}, nil
}

// AuthenticateMachineUser authenticates with the api key.
func (i *Iam) AuthenticateMachineUser(apiKey string) (*v1alpha.MachineUser, error) {
	id, err := i.FromMachineUserAPIKey(apiKey)
	if err != nil {
		return nil, err
	}

	var machineUser *iam.MachineUser
	err = i.db.Transaction(func(tx *gorm.DB) error {
		res, err := i.machineUserQuery.Get(tx, id)
		machineUser = res
		return application.Error(err, application.ErrorMap{
			application.NotFound: id,
		})
	})
	if err != nil {
		return nil, err
	}

	return &v1alpha.MachineUser{
		Name:        fmt.Sprintf("machineusers/%s", machineUser.MachineUserID),
		DisplayName: machineUser.DisplayName,
		Parent:      fmt.Sprintf("projects/%s", machineUser.ParentID),
	}, nil
}

// AuthorizeMachineUser authorizes the machine.
func (i *Iam) AuthorizeMachineUser(machineUser *v1alpha.MachineUser, hierarchyProjectName string, permissionName string) (bool, error) {
	projectIDs := []string{application.FromNameF(hierarchyProjectName)}
	if hierarchyProjectName != "projects//" {
		projectIDs = append(projectIDs, "/")
	}

	var permissions []*iam.Permission

	err := i.db.Transaction(func(tx *gorm.DB) error {
		res, err := i.permissionQuery.Hierachy(tx, machineUser.GetName(), projectIDs)
		permissions = res
		return application.Error(err, application.ErrorMap{})
	})
	if err != nil {
		return false, err
	}

	for _, p := range permissions {
		if fmt.Sprintf("permissions/%s", p.PermissionID) == permissionName {
			return true, nil
		}
	}

	return false, nil
}

// CreateRole creates a role.
func (i *Iam) CreateRole(id, parentID string, role *v1alpha.Role) (*v1alpha.Role, error) {
	err := i.db.Transaction(func(tx *gorm.DB) error {
		var permissionIDs []string
		for _, p := range role.GetPermissions() {
			permissionIDs = append(permissionIDs, application.FromNameF(p))
		}

		permissions, err := i.permissionQuery.List(tx, permissionIDs...)
		if err != nil {
			return application.Error(err, application.ErrorMap{})
		}

		err = i.roleQuery.Create(tx, &iam.Role{
			RoleID:      id,
			ParentID:    parentID,
			DisplayName: role.GetDisplayName(),
		}, permissions...)

		return application.Error(err, application.ErrorMap{
			application.AlreadyExists: id,
		})
	})
	if err != nil {
		return nil, err
	}

	return &v1alpha.Role{
		Name:        fmt.Sprintf("roles/%s", id),
		DisplayName: role.GetDisplayName(),
		Permissions: role.GetPermissions(),
		Parent:      fmt.Sprintf("projects/%s", parentID),
	}, nil
}

// AddRoleBinding creates a role binding.
func (i *Iam) AddRoleBinding(parentID string, roleBinding *v1alpha.RoleBinding) (*v1alpha.RoleBinding, error) {
	err := i.db.Transaction(func(tx *gorm.DB) error {
		roleParentID, roleID := fromRoleNameF(roleBinding.GetRole())
		if parentID != "/" && roleParentID != parentID {
			return application.ErrPermissionDenied
		}

		role, _, err := i.roleQuery.Get(tx, roleID)
		if err != nil {
			return application.Error(err, application.ErrorMap{
				application.NotFound: roleID,
			})
		}

		err = i.roleBindingQuery.Create(tx, &iam.RoleBinding{
			RoleKey:  role.RoleKey,
			User:     roleBinding.GetUser(),
			ParentID: parentID,
		})
		return err
	})
	if err != nil {
		return nil, err
	}

	return &v1alpha.RoleBinding{
		Role:   roleBinding.GetRole(),
		User:   roleBinding.GetUser(),
		Parent: fmt.Sprintf("projects/%s", parentID),
	}, nil
}
