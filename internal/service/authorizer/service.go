package authorizer

import (
	"context"
	"fmt"

	v2 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v2"
	"github.com/nokamoto/demo20-apis/cloud/api"
	"github.com/nokamoto/demo20-apis/cloud/iam/admin/v1alpha"
	"github.com/nokamoto/demo20-apps/internal/service/core/incall"
	"github.com/nokamoto/demo20-apps/pkg/sdk/metadata"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type service struct {
	logger *zap.Logger
	iam    client
	config *api.AuthzConfig
}

// NewService returns AuthorizationServer.
func NewService(logger *zap.Logger, iam client, config *api.AuthzConfig) v2.AuthorizationServer {
	return &service{logger: logger, iam: iam}
}

func (s *service) validateCheckRequest(ctx context.Context, req *v2.CheckRequest) (string, *api.Metadata, error) {
	http := req.GetAttributes().GetRequest().GetHttp()
	header := http.GetHeaders()[metadata.MetadataKey]
	md, err := metadata.Decode(header)
	if err != nil {
		return "", nil, err
	}
	return http.GetPath(), md, nil
}

func (s *service) Check(ctx context.Context, req *v2.CheckRequest) (*v2.CheckResponse, error) {
	logger := incall.NewInCall(s.logger, "Check", req).Logger

	path, md, err := s.validateCheckRequest(ctx, req)
	if err != nil {
		logger.Debug("invalid request", zap.Error(err))
		return permissionDenied(fmt.Sprintf("invalid metadata: %s", metadata.MetadataKey)), nil
	}

	authz, found := s.config.GetConfig()[path]
	if !found {
		logger.Debug("path not found")
		return permissionDenied(fmt.Sprintf("path not found: %s", path)), nil
	}

	if authz.GetInsecure() {
		logger.Debug("insecure: authorization not required")
	}

	res, err := s.iam.AuthorizeMachineUser(context.Background(), &v1alpha.AuthorizeMachineUserRequest{
		ApiKey:     md.GetMachineUserApiKey(),
		Permission: authz.GetPermission(),
		Parent:     md.GetParent(),
	})

	code := status.Code(err)
	if code == codes.OK {
		md.User = &api.Metadata_MachineUser{
			MachineUser: res.GetMachineUser().GetName(),
		}
		value, err := metadata.Encode(md)
		if err != nil {
			logger.Error("failed to encode metadata", zap.Error(err))
			return permissionDenied("internal"), nil
		}
		return ok(value), nil
	}

	logger.Debug("permission denied", zap.Error(err))
	return permissionDenied(err.Error()), nil
}
