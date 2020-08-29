package main

import (
	"context"

	"github.com/nokamoto/demo20-apis/cloud/resourcemanager/v1alpha"
	"github.com/nokamoto/demo20-apps/internal/automatedtest"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	automatedtest.Main(func(con *grpc.ClientConn) automatedtest.Scenarios {
		c := v1alpha.NewResourceManagerClient(con)

		return automatedtest.Scenarios{
			{
				Name: "get the project",
				Run: func(state automatedtest.State, logger *zap.Logger) (automatedtest.State, error) {
					res, err := c.GetProject(context.Background(), &v1alpha.GetProjectRequest{
						Name: "projects/todo",
					})
					if err != nil {
						return nil, err
					}

					expected := &v1alpha.Project{
						Name:        "projects/todo",
						DisplayName: "todo display name",
					}

					return state, automatedtest.Diff(expected, res)
				},
			},
		}
	})
}
