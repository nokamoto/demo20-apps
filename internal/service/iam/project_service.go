package iam

import (
	"context"

	"github.com/nokamoto/demo20-apis/cloud/iam/v1alpha"
	"github.com/nokamoto/demo20-apps/internal/service/core/validation"
	"go.uber.org/zap"
)

type projectService struct {
	v1alpha.UnimplementedIamServer
	service *service
}

// NewService returns v1alpha.IamServer.
func NewService(iam iam, logger *zap.Logger) v1alpha.IamServer {
	return &projectService{service: &service{iam: iam, logger: logger}}
}

func (p *projectService) CreateMachineUser(ctx context.Context, req *v1alpha.CreateMachineUserRequest) (*v1alpha.MachineUser, error) {
	return p.service.CreateMachineUser(ctx, req, validation.ProjectIncomingContext)
}

func (p *projectService) CreateRole(ctx context.Context, req *v1alpha.CreateRoleRequest) (*v1alpha.Role, error) {
	return p.service.CreateRole(ctx, req, validation.ProjectIncomingContext)
}

func (p *projectService) AddRoleBinding(ctx context.Context, req *v1alpha.AddRoleBindingRequest) (*v1alpha.RoleBinding, error) {
	return p.service.AddRoleBinding(ctx, req, validation.ProjectIncomingContext)
}
