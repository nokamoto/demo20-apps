package resourcemanager

import (
	"context"

	"github.com/nokamoto/demo20-apis/cloud/resourcemanager/v1alpha"
	"github.com/nokamoto/demo20-apps/internal/service/core/incall"
	"github.com/nokamoto/demo20-apps/internal/service/core/validation"
	"go.uber.org/zap"
)

type service struct {
	resourcemanager resourcemanager
	v1alpha.UnimplementedResourceManagerServer
	logger *zap.Logger
}

// NewService returns admin.v1alpha.IamServer.
func NewService(resourcemanager resourcemanager, logger *zap.Logger) v1alpha.ResourceManagerServer {
	return &service{resourcemanager: resourcemanager, logger: logger}
}

func (s *service) validateGetProject(ctx context.Context, req *v1alpha.GetProjectRequest) (string, []error) {
	id, err := validation.ProjectIncomingContext(ctx)
	return id, validation.Concat(err)
}

// GetProject returns a project.
func (s *service) GetProject(ctx context.Context, req *v1alpha.GetProjectRequest) (*v1alpha.Project, error) {
	scoped := incall.NewInCall(s.logger, "GetProject", req)

	parentID, errs := s.validateGetProject(ctx, req)
	if len(errs) != 0 {
		return nil, scoped.InvalidArgument(errs)
	}

	res, err := s.resourcemanager.Get(parentID)
	if err != nil {
		return nil, scoped.Error(err)
	}

	return res, nil
}

func (s *service) validateCreateProject(ctx context.Context, req *v1alpha.CreateProjectRequest) (string, []error) {
	return req.GetProjectId(), validation.Concat(
		validation.ID(req.GetProjectId()),
		validation.Empty(req.GetProject().GetName()),
	)
}

// CreateProject creates a project.
func (s *service) CreateProject(ctx context.Context, req *v1alpha.CreateProjectRequest) (*v1alpha.Project, error) {
	scoped := incall.NewInCall(s.logger, "CreateProject", req)

	projectID, errs := s.validateCreateProject(ctx, req)
	if len(errs) != 0 {
		return nil, scoped.InvalidArgument(errs)
	}

	res, err := s.resourcemanager.Create(projectID, req.GetProject())
	if err != nil {
		return nil, scoped.Error(err)
	}

	return res, nil
}
