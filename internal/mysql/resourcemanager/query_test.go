package resourcemanager

import (
	"regexp"
	"testing"

	"github.com/nokamoto/demo20-apps/internal/test"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/nokamoto/demo20-apps/internal/mysql"
)

func TestQuery_Create(t *testing.T) {
	run := func(project Project) mysql.Run {
		return func(t *testing.T, tx *gorm.DB) error {
			return Query{}.Create(tx, &project)
		}
	}

	project := Project{
		ProjectID:   "foo",
		DisplayName: "foo display name",
	}

	xs := mysql.TestCases{
		{
			Name: "OK",
			Run:  run(project),
			Mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `resourcemanager_project` (`project_id`,`display_name`) VALUES (?,?)")).
					WithArgs(project.ProjectID, project.DisplayName).
					WillReturnResult(sqlmock.NewResult(1000, 1))
				mock.ExpectCommit()
			},
			Check: test.Succeeded,
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
				mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `resourcemanager_project`  WHERE (project_id = ?)")).
					WithArgs("foo").
					WillReturnResult(sqlmock.NewResult(1000, 1))
				mock.ExpectCommit()
			},
			Check: test.Succeeded,
		},
	}

	xs.Run(t)
}

func TestQuery_Update(t *testing.T) {
	run := func(project Project) mysql.Run {
		return func(t *testing.T, tx *gorm.DB) error {
			return Query{}.Update(tx, &project)
		}
	}

	project := Project{
		ProjectKey:  100,
		ProjectID:   "foo",
		DisplayName: "foo display name",
	}

	xs := mysql.TestCases{
		{
			Name: "OK",
			Run:  run(project),
			Mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("UPDATE `resourcemanager_project` SET `project_id` = ?, `display_name` = ? WHERE `resourcemanager_project`.`project_key` = ?")).
					WithArgs(project.ProjectID, project.DisplayName, project.ProjectKey).
					WillReturnResult(sqlmock.NewResult(1000, 1))
				mock.ExpectCommit()
			},
			Check: test.Succeeded,
		},
	}

	xs.Run(t)
}

func TestQuery_Get(t *testing.T) {
	rows := func(xs ...Project) *sqlmock.Rows {
		v := sqlmock.NewRows([]string{
			"project_key", "project_id", "display_name",
		})
		for _, x := range xs {
			v.AddRow(x.ProjectKey, x.ProjectID, x.DisplayName)
		}
		return v
	}

	run := func(id string, expected *Project) mysql.Run {
		return func(t *testing.T, tx *gorm.DB) error {
			return test.Diff1(Query{}.Get(tx, id))(t, expected)
		}
	}

	project := Project{
		ProjectKey:  100,
		ProjectID:   "foo",
		DisplayName: "foo display name",
	}

	xs := mysql.TestCases{
		{
			Name: "OK",
			Run:  run(project.ProjectID, &project),
			Mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta("")).
					WithArgs(project.ProjectID).
					WillReturnRows(rows(project))
				mock.ExpectCommit()
			},
			Check: test.Succeeded,
		},
	}

	xs.Run(t)
}
