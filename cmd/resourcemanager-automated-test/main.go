package main

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/proto"

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
				Name: "create a project",
				Run: func(state automatedtest.State, logger *zap.Logger) (automatedtest.State, error) {
					id := automatedtest.RandomID()

					expected := &v1alpha.Project{
						Name:        fmt.Sprintf("projects/%s", id),
						DisplayName: "test project",
					}

					res, err := c.CreateProject(context.Background(), &v1alpha.CreateProjectRequest{
						ProjectId: id,
						Project: &v1alpha.Project{
							DisplayName: expected.GetDisplayName(),
						},
					})
					if err != nil {
						return nil, err
					}

					err = automatedtest.Diff(expected, res)
					if err != nil {
						return nil, err
					}

					state["project"] = proto.MarshalTextString(expected)

					return state, nil
				},
			},
			{
				Name: "get the project",
				Run: func(state automatedtest.State, logger *zap.Logger) (automatedtest.State, error) {
					var expected v1alpha.Project
					err := proto.UnmarshalText(state["project"], &expected)
					if err != nil {
						return nil, err
					}

					res, err := c.GetProject(context.Background(), &v1alpha.GetProjectRequest{
						Name: expected.GetName(),
					})
					if err != nil {
						return nil, err
					}

					return state, automatedtest.Diff(&expected, res)
				},
			},
		}
	})
}
