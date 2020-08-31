package main

import (
	"context"
	"fmt"

	"github.com/nokamoto/demo20-apis/cloud/api"
	"github.com/nokamoto/demo20-apis/cloud/rdb/v1alpha"
	"github.com/nokamoto/demo20-apps/internal/automatedtest"
	"github.com/nokamoto/demo20-apps/pkg/sdk/metadata"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/testing/protocmp"
)

func main() {
	automatedtest.Main(func(con *grpc.ClientConn) automatedtest.Scenarios {
		c := v1alpha.NewRdbClient(con)

		ctx := metadata.AppendToOutgoingContextF(context.Background(), &api.Metadata{
			Parent: "projects/todo",
		})

		return automatedtest.Scenarios{
			{
				Name: "create a cluster",
				Run: func(state automatedtest.State, logger *zap.Logger) (automatedtest.State, error) {
					id := automatedtest.RandomID()

					res, err := c.CreateCluster(ctx, &v1alpha.CreateClusterRequest{
						ClusterId: id,
						Cluster: &v1alpha.Cluster{
							Replicas: 1,
						},
					})
					if err != nil {
						return nil, err
					}

					if len(res.GetInstances()) != 2 {
						return nil, fmt.Errorf("expected len(instances) is 2: %v", res.GetInstances())
					}

					expected := &v1alpha.Cluster{
						Name:     fmt.Sprintf("clusters/%s", id),
						Replicas: 1,
						Parent:   "projects/todo",
					}

					return state, automatedtest.Diff(expected, res, protocmp.IgnoreFields(&v1alpha.Cluster{}, "instances"))
				},
			},
		}
	})
}
