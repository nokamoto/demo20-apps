package iam

import (
	"context"

	"github.com/golang/protobuf/proto"
	admin "github.com/nokamoto/demo20-apis/cloud/iam/admin/v1alpha"
	"github.com/nokamoto/demo20-apis/cloud/iam/v1alpha"
	"github.com/nokamoto/demo20-apps/internal/service/core/incall"
	"github.com/nokamoto/demo20-apps/internal/service/core/validation"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type service struct {
	iam    iam
	logger *zap.Logger
}

type createPermissionRequest interface {
	proto.Message
	GetPermissionId() string
	GetPermission() *v1alpha.Permission
}

type projectFromIncomingContext func(context.Context) (string, error)

func (s *service) validateCreatePermission(ctx context.Context, req createPermissionRequest) []error {
	return validation.Concat(
		validation.ID(req.GetPermissionId()),
		validation.Empty(req.GetPermission().GetName()),
	)
}

// CreatePermission creates a permission.
func (s *service) CreatePermission(ctx context.Context, req createPermissionRequest) (*v1alpha.Permission, error) {
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

type createMachineUserRequest interface {
	proto.Message
	GetMachineUser() *v1alpha.MachineUser
}

func (s *service) validateCreateMachineUser(ctx context.Context, req createMachineUserRequest, project projectFromIncomingContext) (string, []error) {
	id, err := project(ctx)
	return id, validation.Concat(
		err,
		validation.Empty(req.GetMachineUser().GetName()),
		validation.Empty(req.GetMachineUser().GetApiKey()),
	)
}

// CreateMachineUser creates a machine user.
func (s *service) CreateMachineUser(ctx context.Context, req createMachineUserRequest, project projectFromIncomingContext) (*v1alpha.MachineUser, error) {
	scoped := incall.NewInCall(s.logger, "CreateMachineUser", req)

	projectID, errs := s.validateCreateMachineUser(ctx, req, project)
	if len(errs) != 0 {
		return nil, scoped.InvalidArgument(errs)
	}

	res, err := s.iam.CreateMachineUser(projectID, req.GetMachineUser())
	if err != nil {
		return nil, scoped.Error(err)
	}

	return res, nil
}

type authorizeMachineUserRequest interface {
	proto.Message
	GetApiKey() string
	GetParent() string
	GetPermission() string
}

// AuthorizeMachineUser authorizes the machine user.
func (s *service) AuthorizeMachineUser(ctx context.Context, req authorizeMachineUserRequest) (*admin.AuthorizeMachineUserResponse, error) {
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

type createRoleRequest interface {
	proto.Message
	GetRoleId() string
	GetRole() *v1alpha.Role
}

func (s *service) validateCreateRole(ctx context.Context, req createRoleRequest, project projectFromIncomingContext) (string, []error) {
	id, err := project(ctx)
	return id, validation.Concat(
		err,
		validation.ID(req.GetRoleId()),
		validation.Empty(req.GetRole().GetName()),
		validation.Empty(req.GetRole().GetParent()),
	)
}

// CreateRole creates a role.
func (s *service) CreateRole(ctx context.Context, req createRoleRequest, project projectFromIncomingContext) (*v1alpha.Role, error) {
	scoped := incall.NewInCall(s.logger, "CreateRole", req)

	projectID, errs := s.validateCreateRole(ctx, req, project)
	if len(errs) != 0 {
		return nil, scoped.InvalidArgument(errs)
	}

	res, err := s.iam.CreateRole(req.GetRoleId(), projectID, req.GetRole())
	if err != nil {
		return nil, scoped.Error(err)
	}

	return res, nil
}

type addRoleBindingRequest interface {
	proto.Message
	GetRoleBinding() *v1alpha.RoleBinding
}

func (s *service) validateAddRoleBinding(ctx context.Context, req addRoleBindingRequest, project projectFromIncomingContext) (string, []error) {
	id, err := project(ctx)
	return id, validation.Concat(
		err,
		validation.NameOr(req.GetRoleBinding().GetRole(), []string{"roles"}, []string{"projects", "roles"}),
		validation.Empty(req.GetRoleBinding().GetParent()),
	)
}

// AddRoleBinding creates a role binding.
func (s *service) AddRoleBinding(ctx context.Context, req addRoleBindingRequest, project projectFromIncomingContext) (*v1alpha.RoleBinding, error) {
	scoped := incall.NewInCall(s.logger, "AddRoleBinding", req)

	projectID, errs := s.validateAddRoleBinding(ctx, req, project)
	if len(errs) != 0 {
		return nil, scoped.InvalidArgument(errs)
	}

	res, err := s.iam.AddRoleBinding(projectID, req.GetRoleBinding())
	if err != nil {
		return nil, scoped.Error(err)
	}

	return res, nil
}
