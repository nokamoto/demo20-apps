package resourcemanager

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/jinzhu/gorm"
	"github.com/nokamoto/demo20-apis/cloud/resourcemanager/v1alpha"
	"github.com/nokamoto/demo20-apps/internal/mysql/resourcemanager"
	"github.com/nokamoto/demo20-apps/internal/test"
)

type testCase struct {
	name  string
	run   func(*testing.T, ResourceManager) error
	mock  func(*MockprojectQuery)
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

				p := NewMockprojectQuery(ctrl)
				x.mock(p)

				err := x.run(t, ResourceManager{
					projectQuery: p,
					db:           g,
				})

				x.check(t, err)
			})
		})
	}
}

func TestResourceManager_Get(t *testing.T) {
	run := func(id string, expected *v1alpha.Project) func(*testing.T, ResourceManager) error {
		return func(t *testing.T, r ResourceManager) error {
			return test.Diff1IgnoreUnexported(r.Get(id))(t, expected)
		}
	}

	xs := testCases{
		{
			name: "OK",
			run: run("foo", &v1alpha.Project{
				Name:        "projects/foo",
				DisplayName: "foo display name",
			}),
			mock: func(p *MockprojectQuery) {
				p.EXPECT().Get(gomock.Any(), "foo").Return(&resourcemanager.Project{
					ProjectID:   "foo",
					DisplayName: "foo display name",
				}, nil)
			},
			check: test.Succeeded,
			tx:    test.Commit,
		},
	}

	xs.run(t)
}
