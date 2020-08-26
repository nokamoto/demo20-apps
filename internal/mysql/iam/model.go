package iam

// Permission represents a permission table.
type Permission struct {
	PermissionKey int64  `gorm:"column:permission_key;auto_increment;primary_key"`
	PermissionID  string `gorm:"column:permission_id"`
}

// TableName returns a table name.
func (Permission) TableName() string {
	return "iam_permission"
}

// Role represents a role table.
type Role struct {
	RoleKey     int64  `gorm:"column:role_key;auto_increment;primary_key"`
	RoleID      string `gorm:"column:role_id"`
	ParentKey   int64  `gorm:"column:parent_key"`
	DisplayName string `gorm:"column:display_name"`
}

// TableName returns a table name.
func (Role) TableName() string {
	return "iam_role"
}

// RolePermission represents a role_permission table.
type RolePermission struct {
	RoleKey       int64 `gorm:"column:role_key"`
	PermissionKey int64 `gorm:"column:permission_key"`
}

// TableName returns a table name.
func (RolePermission) TableName() string {
	return "iam_role_permission"
}

type bulkRolePermission []*RolePermission

func (xs bulkRolePermission) Args() [][]interface{} {
	var res [][]interface{}
	for _, x := range xs {
		var args []interface{}
		args = append(args, x.RoleKey)
		args = append(args, x.PermissionKey)
		res = append(res, args)
	}
	return res
}

func (xs bulkRolePermission) Fields() []string {
	return []string{"role_key", "permission_key"}
}

// RoleBinding represents a role_binding table.
type RoleBinding struct {
	RoleKey   int64  `gorm:"column:role_key"`
	User      string `gorm:"column:user"`
	ParentKey int64  `gorm:"column:parent_key"`
}

// TableName returns a table name.
func (RoleBinding) TableName() string {
	return "iam_role_binding"
}
