//go:generate mockgen -source=$GOFILE -destination=query_mock.go -package=$GOPACKAGE
package iam

import (
	"github.com/jinzhu/gorm"
	"github.com/nokamoto/demo20-apps/internal/mysql/iam"
)

type permissionQuery interface {
	Create(*gorm.DB, *iam.Permission) error
	List(tx *gorm.DB, permissionIDs ...string) ([]*iam.Permission, error)
	Hierachy(tx *gorm.DB, user string, hierarchy []string) ([]*iam.Permission, error)
}

type machineUserQuery interface {
	Create(*gorm.DB, *iam.MachineUser) error
	Get(*gorm.DB, string) (*iam.MachineUser, error)
}

type roleQuery interface {
	Create(tx *gorm.DB, role *iam.Role, permissions ...*iam.Permission) error
	Get(tx *gorm.DB, id string) (*iam.Role, []*iam.RolePermission, error)
}

type roleBindingQuery interface {
	Create(tx *gorm.DB, roleBinding *iam.RoleBinding) error
}
