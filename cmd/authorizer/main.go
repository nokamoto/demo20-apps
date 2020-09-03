package main

import (
	"os"

	v2 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v2"
	"github.com/nokamoto/demo20-apis/cloud/iam/admin/v1alpha"
	"github.com/nokamoto/demo20-apps/internal/server"
	"github.com/nokamoto/demo20-apps/internal/service/authorizer"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

const (
	iamGrpcAddress = "IAM_GRPC_ADDRESS"
	configPath     = "CONFIG_PATH"
)

func main() {
	server.MainWithoutMySQL(func(logger *zap.Logger, s *grpc.Server) {
		address := os.Getenv(iamGrpcAddress)
		conn, err := grpc.Dial(address, grpc.WithInsecure())
		if err != nil {
			logger.Fatal("failed to create conn", zap.Error(err))
		}

		loader := authorizer.ConfigLoader{}
		cfg, err := loader.Read(os.Getenv(configPath))
		if err != nil {
			logger.Fatal("failed to load config", zap.Error(err))
		}

		v2.RegisterAuthorizationServer(s, authorizer.NewService(logger, v1alpha.NewIamClient(conn), cfg))
	})
}
