package compute

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/jinzhu/gorm"
	"github.com/nokamoto/demo20-apis/cloud/compute/v1alpha"
	"github.com/nokamoto/demo20-apps/internal/compute/mock"
	"github.com/nokamoto/demo20-apps/internal/mysql/compute"
	"github.com/nokamoto/demo20-apps/internal/mysql/resourcemanager"
	"github.com/nokamoto/demo20-apps/internal/test"
)

type testCase struct {
	name  string
	run   func(*testing.T, Compute) error
	mock  func(*mock.MockinstanceQuery, *mock.MockprojectQuery)
	check test.Check
	tx    func(sqlmock.Sqlmock)
}

type testCases []testCase

func (xs testCases) run(t *testing.T) {
	for _, x := range xs {
		t.Run(x.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			db, m, err := sqlmock.New()
			if err != nil {
				t.Fatal(err)
			}
			defer db.Close()

			if x.tx == nil {
				t.Fatal("no tx")
			}

			x.tx(m)

			g, err := gorm.Open("mysql", db)
			if err != nil {
				t.Fatal(err)
			}
			defer g.Close()

			if x.mock == nil {
				t.Fatal("no mock")
			}

			i := mock.NewMockinstanceQuery(ctrl)
			p := mock.NewMockprojectQuery(ctrl)
			x.mock(i, p)

			if x.run == nil {
				t.Fatal("no run")
			}

			err = x.run(t, Compute{
				instanceQuery: i,
				projectQuery:  p,
				db:            g,
			})

			if x.check == nil {
				t.Fatal("no check")
			}

			x.check(t, err)
		})
	}
}

func commit(mock sqlmock.Sqlmock) {
	mock.ExpectBegin()
	mock.ExpectCommit()
}

func TestCompute_Create(t *testing.T) {
	run := func(id, projectID string, req, expected *v1alpha.Instance) func(*testing.T, Compute) error {
		return func(t *testing.T, c Compute) error {
			return test.Diff1IgnoreUnexported(c.Create(id, projectID, req))(t, expected, v1alpha.Instance{})
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
			mock: func(i *mock.MockinstanceQuery, p *mock.MockprojectQuery) {
				gomock.InOrder(
					p.EXPECT().Get(gomock.Any(), "bar").Return(&resourcemanager.Project{
						ProjectKey: 100,
					}, nil),
					i.EXPECT().Create(gomock.Any(), &compute.Instance{
						InstanceID: "foo",
						ParentKey:  100,
						Labels:     "baz,qux",
					}).Return(nil),
				)
			},
			check: test.Succeeded,
			tx:    commit,
		},
	}

	xs.run(t)
}
