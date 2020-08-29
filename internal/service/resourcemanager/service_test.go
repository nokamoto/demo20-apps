package resourcemanager

import (
	"context"
	"testing"

	"github.com/nokamoto/demo20-apis/cloud/resourcemanager/v1alpha"

	"github.com/golang/mock/gomock"
	"github.com/nokamoto/demo20-apps/internal/test"
	"go.uber.org/zap/zaptest"
)

type testCase struct {
	name  string
	run   func(*testing.T, *service) error
	mock  func(*Mockresourcemanager)
	check test.Check
}

type testCases []testCase

func (xs testCases) run(t *testing.T) {
	for _, x := range xs {
		t.Run(x.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			r := NewMockresourcemanager(ctrl)
			x.mock(r)

			err := x.run(t, &service{
				resourcemanager: r,
				logger:          zaptest.NewLogger(t),
			})

			x.check(t, err)
		})
	}
}

func Test_service_GetProject(t *testing.T) {
	run := func(req *v1alpha.GetProjectRequest, expected *v1alpha.Project) func(*testing.T, *service) error {
		return func(t *testing.T, s *service) error {
			return test.Diff1IgnoreUnexported(s.GetProject(context.Background(), req))(t, expected)
		}
	}

	xs := testCases{
		{
			name: "OK",
			run: run(
				&v1alpha.GetProjectRequest{
					Name: "projects/foo",
				},
				&v1alpha.Project{
					Name:        "projects/foo",
					DisplayName: "foo display name",
				},
			),
			mock: func(r *Mockresourcemanager) {
				r.EXPECT().Get("foo").Return(&v1alpha.Project{
					Name:        "projects/foo",
					DisplayName: "foo display name",
				}, nil)
			},
			check: test.Succeeded,
		},
	}

	xs.run(t)
}
