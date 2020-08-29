package test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"

	"google.golang.org/protobuf/testing/protocmp"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/go-cmp/cmp"
)

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

// Code asserts that err has the expected code.
func Code(expected codes.Code) Check {
	return func(t *testing.T, err error) {
		t.Helper()
		if actual := status.Code(err); actual != expected {
			t.Errorf("expected %v but actual %v", expected, actual)
		}
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

// Diff1IgnoreUnexported asserts that a1 is equal to e1 with protocmp.Transform.
func Diff1IgnoreUnexported(a1 interface{}, err error) func(*testing.T, interface{}) error {
	return func(t *testing.T, e1 interface{}) error {
		if diff := cmp.Diff(e1, a1, protocmp.Transform()); len(diff) != 0 {
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

// UseGorm runs f with mocked sqls.
func UseGorm(t *testing.T, f func(sqlmock.Sqlmock, *gorm.DB)) {
	db, m, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	g, err := gorm.Open("mysql", db)
	if err != nil {
		t.Fatal(err)
	}
	defer g.Close()

	f(m, g)
}

// Commit asserts that the transaction is comitted.
func Commit(mock sqlmock.Sqlmock) {
	mock.ExpectBegin()
	mock.ExpectCommit()
}
