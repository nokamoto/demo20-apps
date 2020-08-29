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

func (s *service) validateGetProject(ctx context.Context, req *v1alpha.GetProjectRequest) ([]string, []error) {
	var ids []string
	return ids, validation.Concat(validation.FromName(req.GetName(), &ids, "projects"))
}

// GetProject returns a project.
func (s *service) GetProject(ctx context.Context, req *v1alpha.GetProjectRequest) (*v1alpha.Project, error) {
	scoped := incall.NewInCall(s.logger, "GetProject", req)

	ids, errs := s.validateGetProject(ctx, req)
	if len(errs) != 0 {
		return nil, scoped.InvalidArgument(errs)
	}

	res, err := s.resourcemanager.Get(ids[0])
	if err != nil {
		return nil, scoped.Error(err)
	}

	return res, nil
}
