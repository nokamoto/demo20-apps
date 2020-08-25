package iam

// Permission represents a permission table.
type Permission struct {
	permissionKey int64  `gorm:"column:permission_key;auto_increment;primary_key"`
	permissionID  string `gorm:"column:permission_id"`
}

// TableName returns a table name.
func (Permission) TableName() string {
	return "iam_permission"
}

// Role represents a role table.
type Role struct {
	roleKey     int64  `gorm:"column:role_key;auto_increment;primary_key"`
	roleID      string `gorm:"column:role_id"`
	parentKey   int64  `gorm:"column:parent_key"`
	displayName string `gorm:"column:display_name"`
}

// TableName returns a table name.
func (Role) TableName() string {
	return "iam_role"
}

// RolePermission represents a role_permission table.
type RolePermission struct {
	rokeKey       int64 `gorm:"column:role_key"`
	permissionKey int64 `gorm:"column:permission_key"`
}

// TableName returns a table name.
func (RolePermission) TableName() string {
	return "iam_role_permission"
}

// RoleBinding represents a role_binding table.
type RoleBinding struct {
	roleKey   int64  `gorm:"column:role_key"`
	user      string `gorm:"column:user"`
	parentKey int64  `gorm:"column:parent_key"`
}

// TableName returns a table name.
func (RoleBinding) TableName() string {
	return "iam_role_binding"
}
