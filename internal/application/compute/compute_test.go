package compute

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/jinzhu/gorm"
	"github.com/nokamoto/demo20-apis/cloud/compute/v1alpha"
	"github.com/nokamoto/demo20-apps/internal/mysql/compute"
	"github.com/nokamoto/demo20-apps/internal/test"
)

type testCase struct {
	name  string
	run   func(*testing.T, Compute) error
	mock  func(*MockinstanceQuery)
	check test.Check
	tx    func(sqlmock.Sqlmock)
}

type testCases []testCase

func (xs testCases) run(t *testing.T) {
	for _, x := range xs {
		t.Run(x.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			test.UseGorm(t, func(m sqlmock.Sqlmock, g *gorm.DB) {
				x.tx(m)

				i := NewMockinstanceQuery(ctrl)
				x.mock(i)

				err := x.run(t, Compute{
					instanceQuery: i,
					db:            g,
				})

				x.check(t, err)
			})
		})
	}
}

func TestCompute_Create(t *testing.T) {
	run := func(id, projectID string, req, expected *v1alpha.Instance) func(*testing.T, Compute) error {
		return func(t *testing.T, c Compute) error {
			return test.Diff1IgnoreUnexported(c.Create(id, projectID, req))(t, expected)
		}
	}

	xs := testCases{
		{
			name: "OK",
			run: run(
				"foo",
				"bar",
				&v1alpha.Instance{
					Labels: []string{"baz", "qux"},
				},
				&v1alpha.Instance{
					Name:   "instances/foo",
					Parent: "projects/bar",
					Labels: []string{"baz", "qux"},
				},
			),
			mock: func(i *MockinstanceQuery) {
				i.EXPECT().Create(gomock.Any(), &compute.Instance{
					InstanceID: "foo",
					ParentID:   "bar",
					Labels:     "baz,qux",
				}).Return(nil)
			},
			check: test.Succeeded,
			tx:    test.Commit,
		},
	}

	xs.run(t)
}
