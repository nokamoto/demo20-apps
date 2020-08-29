package iam

import (
	"testing"

	"github.com/nokamoto/demo20-apis/cloud/iam/v1alpha"
	"github.com/nokamoto/demo20-apps/internal/mysql/iam"

	"github.com/DATA-DOG/go-sqlmock"
	gomock "github.com/golang/mock/gomock"
	"github.com/jinzhu/gorm"
	"github.com/nokamoto/demo20-apps/internal/test"
)

type testCase struct {
	name  string
	run   func(*testing.T, Iam) error
	mock  func(*MockpermissionQuery)
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

				p := NewMockpermissionQuery(ctrl)
				x.mock(p)

				err := x.run(t, Iam{
					permissionQuery: p,
					db:              g,
				})

				x.check(t, err)
			})
		})
	}
}

func TestIam_Create(t *testing.T) {
	run := func(id string, expected *v1alpha.Permission) func(*testing.T, Iam) error {
		return func(t *testing.T, i Iam) error {
			return test.Diff1IgnoreUnexported(i.Create(id))(t, expected)
		}
	}

	xs := testCases{
		{
			name: "OK",
			run: run("foo", &v1alpha.Permission{
				Name: "permissions/foo",
			}),
			mock: func(p *MockpermissionQuery) {
				gomock.InOrder(
					p.EXPECT().Create(gomock.Any(), &iam.Permission{
						PermissionID: "foo",
					}).Return(nil),
				)
			},
			check: test.Succeeded,
			tx:    test.Commit,
		},
	}

	xs.run(t)
}
