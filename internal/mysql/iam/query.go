package iam

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/nokamoto/demo20-apps/internal/mysql"
)

// PermissionQuery defines quries for a permission table within a transaction.
type PermissionQuery struct{}

// Create inserts the permission record.
func (PermissionQuery) Create(tx *gorm.DB, permission *Permission) error {
	return mysql.Create(tx, permission)
}

// List returns permission records by permission ids.
func (PermissionQuery) List(tx *gorm.DB, permissionIDs ...string) ([]*Permission, error) {
	var permissions []*Permission
	res := tx.Debug().Where("permission_id in (?)", permissionIDs).Find(&permissions)
	if res.Error != nil {
		return nil, mysql.Translate(res.Error)
	}
	if len(permissionIDs) != len(permissions) {
		return nil, fmt.Errorf("permission: %w", mysql.ErrNotFound)
	}
	return permissions, nil
}

// RoleQuery defines quries for role tables within a transaction.
type RoleQuery struct{}

func newBulkRolePermission(roleKey int64, permissions ...*Permission) bulkRolePermission {
	var bulk bulkRolePermission
	for _, p := range permissions {
		bulk = append(bulk, &RolePermission{
			RoleKey:       roleKey,
			PermissionKey: p.PermissionKey,
		})
	}
	return bulk
}

// Create inserts role and role-permission records.
func (q RoleQuery) Create(tx *gorm.DB, role *Role, permissions ...*Permission) error {
	res := tx.Debug().Create(role)
	if res.Error != nil {
		return mysql.Translate(res.Error)
	}

	return mysql.BulkInsert(
		tx,
		RolePermission{}.TableName(),
		newBulkRolePermission(role.RoleKey, permissions...),
	)
}

// Delete deletes role and role-permission records.
func (RoleQuery) Delete(tx *gorm.DB, role *Role) error {
	res := tx.Debug().Where("role_key = ?", role.RoleKey).Delete(&RolePermission{})
	if res.Error != nil {
		return mysql.Translate(res.Error)
	}

	res = tx.Debug().Where("role_key = ?", role.RoleKey).Delete(&Role{})
	if res.Error != nil {
		return mysql.Translate(res.Error)
	}

	return nil
}

// Get returns role and role-permission records by the role id.
func (RoleQuery) Get(tx *gorm.DB, id string) (*Role, []*RolePermission, error) {
	var role Role
	res := tx.Debug().Where("role_id = ?", id).Take(&role)
	if res.RecordNotFound() {
		return nil, nil, mysql.ErrNotFound
	}
	if res.Error != nil {
		return nil, nil, mysql.Translate(res.Error)
	}

	var permissions []*RolePermission
	res = tx.Debug().Where("role_key = ?", role.RoleKey).Find(&permissions)
	if res.Error != nil {
		return nil, nil, mysql.Translate(res.Error)
	}

	return &role, permissions, nil
}

// Update updates role and role-permission records.
func (q RoleQuery) Update(tx *gorm.DB, role *Role, permissions ...*Permission) error {
	res := tx.Debug().Save(role)
	if res.Error != nil {
		return mysql.Translate(res.Error)
	}

	res = tx.Debug().Where("role_key = ?", role.RoleKey).Delete(&RolePermission{})
	if res.Error != nil {
		return mysql.Translate(res.Error)
	}

	return mysql.BulkInsert(
		tx,
		RolePermission{}.TableName(),
		newBulkRolePermission(role.RoleKey, permissions...),
	)
}

// List returns role and role-permission records by the parent key.
func (RoleQuery) List(tx *gorm.DB, parentID string, offset, limit int) ([]*Role, []*RolePermission, error) {
	var roles []*Role
	err := mysql.List(tx, &roles, offset, limit, "parent_id = ?", parentID)
	if err != nil {
		return nil, nil, err
	}

	if len(roles) == 0 {
		return nil, nil, nil
	}

	var keys []interface{}
	for _, r := range roles {
		keys = append(keys, r.RoleKey)
	}

	var permissions []*RolePermission
	return roles, permissions, mysql.ListAll(tx, &permissions, "role_key in (?)", keys...)
}

// RoleBindingQuery defines queries for a role binding table within a transaction.
type RoleBindingQuery struct{}

// Create inserts the role binding record.
func (RoleBindingQuery) Create(tx *gorm.DB, roleBinding *RoleBinding) error {
	return mysql.Create(tx, roleBinding)
}

// Delete deletes the role binding record.
func (RoleBindingQuery) Delete(tx *gorm.DB, roleBinding *RoleBinding) error {
	return mysql.Delete(
		tx,
		&RoleBinding{},
		"role_key = ? and user = ? and parent_key = ?",
		roleBinding.RoleKey, roleBinding.User, roleBinding.ParentKey,
	)
}

// List returns role binding records by the parent key.
func (RoleBindingQuery) List(tx *gorm.DB, parentKey int64, offset, limit int) ([]*RoleBinding, error) {
	var roleBindings []*RoleBinding
	return roleBindings, mysql.List(tx, &roleBindings, offset, limit, "parent_key = ?", parentKey)
}

// MachineUserQuery defines queries for a machine user table within a transaction.
type MachineUserQuery struct{}

// Create inserts the machine user record.
func (MachineUserQuery) Create(tx *gorm.DB, machineUser *MachineUser) error {
	return mysql.Create(tx, machineUser)
}
