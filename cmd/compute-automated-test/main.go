package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-cmp/cmp"
	"github.com/nokamoto/demo20-apis/cloud/compute/v1alpha"
	"github.com/nokamoto/demo20-apps/internal/automatedtest"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/testing/protocmp"
)

func main() {
	automatedtest.Main(func(con *grpc.ClientConn) automatedtest.Scenarios {
		c := v1alpha.NewComputeClient(con)

		return automatedtest.Scenarios{
			{
				Name: "create an instance",
				Run: func(state automatedtest.State, logger *zap.Logger) (automatedtest.State, error) {
					res, err := c.CreateInstance(context.Background(), &v1alpha.CreateInstanceRequest{
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

					logger.Debug("ignore fields", zap.String("name", res.GetName()))
					res.Name = ""

					expected := &v1alpha.Instance{
						Parent: "projects/todo",
						Labels: []string{"foo", "bar"},
					}

					if diff := cmp.Diff(expected, res, protocmp.Transform()); len(diff) != 0 {
						fmt.Println(diff)
						return nil, fmt.Errorf("unexpected response: %s", diff)
					}

					return state, nil
				},
			},
		}
	})
}
