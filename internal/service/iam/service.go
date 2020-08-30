package iam

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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
	scoped := incall.NewInCall(s.logger, "CreateMachineUser", req)

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

// AuthorizeMachineUser authorizes the machine user.
func (s *service) AuthorizeMachineUser(ctx context.Context, req *admin.AuthorizeMachineUserRequest) (*admin.AuthorizeMachineUserResponse, error) {
	scoped := incall.NewInCall(s.logger, "AuthorizeMachineUser", req)

	authn, err := s.iam.AuthenticateMachineUser(req.GetApiKey())
	if err != nil {
		return nil, scoped.Error(err)
	}

	authz, err := s.iam.AuthorizeMachineUser(authn, req.GetParent(), req.GetPermission())
	if err != nil {
		return nil, scoped.Error(err)
	}

	if !authz {
		return nil, status.Error(codes.PermissionDenied, "unauthorized")
	}

	return &admin.AuthorizeMachineUserResponse{
		MachineUser: authn,
	}, nil
}

func (s *service) validateCreateRole(ctx context.Context, req *admin.CreateRoleRequest) []error {
	return validation.Concat(
		validation.ID(req.GetRoleId()),
		validation.Empty(req.GetRole().GetName()),
		validation.Empty(req.GetRole().GetParent()),
	)
}

// CreateRole creates a role.
func (s *service) CreateRole(ctx context.Context, req *admin.CreateRoleRequest) (*v1alpha.Role, error) {
	scoped := incall.NewInCall(s.logger, "CreateRole", req)

	errs := s.validateCreateRole(ctx, req)
	if len(errs) != 0 {
		return nil, scoped.InvalidArgument(errs)
	}

	res, err := s.iam.CreateRole(req.GetRoleId(), "/", req.GetRole())
	if err != nil {
		return nil, scoped.Error(err)
	}

	return res, nil
}

func (s *service) validateAddRoleBinding(ctx context.Context, req *admin.AddRoleBindingRequest) []error {
	return validation.Concat(
		validation.NameOr(req.GetRoleBinding().GetRole(), []string{"roles"}, []string{"projects", "roles"}),
		validation.Empty(req.GetRoleBinding().GetParent()),
	)
}

// AddRoleBinding creates a role binding.
func (s *service) AddRoleBinding(ctx context.Context, req *admin.AddRoleBindingRequest) (*v1alpha.RoleBinding, error) {
	scoped := incall.NewInCall(s.logger, "AddRoleBinding", req)

	errs := s.validateAddRoleBinding(ctx, req)
	if len(errs) != 0 {
		return nil, scoped.InvalidArgument(errs)
	}

	res, err := s.iam.AddRoleBinding("/", req.GetRoleBinding())
	if err != nil {
		return nil, scoped.Error(err)
	}

	return res, nil
}
