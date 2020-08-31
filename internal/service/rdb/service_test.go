package rdb

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/nokamoto/demo20-apis/cloud/api"
	"github.com/nokamoto/demo20-apis/cloud/rdb/v1alpha"
	"github.com/nokamoto/demo20-apps/internal/test"
	"github.com/nokamoto/demo20-apps/pkg/sdk/metadata"
	"go.uber.org/zap/zaptest"
)

type testCase struct {
	name  string
	run   func(*testing.T, *service) error
	mock  func(*Mockrdb)
	check test.Check
}

type testCases []testCase

func (xs testCases) run(t *testing.T) {
	for _, x := range xs {
		t.Run(x.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			r := NewMockrdb(ctrl)
			x.mock(r)

			err := x.run(t, &service{
				rdb:    r,
				logger: zaptest.NewLogger(t),
			})

			x.check(t, err)
		})
	}
}

func Test_service_CreateCluster(t *testing.T) {
	run := func(ctx context.Context, req *v1alpha.CreateClusterRequest, expected *v1alpha.Cluster) func(*testing.T, *service) error {
		return func(t *testing.T, s *service) error {
			return test.Diff1IgnoreUnexported(s.CreateCluster(ctx, req))(t, expected)
		}
	}

	xs := testCases{
		{
			name: "OK",
			run: run(
				metadata.NewIncomingContextF(context.Background(), &api.Metadata{
					Parent: "projects/test",
				}),
				&v1alpha.CreateClusterRequest{
					ClusterId: "foo",
					Cluster: &v1alpha.Cluster{
						Replicas: 1,
					},
				},
				&v1alpha.Cluster{
					Name:      "clusters/foo",
					Replicas:  1,
					Instances: []string{"bar", "baz"},
					Parent:    "projects/test",
				},
			),
			mock: func(r *Mockrdb) {
				r.EXPECT().Create(
					"foo", "test",
					test.ProtoEq(&v1alpha.Cluster{
						Replicas: 1,
					}),
				).Return(&v1alpha.Cluster{
					Name:      "clusters/foo",
					Replicas:  1,
					Instances: []string{"bar", "baz"},
					Parent:    "projects/test",
				}, nil)
			},
			check: test.Succeeded,
		},
	}

	xs.run(t)
}
