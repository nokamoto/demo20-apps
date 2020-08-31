package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/nokamoto/demo20-apis/cloud/api"
	"github.com/nokamoto/demo20-apps/pkg/sdk/metadata"

	"github.com/nokamoto/demo20-apis/cloud/compute/v1alpha"
	"github.com/nokamoto/demo20-apps/internal/automatedtest"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/testing/protocmp"
)

func main() {
	automatedtest.Main(func(con *grpc.ClientConn) automatedtest.Scenarios {
		c := v1alpha.NewComputeClient(con)

		ctx := metadata.AppendToOutgoingContextF(context.Background(), &api.Metadata{
			Parent: "projects/todo",
		})

		return automatedtest.Scenarios{
			{
				Name: "create an instance",
				Run: func(state automatedtest.State, logger *zap.Logger) (automatedtest.State, error) {
					res, err := c.CreateInstance(ctx, &v1alpha.CreateInstanceRequest{
						Instance: &v1alpha.Instance{
							Labels: []string{"foo", "bar"},
						},
					})
					if err != nil {
						return nil, err
					}

					if !strings.HasPrefix(res.GetName(), "instances/") {
						return nil, fmt.Errorf("unexpected prefix: %v", res)
					}

					expected := &v1alpha.Instance{
						Parent: "projects/todo",
						Labels: []string{"foo", "bar"},
					}

					return state, automatedtest.Diff(expected, res, protocmp.IgnoreFields(&v1alpha.Instance{}, "name"))
				},
			},
		}
	})
}
