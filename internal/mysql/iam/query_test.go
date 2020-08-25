package iam

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/go-cmp/cmp"
	"github.com/jinzhu/gorm"
	"github.com/nokamoto/demo20-apps/internal/mysql/core"
)

func TestPermissionQuery_Create(t *testing.T) {
	run := func(permission Permission) core.Run {
		return func(t *testing.T, tx *gorm.DB) error {
			return PermissionQuery{}.Create(tx, &permission)
		}
	}

	permission := Permission{
		PermissionID: "foo",
	}

	xs := core.TestCases{
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
			Check: core.Succeeded,
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

	run := func(permissionIDs []string, expected []*Permission) core.Run {
		return func(t *testing.T, tx *gorm.DB) error {
			actual, err := PermissionQuery{}.List(tx, permissionIDs...)
			if diff := cmp.Diff(expected, actual); len(diff) != 0 {
				t.Error(diff)
			}
			return err
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

	xs := core.TestCases{
		{
			Name: "OK",
			Run:  run([]string{"foo", "bar"}, []*Permission{&foo, &bar}),
			Mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `iam_permission`  WHERE (permission_id in (?,?))")).
					WithArgs("foo", "bar").
					WillReturnRows(rows(foo, bar))
			},
			Check: core.Succeeded,
		},
	}

	xs.Run(t)
}

func TestRoleQuery_Create(t *testing.T) {
	run := func(role Role, permissions ...*Permission) core.Run {
		return func(t *testing.T, tx *gorm.DB) error {
			return RoleQuery{}.Create(tx, &role, permissions...)
		}
	}

	role := Role{
		RoleID:      "foo",
		ParentKey:   100,
		DisplayName: "foo display name",
	}

	permission := Permission{
		PermissionKey: 200,
		PermissionID:  "bar",
	}

	xs := core.TestCases{
		{
			Name: "OK",
			Run:  run(role, &permission),
			Mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `iam_role` (`role_id`,`parent_key`,`display_name`) VALUES (?,?,?)")).
					WithArgs(role.RoleID, role.ParentKey, role.DisplayName).
					WillReturnResult(sqlmock.NewResult(1000, 1))
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `iam_role_permission` (`permission_key`,`role_key`) VALUES(?,?)")).
					WithArgs(role.RoleKey, permission.PermissionKey).
					WillReturnResult(sqlmock.NewResult(2000, 1))
				mock.ExpectCommit()
			},
			Check: core.Succeeded,
		},
	}

	xs.Run(t)
}
