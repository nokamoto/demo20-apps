package compute

import (
	"context"
	"testing"

	"github.com/nokamoto/demo20-apis/cloud/compute/v1alpha"
	"google.golang.org/grpc/codes"

	"github.com/golang/mock/gomock"
	"github.com/nokamoto/demo20-apps/internal/service/compute/mock"
	"github.com/nokamoto/demo20-apps/internal/test"
	"go.uber.org/zap/zaptest"
)

type testCase struct {
	name  string
	run   func(*testing.T, *service) error
	mock  func(*mock.Mockcompute)
	check test.Check
}

type testCases []testCase

func (xs testCases) run(t *testing.T) {
	for _, x := range xs {
		t.Run(x.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			if x.mock == nil {
				t.Fatal("no mock")
			}

			c := mock.NewMockcompute(ctrl)
			x.mock(c)

			if x.run == nil {
				t.Fatal("no run")
			}

			err := x.run(t, &service{
				compute: c,
				logger:  zaptest.NewLogger(t),
			})

			if x.check == nil {
				t.Fatal("no check")
			}

			x.check(t, err)
		})
	}
}

func Test_service_CreateInstance(t *testing.T) {
	run := func(req *v1alpha.CreateInstanceRequest, expected *v1alpha.Instance) func(*testing.T, *service) error {
		return func(t *testing.T, s *service) error {
			return test.Diff1IgnoreUnexported(s.CreateInstance(context.Background(), req))(t, expected)
		}
	}

	xs := testCases{
		{
			name: "OK",
			run: run(
				&v1alpha.CreateInstanceRequest{
					Instance: &v1alpha.Instance{
						Labels: []string{"baz", "qux"},
					},
				},
				&v1alpha.Instance{
					Name:   "instances/todo-xyz",
					Parent: "projects/todo",
					Labels: []string{"baz", "qux"},
				},
			),
			mock: func(c *mock.Mockcompute) {
				gomock.InOrder(
					c.EXPECT().RandomName("todo").Return("todo-xyz"),
					c.EXPECT().Create(
						"todo-xyz",
						"todo",
						&v1alpha.Instance{
							Labels: []string{"baz", "qux"},
						},
					).Return(&v1alpha.Instance{
						Name:   "instances/todo-xyz",
						Parent: "projects/todo",
						Labels: []string{"baz", "qux"},
					}, nil),
				)
			},
			check: test.Code(codes.OK),
		},
	}

	xs.run(t)
}
