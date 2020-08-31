package main

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/google/go-cmp/cmp"
	admin "github.com/nokamoto/demo20-apis/cloud/iam/admin/v1alpha"
	"github.com/nokamoto/demo20-apis/cloud/iam/v1alpha"
	"github.com/nokamoto/demo20-apps/internal/automatedtest"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/testing/protocmp"
)

func main() {
	automatedtest.Main(func(con *grpc.ClientConn) automatedtest.Scenarios {
		adminClient := admin.NewIamClient(con)
		projectClient := v1alpha.NewIamClient(con)

		return automatedtest.Scenarios{
			{
				Name: "create a permission",
				Run: func(state automatedtest.State, logger *zap.Logger) (automatedtest.State, error) {
					id := automatedtest.RandomID()

					res, err := adminClient.CreatePermission(context.Background(), &admin.CreatePermissionRequest{
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

					state["permission"] = proto.MarshalTextString(res)

					return state, nil
				},
			},
			{
				Name: "create a machine user",
				Run: func(state automatedtest.State, logger *zap.Logger) (automatedtest.State, error) {
					res, err := adminClient.CreateMachineUser(context.Background(), &admin.CreateMachineUserRequest{
						MachineUser: &v1alpha.MachineUser{
							DisplayName: "test machine user",
						},
					})
					if err != nil {
						return nil, err
					}

					expected := &v1alpha.MachineUser{
						DisplayName: "test machine user",
						Parent:      "projects//",
					}

					if diff := cmp.Diff(expected, res, protocmp.Transform(), protocmp.IgnoreFields(&v1alpha.MachineUser{}, "name", "api_key")); len(diff) != 0 {
						return nil, fmt.Errorf("unexpected response: %s", diff)
					}

					state["machineuser"] = proto.MarshalTextString(res)

					return state, nil
				},
			},
			{
				Name: "create a role",
				Run: func(state automatedtest.State, logger *zap.Logger) (automatedtest.State, error) {
					id := automatedtest.RandomID()

					var permission v1alpha.Permission
					err := proto.UnmarshalText(state["permission"], &permission)
					if err != nil {
						return nil, err
					}

					res, err := adminClient.CreateRole(context.Background(), &admin.CreateRoleRequest{
						RoleId: id,
						Role: &v1alpha.Role{
							DisplayName: "test display name",
							Permissions: []string{permission.GetName()},
						},
					})
					if err != nil {
						return nil, err
					}

					expected := &v1alpha.Role{
						Name:        fmt.Sprintf("roles/%s", id),
						DisplayName: "test display name",
						Permissions: []string{permission.GetName()},
						Parent:      "projects//",
					}

					if diff := cmp.Diff(expected, res, protocmp.Transform()); len(diff) != 0 {
						return nil, fmt.Errorf("unexpected response: %s", diff)
					}

					state["role"] = proto.MarshalTextString(res)

					return state, nil
				},
			},
			{
				Name: "add a role binding",
				Run: func(state automatedtest.State, logger *zap.Logger) (automatedtest.State, error) {
					var role v1alpha.Role
					err := proto.UnmarshalText(state["role"], &role)
					if err != nil {
						return nil, err
					}

					var machineUser v1alpha.MachineUser
					err = proto.UnmarshalText(state["machineuser"], &machineUser)
					if err != nil {
						return nil, err
					}

					res, err := adminClient.AddRoleBinding(context.Background(), &admin.AddRoleBindingRequest{
						RoleBinding: &v1alpha.RoleBinding{
							Role: role.GetName(),
							User: machineUser.GetName(),
						},
					})
					if err != nil {
						return nil, err
					}

					expected := &v1alpha.RoleBinding{
						Role:   role.GetName(),
						User:   machineUser.GetName(),
						Parent: "projects//",
					}

					if diff := cmp.Diff(expected, res, protocmp.Transform()); len(diff) != 0 {
						return nil, fmt.Errorf("unexpected response: %s", diff)
					}

					state["rolebinding"] = proto.MarshalTextString(res)

					return state, nil
				},
			},
			{
				Name: "authorize the machine user",
				Run: func(state automatedtest.State, logger *zap.Logger) (automatedtest.State, error) {
					var permission v1alpha.Permission
					err := proto.UnmarshalText(state["permission"], &permission)
					if err != nil {
						return nil, err
					}

					var machineUser v1alpha.MachineUser
					err = proto.UnmarshalText(state["machineuser"], &machineUser)
					if err != nil {
						return nil, err
					}

					res, err := adminClient.AuthorizeMachineUser(context.Background(), &admin.AuthorizeMachineUserRequest{
						ApiKey:     machineUser.GetApiKey(),
						Permission: permission.GetName(),
						Parent:     "projects//",
					})
					if err != nil {
						return nil, err
					}

					expected := &admin.AuthorizeMachineUserResponse{
						MachineUser: &v1alpha.MachineUser{
							Name:        machineUser.GetName(),
							DisplayName: machineUser.GetDisplayName(),
							Parent:      machineUser.GetParent(),
						},
					}

					if diff := cmp.Diff(expected, res, protocmp.Transform()); len(diff) != 0 {
						return nil, fmt.Errorf("unexpected response: %s", diff)
					}

					return state, nil
				},
			},
			{
				Name: "create a project role",
				Run: func(state automatedtest.State, logger *zap.Logger) (automatedtest.State, error) {
					permissionID := automatedtest.RandomID()

					permission, err := adminClient.CreatePermission(context.Background(), &admin.CreatePermissionRequest{
						PermissionId: permissionID,
					})
					if err != nil {
						return nil, err
					}

					roleID := automatedtest.RandomID()

					res, err := projectClient.CreateRole(context.Background(), &v1alpha.CreateRoleRequest{
						RoleId: roleID,
						Role: &v1alpha.Role{
							DisplayName: "test project role",
							Permissions: []string{permission.GetName()},
						},
					})
					if err != nil {
						return nil, err
					}

					expected := &v1alpha.Role{
						Name:        fmt.Sprintf("projects/%s/roles/%s", "todo", roleID),
						DisplayName: "test project role",
						Permissions: []string{permission.GetName()},
						Parent:      "projects/todo",
					}

					if diff := cmp.Diff(expected, res, protocmp.Transform()); len(diff) != 0 {
						return nil, fmt.Errorf("unexpected response: %s", diff)
					}

					state["projectpermission"] = proto.MarshalTextString(permission)
					state["projectrole"] = proto.MarshalTextString(res)

					return state, nil
				},
			},
			{
				Name: "create a project machine user",
				Run: func(state automatedtest.State, logger *zap.Logger) (automatedtest.State, error) {
					res, err := projectClient.CreateMachineUser(context.Background(), &v1alpha.CreateMachineUserRequest{
						MachineUser: &v1alpha.MachineUser{
							DisplayName: "test project machine user",
						},
					})
					if err != nil {
						return nil, err
					}

					expected := &v1alpha.MachineUser{
						DisplayName: "test project machine user",
						Parent:      "projects/todo",
					}

					if diff := cmp.Diff(expected, res, protocmp.Transform(), protocmp.IgnoreFields(&v1alpha.MachineUser{}, "name", "api_key")); len(diff) != 0 {
						return nil, fmt.Errorf("unexpected response: %s", diff)
					}

					state["projectmachineuser"] = proto.MarshalTextString(res)

					return state, nil
				},
			},
			{
				Name: "create a project role binding",
				Run: func(state automatedtest.State, logger *zap.Logger) (automatedtest.State, error) {
					var role v1alpha.Role
					err := proto.UnmarshalText(state["projectrole"], &role)
					if err != nil {
						return nil, err
					}

					var machineUser v1alpha.MachineUser
					err = proto.UnmarshalText(state["projectmachineuser"], &machineUser)
					if err != nil {
						return nil, err
					}

					res, err := projectClient.AddRoleBinding(context.Background(), &v1alpha.AddRoleBindingRequest{
						RoleBinding: &v1alpha.RoleBinding{
							Role: role.GetName(),
							User: machineUser.GetName(),
						},
					})
					if err != nil {
						return nil, err
					}

					expected := &v1alpha.RoleBinding{
						Role:   role.GetName(),
						User:   machineUser.GetName(),
						Parent: "projects/todo",
					}

					if diff := cmp.Diff(expected, res, protocmp.Transform()); len(diff) != 0 {
						return nil, fmt.Errorf("unexpected response: %s", diff)
					}

					return state, nil
				},
			},
			{
				Name: "authorize the project machine user",
				Run: func(state automatedtest.State, logger *zap.Logger) (automatedtest.State, error) {
					var permission v1alpha.Permission
					err := proto.UnmarshalText(state["projectpermission"], &permission)
					if err != nil {
						return nil, err
					}

					var machineUser v1alpha.MachineUser
					err = proto.UnmarshalText(state["projectmachineuser"], &machineUser)
					if err != nil {
						return nil, err
					}

					res, err := adminClient.AuthorizeMachineUser(context.Background(), &admin.AuthorizeMachineUserRequest{
						ApiKey:     machineUser.GetApiKey(),
						Permission: permission.GetName(),
						Parent:     "projects/todo",
					})
					if err != nil {
						return nil, err
					}

					expected := &admin.AuthorizeMachineUserResponse{
						MachineUser: &v1alpha.MachineUser{
							Name:        machineUser.GetName(),
							DisplayName: machineUser.GetDisplayName(),
							Parent:      machineUser.GetParent(),
						},
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
