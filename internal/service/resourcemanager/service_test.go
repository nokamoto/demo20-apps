package resourcemanager

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/nokamoto/demo20-apis/cloud/api"
	"github.com/nokamoto/demo20-apis/cloud/resourcemanager/v1alpha"
	"github.com/nokamoto/demo20-apps/internal/test"
	"github.com/nokamoto/demo20-apps/pkg/sdk/metadata"
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
	run := func(ctx context.Context, req *v1alpha.GetProjectRequest, expected *v1alpha.Project) func(*testing.T, *service) error {
		return func(t *testing.T, s *service) error {
			return test.Diff1IgnoreUnexported(s.GetProject(ctx, req))(t, expected)
		}
	}

	xs := testCases{
		{
			name: "OK",
			run: run(
				metadata.NewIncomingContextF(context.Background(), &api.Metadata{
					Parent: "projects/foo",
				}),
				&v1alpha.GetProjectRequest{},
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

func Test_service_CreateProject(t *testing.T) {
	run := func(req *v1alpha.CreateProjectRequest, expected *v1alpha.Project) func(*testing.T, *service) error {
		return func(t *testing.T, s *service) error {
			return test.Diff1IgnoreUnexported(s.CreateProject(context.Background(), req))(t, expected)
		}
	}

	xs := testCases{
		{
			name: "OK",
			run: run(
				&v1alpha.CreateProjectRequest{
					ProjectId: "foo",
					Project: &v1alpha.Project{
						DisplayName: "foo display name",
					},
				},
				&v1alpha.Project{
					Name:        "projects/foo",
					DisplayName: "foo display name",
				},
			),
			mock: func(r *Mockresourcemanager) {
				r.EXPECT().Create("foo", test.ProtoEq(&v1alpha.Project{
					DisplayName: "foo display name",
				})).Return(&v1alpha.Project{
					Name:        "projects/foo",
					DisplayName: "foo display name",
				}, nil)
			},
			check: test.Succeeded,
		},
	}

	xs.run(t)
}
