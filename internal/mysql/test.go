package mysql

import (
	"testing"

	"github.com/nokamoto/demo20-apps/internal/test"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
)

// TestCase represents a single test case using sql-mock.
type TestCase struct {
	Name  string
	Run   Run
	Mock  func(sqlmock.Sqlmock)
	Check test.Check
}

// Run represents a test execution.
type Run func(*testing.T, *gorm.DB) error

// TestCases represents a list of TestCase.
type TestCases []TestCase

// Run runs all test cases.
func (xs TestCases) Run(t *testing.T) {
	for _, x := range xs {
		t.Run(x.Name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatal(err)
			}
			defer db.Close()

			x.Mock(mock)

			g, err := gorm.Open("mysql", db)
			if err != nil {
				t.Fatal(err)
			}
			defer g.Close()

			err = g.Transaction(func(tx *gorm.DB) error {
				return x.Run(t, tx)
			})

			x.Check(t, err)

			err = mock.ExpectationsWereMet()
			if err != nil {
				t.Error(err)
			}
		})
	}
}
