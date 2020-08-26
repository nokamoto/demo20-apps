package core

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
)

// TestCase represents a single test case using sql-mock.
type TestCase struct {
	Name  string
	Run   Run
	Mock  func(sqlmock.Sqlmock)
	Check Check
}

// Run represents a test execution.
type Run func(*testing.T, *gorm.DB) error

// Check represents an assertion of TestCase.
type Check func(*testing.T, error)

// Succeeded asserts that err is nil.
func Succeeded(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("expected no error: %v", err)
	}
}

// Failed asserts that err is expected.
func Failed(expected error) Check {
	return func(t *testing.T, actual error) {
		t.Helper()
		if !errors.Is(actual, expected) {
			t.Errorf("expected %v but actual %v", expected, actual)
		}
	}
}

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

			if x.Mock != nil {
				x.Mock(mock)
			}

			g, err := gorm.Open("mysql", db)
			if err != nil {
				t.Fatal(err)
			}
			defer g.Close()

			if x.Run == nil {
				t.Fatal("no execution")
			}

			err = g.Transaction(func(tx *gorm.DB) error {
				return x.Run(t, tx)
			})

			if x.Check == nil {
				t.Fatal("no check")
			}

			x.Check(t, err)

			err = mock.ExpectationsWereMet()
			if err != nil {
				t.Error(err)
			}
		})
	}
}

// Diff1 asserts that a1 is equal to e1.
func Diff1(a1 interface{}, err error) func(*testing.T, interface{}) error {
	return func(t *testing.T, e1 interface{}) error {
		if diff := cmp.Diff(e1, a1); len(diff) != 0 {
			t.Error(diff)
		}
		return err
	}
}

// Diff2 asserts that a1 is equal to e1, and so on.
func Diff2(a1, a2 interface{}, err error) func(*testing.T, interface{}, interface{}) error {
	return func(t *testing.T, e1, e2 interface{}) error {
		if diff := cmp.Diff(e1, a1); len(diff) != 0 {
			t.Error(diff)
		}
		if diff := cmp.Diff(e2, a2); len(diff) != 0 {
			t.Error(diff)
		}
		return err
	}
}
