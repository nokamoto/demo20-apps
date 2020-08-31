package rdb

import (
	"context"

	"github.com/nokamoto/demo20-apis/cloud/rdb/v1alpha"
	"github.com/nokamoto/demo20-apps/internal/service/core/incall"
	"github.com/nokamoto/demo20-apps/internal/service/core/validation"
	"go.uber.org/zap"
)

type service struct {
	rdb rdb
	v1alpha.UnimplementedRdbServer
	logger *zap.Logger
}

// NewService returns admin.v1alpha.IamServer.
func NewService(rdb rdb, logger *zap.Logger) v1alpha.RdbServer {
	return &service{rdb: rdb, logger: logger}
}

func (s *service) validateCreateCluster(ctx context.Context, req *v1alpha.CreateClusterRequest) (string, []error) {
	parentID, err := validation.ProjectIncomingContext(ctx)
	return parentID, validation.Concat(
		err,
		validation.ID(req.GetClusterId()),
		validation.Empty(req.GetCluster().GetName()),
		validation.Empty(req.GetCluster().GetParent()),
		validation.Range(int(req.GetCluster().GetReplicas()), 0, 32),
		validation.EmptyStrings(req.GetCluster().GetInstances()),
	)
}

// CreateCluster creates a cluster.
func (s *service) CreateCluster(ctx context.Context, req *v1alpha.CreateClusterRequest) (*v1alpha.Cluster, error) {
	scoped := incall.NewInCall(s.logger, "GetProject", req)

	parentID, errs := s.validateCreateCluster(ctx, req)
	if len(errs) != 0 {
		return nil, scoped.InvalidArgument(errs)
	}

	res, err := s.rdb.Create(req.GetClusterId(), parentID, req.GetCluster())
	if err != nil {
		return nil, scoped.Error(err)
	}

	return res, nil
}
