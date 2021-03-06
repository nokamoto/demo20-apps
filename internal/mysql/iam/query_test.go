package iam

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/nokamoto/demo20-apps/internal/mysql"
	"github.com/nokamoto/demo20-apps/internal/test"
)

func TestPermissionQuery_Create(t *testing.T) {
	run := func(permission Permission) mysql.Run {
		return func(t *testing.T, tx *gorm.DB) error {
			return PermissionQuery{}.Create(tx, &permission)
		}
	}

	permission := Permission{
		PermissionID: "foo",
	}

	xs := mysql.TestCases{
		{
			Name: "OK",
			Run:  run(permission),
			Mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `iam_permission` (`permission_id`) VALUES (?)")).
					WithArgs(permission.PermissionID).
					WillReturnResult(sqlmock.NewResult(1000, 1))
				mock.ExpectCommit()
			},
			Check: test.Succeeded,
		},
	}

	xs.Run(t)
}

func TestPermissionQuery_List(t *testing.T) {
	rows := func(xs ...Permission) *sqlmock.Rows {
		v := sqlmock.NewRows([]string{
			"permission_key", "permission_id",
		})
		for _, x := range xs {
			v.AddRow(x.PermissionKey, x.PermissionID)
		}
		return v
	}

	run := func(permissionIDs []string, expected []*Permission) mysql.Run {
		return func(t *testing.T, tx *gorm.DB) error {
			return test.Diff1(PermissionQuery{}.List(tx, permissionIDs...))(t, expected)
		}
	}

	foo := Permission{
		PermissionKey: 100,
		PermissionID:  "foo",
	}

	bar := Permission{
		PermissionKey: 200,
		PermissionID:  "bar",
	}

	xs := mysql.TestCases{
		{
			Name: "OK",
			Run:  run([]string{"foo", "bar"}, []*Permission{&foo, &bar}),
			Mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `iam_permission`  WHERE (permission_id in (?,?))")).
					WithArgs("foo", "bar").
					WillReturnRows(rows(foo, bar))
				mock.ExpectCommit()
			},
			Check: test.Succeeded,
		},
	}

	xs.Run(t)
}

func TestRoleQuery_Create(t *testing.T) {
	run := func(role Role, permissions ...*Permission) mysql.Run {
		return func(t *testing.T, tx *gorm.DB) error {
			return RoleQuery{}.Create(tx, &role, permissions...)
		}
	}

	role := Role{
		RoleID:      "foo",
		ParentID:    "baz",
		DisplayName: "foo display name",
	}

	permission := Permission{
		PermissionKey: 200,
		PermissionID:  "bar",
	}

	xs := mysql.TestCases{
		{
			Name: "OK",
			Run:  run(role, &permission),
			Mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `iam_role` (`role_id`,`parent_id`,`display_name`) VALUES (?,?,?)")).
					WithArgs(role.RoleID, role.ParentID, role.DisplayName).
					WillReturnResult(sqlmock.NewResult(1000, 1))
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `iam_role_permission` (`role_key`,`permission_key`) VALUES (?,?)")).
					WithArgs(1000, permission.PermissionKey).
					WillReturnResult(sqlmock.NewResult(2000, 1))
				mock.ExpectCommit()
			},
			Check: test.Succeeded,
		},
	}

	xs.Run(t)
}

func TestRoleQuery_Delete(t *testing.T) {
	run := func(role Role) mysql.Run {
		return func(t *testing.T, tx *gorm.DB) error {
			return RoleQuery{}.Delete(tx, &role)
		}
	}

	role := Role{
		RoleKey: 100,
	}

	xs := mysql.TestCases{
		{
			Name: "OK",
			Run:  run(role),
			Mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `iam_role_permission` WHERE (role_key = ?)")).
					WithArgs(role.RoleKey).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `iam_role` WHERE (role_key = ?)")).
					WithArgs(role.RoleKey).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			Check: test.Succeeded,
		},
	}

	xs.Run(t)
}

