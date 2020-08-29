package main

import (
	"google.golang.org/protobuf/testing/protocmp"
	"context"
	"fmt"

	"github.com/google/go-cmp/cmp"
	"github.com/nokamoto/demo20-apis/cloud/iam/v1alpha"

	admin "github.com/nokamoto/demo20-apis/cloud/iam/admin/v1alpha"
	"github.com/nokamoto/demo20-apps/internal/automatedtest"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	automatedtest.Main(func(con *grpc.ClientConn) automatedtest.Scenarios {
		c := admin.NewIamClient(con)

		return automatedtest.Scenarios{
			{
				Name: "create a permission",
				Run: func(state automatedtest.State, logger *zap.Logger) (automatedtest.State, error) {
					id := automatedtest.RandomID()

					res, err := c.CreatePermission(context.Background(), &admin.CreatePermissionRequest{
						PermissionId: id,
					})
					if err != nil {
						return nil, err
					}

					expected := &v1alpha.Permission{
						Name: fmt.Sprintf("permissions/%s", id),
					}

					if diff := cmp.Diff(expected, res, protocmp.Transform()); len(diff) != 0 {
						return nil, fmt.Errorf("unexpected response: %s", diff)
					}

					return state, nil
				},
			},
		}
	})
}
