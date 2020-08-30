package iam

import (
	"context"

	"github.com/nokamoto/demo20-apps/internal/service/core/incall"
	"github.com/nokamoto/demo20-apps/internal/service/core/validation"

	admin "github.com/nokamoto/demo20-apis/cloud/iam/admin/v1alpha"
	"github.com/nokamoto/demo20-apis/cloud/iam/v1alpha"
	"go.uber.org/zap"
)

type service struct {
	iam iam
	admin.UnimplementedIamServer
	logger *zap.Logger
}

// NewAdminService returns admin.v1alpha.IamServer.
func NewAdminService(iam iam, logger *zap.Logger) admin.IamServer {
	return &service{iam: iam, logger: logger}
}

func (s *service) validateCreatePermission(ctx context.Context, req *admin.CreatePermissionRequest) []error {
	return validation.Concat(
		validation.ID(req.GetPermissionId()),
		validation.Empty(req.GetPermission().GetName()),
	)
}

// CreatePermission creates a permission.
func (s *service) CreatePermission(ctx context.Context, req *admin.CreatePermissionRequest) (*v1alpha.Permission, error) {
	scoped := incall.NewInCall(s.logger, "CreatePermission", req)

	errs := s.validateCreatePermission(ctx, req)
	if len(errs) != 0 {
		return nil, scoped.InvalidArgument(errs)
	}

	res, err := s.iam.Create(req.GetPermissionId())
	if err != nil {
		return nil, scoped.Error(err)
	}

	return res, nil
}

func (s *service) validateCreateMachineUser(ctx context.Context, req *admin.CreateMachineUserRequest) []error {
	return validation.Concat(
		validation.Empty(req.GetMachineUser().GetName()),
		validation.Empty(req.GetMachineUser().GetApiKey()),
	)
}

// CreateMachineUser creates a machine user.
func (s *service) CreateMachineUser(ctx context.Context, req *admin.CreateMachineUserRequest) (*v1alpha.MachineUser, error) {
	scoped := incall.NewInCall(s.logger, "CreatePermission", req)

	errs := s.validateCreateMachineUser(ctx, req)
	if len(errs) != 0 {
		return nil, scoped.InvalidArgument(errs)
	}

	res, err := s.iam.CreateMachineUser("/", req.GetMachineUser())
	if err != nil {
		return nil, scoped.Error(err)
	}

	return res, nil
}
