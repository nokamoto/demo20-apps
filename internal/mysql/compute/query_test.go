package compute

import (
	"regexp"
	"testing"

	"github.com/nokamoto/demo20-apps/internal/test"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/nokamoto/demo20-apps/internal/mysql"
)

func rows(xs ...Instance) *sqlmock.Rows {
	v := sqlmock.NewRows([]string{
		"instance_key", "instance_id", "parent_key", "labels",
	})
	for _, x := range xs {
		v.AddRow(x.InstanceKey, x.InstanceID, x.ParentKey, x.Labels)
	}
	return v
}

func TestQuery_Create(t *testing.T) {
	run := func(instance Instance) mysql.Run {
		return func(t *testing.T, tx *gorm.DB) error {
			return Query{}.Create(tx, &instance)
		}
	}

	instance := Instance{
		InstanceID: "foo",
		ParentKey:  100,
		Labels:     "bar,baz",
	}

	xs := mysql.TestCases{
		{
			Name: "OK",
			Run:  run(instance),
			Mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `compute_instance` (`instance_id`,`parent_key`,`labels`) VALUES (?,?,?)")).
					WithArgs(instance.InstanceID, instance.ParentKey, instance.Labels).
					WillReturnResult(sqlmock.NewResult(1000, 1))
				mock.ExpectCommit()
			},
			Check: test.Succeeded,
		},
		{
			Name: "AlreadyExists",
			Run:  run(instance),
			Mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `compute_instance` (`instance_id`,`parent_key`,`labels`) VALUES (?,?,?)")).
					WithArgs(instance.InstanceID, instance.ParentKey, instance.Labels).
					WillReturnError(&mysql.TestDupEntryErr)
				mock.ExpectRollback()
			},
			Check: test.Failed(mysql.ErrAlreadyExists),
		},
	}

	xs.Run(t)
}

func TestQuery_Delete(t *testing.T) {
	run := func(id string) mysql.Run {
		return func(t *testing.T, tx *gorm.DB) error {
			return Query{}.Delete(tx, id)
		}
	}

	xs := mysql.TestCases{
		{
			Name: "OK",
			Run:  run("foo"),
			Mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `compute_instance` WHERE (instance_id = ?)")).
					WithArgs("foo").
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			Check: test.Succeeded,
		},
	}

	xs.Run(t)
}

func TestQuery_Get(t *testing.T) {
	run := func(id string, expected *Instance) mysql.Run {
		return func(t *testing.T, tx *gorm.DB) error {
			return test.Diff1(Query{}.Get(tx, id))(t, expected)
		}
	}

	instance := Instance{
		InstanceKey: 1000,
		InstanceID:  "foo",
		ParentKey:   100,
		Labels:      "bar,baz",
	}

	xs := mysql.TestCases{
		{
			Name: "OK",
			Run:  run("foo", &instance),
			Mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `compute_instance` WHERE (instance_id = ?) LIMIT 1")).
					WithArgs("foo").
					WillReturnRows(rows(instance))
				mock.ExpectCommit()
			},
			Check: test.Succeeded,
		},
		{
			Name: "NotFound",
			Run:  run("foo", &Instance{}),
			Mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `compute_instance` WHERE (instance_id = ?) LIMIT 1")).
					WithArgs("foo").
					WillReturnRows(rows())
				mock.ExpectRollback()
			},
			Check: test.Failed(mysql.ErrNotFound),
		},
	}

	xs.Run(t)
}

func TestQuery_List(t *testing.T) {
	offset, limit := 100, 200

	run := func(parentKey int64, expected []*Instance) mysql.Run {
		return func(t *testing.T, tx *gorm.DB) error {
			return test.Diff1(Query{}.List(tx, parentKey, offset, limit))(t, expected)
		}
	}

	foo := Instance{
		InstanceKey: 1000,
		InstanceID:  "foo",
		ParentKey:   100,
		Labels:      "bar,baz",
	}

	bar := Instance{
		InstanceKey: 2000,
		InstanceID:  "bar",
		ParentKey:   200,
		Labels:      "foo,baz",
	}

	xs := mysql.TestCases{
		{
			Name: "OK empty",
			Run:  run(3000, []*Instance{}),
			Mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `compute_instance`  WHERE (parent_key = ?) LIMIT 200 OFFSET 100")).
					WithArgs(3000).
					WillReturnRows(rows())
				mock.ExpectCommit()
			},
			Check: test.Succeeded,
		},
		{
			Name: "OK",
			Run:  run(3000, []*Instance{&foo, &bar}),
			Mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `compute_instance`  WHERE (parent_key = ?) LIMIT 200 OFFSET 100")).
					WithArgs(3000).
					WillReturnRows(rows(foo, bar))
				mock.ExpectCommit()
			},
			Check: test.Succeeded,
		},
	}

	xs.Run(t)
}
