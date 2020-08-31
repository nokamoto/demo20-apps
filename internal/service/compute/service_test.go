package compute

import (
	"context"
	"testing"

	"github.com/nokamoto/demo20-apis/cloud/api"
	"github.com/nokamoto/demo20-apps/pkg/sdk/metadata"

	"github.com/nokamoto/demo20-apis/cloud/compute/v1alpha"
	"google.golang.org/grpc/codes"

	"github.com/golang/mock/gomock"
	"github.com/nokamoto/demo20-apps/internal/test"
	"go.uber.org/zap/zaptest"
)

type testCase struct {
	name  string
	run   func(*testing.T, *service) error
	mock  func(*Mockcompute)
	check test.Check
}

type testCases []testCase

func (xs testCases) run(t *testing.T) {
	for _, x := range xs {
		t.Run(x.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			c := NewMockcompute(ctrl)
			x.mock(c)

			err := x.run(t, &service{
				compute: c,
				logger:  zaptest.NewLogger(t),
			})

			x.check(t, err)
		})
	}
}

func Test_service_CreateInstance(t *testing.T) {
	run := func(ctx context.Context, req *v1alpha.CreateInstanceRequest, expected *v1alpha.Instance) func(*testing.T, *service) error {
		return func(t *testing.T, s *service) error {
			return test.Diff1IgnoreUnexported(s.CreateInstance(ctx, req))(t, expected)
		}
	}

	xs := testCases{
		{
			name: "OK",
			run: run(
				metadata.NewIncomingContextF(context.Background(), &api.Metadata{
					Parent: "projects/test",
				}),
				&v1alpha.CreateInstanceRequest{
					Instance: &v1alpha.Instance{
						Labels: []string{"baz", "qux"},
					},
				},
				&v1alpha.Instance{
					Name:   "instances/test-xyz",
					Parent: "projects/test",
					Labels: []string{"baz", "qux"},
				},
			),
			mock: func(c *Mockcompute) {
				gomock.InOrder(
					c.EXPECT().RandomName("test").Return("test-xyz"),
					c.EXPECT().Create(
						"test-xyz",
						"test",
						&v1alpha.Instance{
							Labels: []string{"baz", "qux"},
						},
					).Return(&v1alpha.Instance{
						Name:   "instances/test-xyz",
						Parent: "projects/test",
						Labels: []string{"baz", "qux"},
					}, nil),
				)
			},
			check: test.Code(codes.OK),
		},
	}

	xs.run(t)
}
