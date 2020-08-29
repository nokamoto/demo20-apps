package iam

import (
	"testing"

	admin "github.com/nokamoto/demo20-apis/cloud/iam/admin/v1alpha"
	"github.com/nokamoto/demo20-apis/cloud/iam/v1alpha"

	"github.com/golang/mock/gomock"
	"github.com/nokamoto/demo20-apps/internal/test"
	"go.uber.org/zap/zaptest"
)

type testCase struct {
	name  string
	run   func(*testing.T, *service) error
	mock  func(*Mockiam)
	check test.Check
}

type testCases []testCase

func (xs testCases) run(t *testing.T) {
	for _, x := range xs {
		t.Run(x.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			i := NewMockiam(ctrl)
			x.mock(i)

			err := x.run(t, &service{
				iam:    i,
				logger: zaptest.NewLogger(t),
			})

			x.check(t, err)
		})
	}
}

func Test_service_CreatePermission(t *testing.T) {
	run := func(req *admin.CreatePermissionRequest, expected *v1alpha.Permission) func(*testing.T, *service) error {
		return func(t *testing.T, s *service) error {
			return test.Diff1IgnoreUnexported(s.CreatePermission(nil, req))(t, expected)
		}
	}

	xs := testCases{
		{
			name: "OK",
			run: run(
				&admin.CreatePermissionRequest{
					PermissionId: "foo",
				},
				&v1alpha.Permission{
					Name: "permissions/foo",
				},
			),
			mock: func(i *Mockiam) {
				i.EXPECT().Create("foo").Return(&v1alpha.Permission{
					Name: "permissions/foo",
				}, nil)
			},
			check: test.Succeeded,
		},
	}

	xs.run(t)
}
