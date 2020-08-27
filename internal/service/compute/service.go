package compute

import (
	"context"

	"github.com/nokamoto/demo20-apis/cloud/compute/v1alpha"
	"github.com/nokamoto/demo20-apps/internal/service/core/incall"
	"go.uber.org/zap"
)

type service struct {
	compute compute
	logger  *zap.Logger
	v1alpha.UnimplementedComputeServer
}

func (s *service) validateCreateInstance(ctx context.Context, req *v1alpha.CreateInstanceRequest) (string, []error) {
	return "todo", nil
}

func (s *service) CreateInstance(ctx context.Context, req *v1alpha.CreateInstanceRequest) (*v1alpha.Instance, error) {
	scoped := incall.NewInCall(s.logger, "CreateInstance", req)

	parentID, errs := s.validateCreateInstance(ctx, req)
	if len(errs) != 0 {
		return nil, scoped.InvalidArgument(errs)
	}

	res, err := s.compute.Create(
		s.compute.RandomName(parentID),
		parentID,
		&v1alpha.Instance{
			Labels: req.GetInstance().GetLabels(),
		},
	)
	if err != nil {
		return nil, scoped.Error(err)
	}

	return res, nil
}