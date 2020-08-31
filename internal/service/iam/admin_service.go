package iam

import (
	"context"

	admin "github.com/nokamoto/demo20-apis/cloud/iam/admin/v1alpha"
	"github.com/nokamoto/demo20-apis/cloud/iam/v1alpha"
	"go.uber.org/zap"
)

type adminService struct {
	admin.UnimplementedIamServer
	service *service
}

// NewAdminService returns admin.v1alpha.IamServer.
func NewAdminService(iam iam, logger *zap.Logger) admin.IamServer {
	return &adminService{service: &service{iam: iam, logger: logger}}
}

func (a *adminService) CreatePermission(ctx context.Context, req *admin.CreatePermissionRequest) (*v1alpha.Permission, error) {
	return a.service.CreatePermission(ctx, req)
}

func (a *adminService) CreateMachineUser(ctx context.Context, req *admin.CreateMachineUserRequest) (*v1alpha.MachineUser, error) {
	return a.service.CreateMachineUser(ctx, req, "/")
}

func (a *adminService) AuthorizeMachineUser(ctx context.Context, req *admin.AuthorizeMachineUserRequest) (*admin.AuthorizeMachineUserResponse, error) {
	return a.service.AuthorizeMachineUser(ctx, req)
}

func (a *adminService) CreateRole(ctx context.Context, req *admin.CreateRoleRequest) (*v1alpha.Role, error) {
	return a.service.CreateRole(ctx, req, "/")
}

func (a *adminService) AddRoleBinding(ctx context.Context, req *admin.AddRoleBindingRequest) (*v1alpha.RoleBinding, error) {
	return a.service.AddRoleBinding(ctx, req, "/")
}