func roleRows(roles ...Role) *sqlmock.Rows {
	v := sqlmock.NewRows([]string{
		"role_key", "role_id", "parent_id", "display_name",
	})
	for _, x := range roles {
		v.AddRow(x.RoleKey, x.RoleID, x.ParentID, x.DisplayName)
	}
	return v
}

func rolePermissionRows(permissions ...RolePermission) *sqlmock.Rows {
	v := sqlmock.NewRows([]string{
		"role_key", "permission_key",
	})
	for _, x := range permissions {
		v.AddRow(x.RoleKey, x.PermissionKey)
	}
	return v
}

func TestRoleQuery_Get(t *testing.T) {
	run := func(id string, rexpected *Role, pexpected ...*RolePermission) mysql.Run {
		return func(t *testing.T, tx *gorm.DB) error {
			return test.Diff2(RoleQuery{}.Get(tx, id))(t, rexpected, pexpected)
		}
	}

	role := Role{
		RoleKey: 100,
	}

	permission := RolePermission{
		RoleKey:       100,
		PermissionKey: 200,
	}

	xs := mysql.TestCases{
		{
			Name: "OK",
			Run:  run("foo", &role, &permission),
			Mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `iam_role` WHERE (role_id = ?) LIMIT 1")).
					WithArgs("foo").
					WillReturnRows(roleRows(role))
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `iam_role_permission` WHERE (role_key = ?)")).
					WithArgs(role.RoleKey).
					WillReturnRows(rolePermissionRows(permission))
				mock.ExpectCommit()
			},
			Check: test.Succeeded,
		},
	}

	xs.Run(t)
}

func TestRoleQuery_Update(t *testing.T) {
	run := func(role *Role, permissions ...*Permission) mysql.Run {
		return func(t *testing.T, tx *gorm.DB) error {
			return RoleQuery{}.Update(tx, role, permissions...)
		}
	}

	role := Role{
		RoleKey:     100,
		RoleID:      "foo",
		ParentID:    "bar",
		DisplayName: "foo display name",
	}

	permission := Permission{
		PermissionKey: 200,
	}

	xs := mysql.TestCases{
		{
			Name: "OK",
			Run:  run(&role, &permission),
			Mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("UPDATE `iam_role` SET `role_id` = ?, `parent_id` = ?, `display_name` = ? WHERE `iam_role`.`role_key` = ?")).
					WithArgs(role.RoleID, role.ParentID, role.DisplayName, role.RoleKey).
					WillReturnResult(sqlmock.NewResult(1000, 1))
				mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `iam_role_permission` WHERE (role_key = ?)")).
					WithArgs(role.RoleKey).
					WillReturnResult(sqlmock.NewResult(2000, 1))
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `iam_role_permission` (`role_key`,`permission_key`) VALUES (?,?)")).
					WithArgs(role.RoleKey, permission.PermissionKey).
					WillReturnResult(sqlmock.NewResult(3000, 1))
				mock.ExpectCommit()
			},
			Check: test.Succeeded,
		},
	}

	xs.Run(t)
}

func TestRoleQuery_List(t *testing.T) {
	offset, limit := 100, 200

	run := func(parentID string, rexpected []*Role, pexpected []*RolePermission) mysql.Run {
		return func(t *testing.T, tx *gorm.DB) error {
			return test.Diff2(RoleQuery{}.List(tx, parentID, offset, limit))(t, rexpected, pexpected)
		}
	}

	role := Role{
		RoleKey: 300,
	}

	permission := RolePermission{
		RoleKey:       300,
		PermissionKey: 400,
	}

	xs := mysql.TestCases{
		{
			Name: "OK",
			Run:  run("foo", []*Role{&role}, []*RolePermission{&permission}),
			Mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `iam_role` WHERE (parent_id = ?) LIMIT 200 OFFSET 100")).
					WithArgs("foo").
					WillReturnRows(roleRows(role))
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `iam_role_permission` WHERE (role_key in (?))")).
					WithArgs(role.RoleKey).
					WillReturnRows(rolePermissionRows(permission))
				mock.ExpectCommit()
			},
			Check: test.Succeeded,
		},
	}

	xs.Run(t)
}
