//go:generate mockgen -source=$GOFILE -destination=iam_mock.go -package=$GOPACKAGE
package iam

import (
	"github.com/nokamoto/demo20-apis/cloud/iam/v1alpha"
)

type iam interface {
	Create(string) (*v1alpha.Permission, error)
	CreateMachineUser(string, *v1alpha.MachineUser) (*v1alpha.MachineUser, error)
	AuthenticateMachineUser(string) (*v1alpha.MachineUser, error)
	AuthorizeMachineUser(*v1alpha.MachineUser, string, string) (bool, error)
	CreateRole(id, parentID string, role *v1alpha.Role) (*v1alpha.Role, error)
	AddRoleBinding(parentID string, roleBinding *v1alpha.RoleBinding) (*v1alpha.RoleBinding, error)
}
